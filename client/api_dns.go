package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/libdns/libdns"
	"github.com/pbergman/provider"
)

type DNSRecord struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

func MarshallRRRecord(data *DNSRecord, zone string) *libdns.RR {
	return &libdns.RR{
		Name: libdns.RelativeName(data.Name, zone),
		Type: data.Type,
		TTL:  time.Duration(data.TTL) * time.Second,
		Data: data.Value,
	}
}

func MarshallDNSRecords(data *libdns.RR, zone string) *DNSRecord {
	var record = &DNSRecord{
		Type:  data.Type,
		Name:  libdns.AbsoluteName(data.Name, zone),
		Value: data.Data,
	}

	if record.Name[len(record.Name)-1] != '.' {
		record.Name += "."
	}

	switch data.TTL.Seconds() {
	case 300, 900, 3600, 10800, 21600, 43200, 86400:
		record.TTL = int(data.TTL.Seconds())
	default:
		record.TTL = 900
	}

	return record
}

func (a *ApiClient) SetDNSList(ctx context.Context, zone string, records []*libdns.RR) error {

	type setPayload struct {
		Records []*DNSRecord `json:"records"`
	}

	var payload = &setPayload{
		Records: make([]*DNSRecord, len(records)),
	}

	for i, c := 0, len(records); i < c; i++ {
		payload.Records[i] = MarshallDNSRecords(records[i], zone)
	}

	var buf = new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return err
	}

	var object status

	if err := a.fetch(ctx, a.toDnsPath(zone), http.MethodPut, buf, &object); err != nil {
		return err
	}

	if err := object.Error(); err != nil {
		return err
	}

	return nil
}
func (a *ApiClient) GetDNSList(ctx context.Context, zone string) ([]libdns.Record, error) {

	type recordData struct {
		Domain  string `json:"domain"`
		Records []*DNSRecord
	}

	var object struct {
		status
		Data *recordData `json:"data"`
	}

	if err := a.fetch(ctx, a.toDnsPath(zone), http.MethodGet, nil, &object); err != nil {
		return nil, err
	}

	var items = make([]libdns.Record, len(object.Data.Records))

	for i, record := range object.Data.Records {
		items[i] = MarshallRRRecord(record, zone)
	}

	return items, nil
}

var (
	_ provider.Client          = (*ApiClient)(nil)
	_ provider.ZoneAwareClient = (*ApiClient)(nil)
)
