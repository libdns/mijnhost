// Package libdnstemplate implements a DNS record management client compatible
// with the libdns interfaces for <PROVIDER NAME>. TODO: This package is a
// template only. Customize all godocs for actual implementation.
package mijnhost

import (
	"context"
	"fmt"
	"net/http"

	"github.com/libdns/libdns"
)

// TODO: Providers must not require additional provisioning steps by the callers; it
// should work simply by populating a struct and calling methods on it. If your DNS
// service requires long-lived state or some extra provisioning step, do it implicitly
// when methods are called; sync.Once can help with this, and/or you can use a
// sync.(RW)Mutex in your Provider struct to synchronize implicit provisioning.

// Provider facilitates DNS record manipulation with <TODO: PROVIDER NAME>.
type Provider struct {
	// TODO: put config fields here (with snake_case json
	// struct tags on exported fields), for example:
	APIToken string `json:"api_token,omitempty"`
	ApiURL   string `json:"api_url,omitempty"`
}

func (p *Provider) setDefaults() {
	if p.ApiURL == "" {
		p.ApiURL = "https://mijn.host/api/v2"
	}
}

func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.setDefaults()

	zone = normalizeZone(zone)
	reqURL := fmt.Sprintf("%s/domains/%s/dns", p.ApiURL, zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	var result RecordsResponse
	err = p.doAPIRequest(req, &result)

	recs := make([]libdns.Record, 0, len(result.DNSRecords))
	for _, r := range result.DNSRecords {
		recs = append(recs, r.libDNSRecord(zone))
	}

	return recs, err
}

func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.setDefaults()

	zone = normalizeZone(zone)
	var created []libdns.Record
	for _, record := range records {
		result, err := p.updateRecord(ctx, zone, record)
		if err != nil {
			return nil, err
		}
		created = append(created, result.libDNSRecord(zone))
	}

	return created, nil
}

func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.setDefaults()

	zone = normalizeZone(zone)
	var deleted []libdns.Record
	for _, record := range records {
		err := p.deleteRecord(ctx, zone, record)
		if err != nil {
			return nil, err
		}
		deleted = append(deleted, record)
	}

	return deleted, nil
}

func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.setDefaults()

	zone = normalizeZone(zone)
	zoneRecords, err := p.GetRecords(ctx, zone)
	if err != nil {
		return nil, err
	}

	var results []libdns.Record
	var resultErr error
	for _, libRecord := range records {
	
		record, err := p.updateRecord(ctx, zone, libRecord)
		if err != nil {
			resultErr = err
		}
		results = append(results, record.libDNSRecord(zone))
		
	}

	return results, resultErr
}

var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
