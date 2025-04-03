package mijnhost

import (
	"fmt"
	"time"

	"github.com/libdns/libdns"
)

type RecordResponse struct {
	Type  string        `json:"type"`
	Name  string        `json:"name"`
	Value string        `json:"value"`
	TTL   time.Duration `json:"ttl"`
}

type RecordsResponse struct {
	DNSRecords []RecordResponse `json:"data.records"`
}

type Record struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

type SavedRecord struct {
	DNSRecord struct {
		ID   uint   `json:"id"`
		Type string `json:"type"`
		Data struct {
			Priority  uint   `json:"priority,omitempty"`
			Value     string `json:"value"`
			Subdomain string `json:"subdomain"`
		}
	} `json:"dns_record"`
}

func (r *SavedRecord) libDNSRecord(zone string) libdns.Record {
	return libdns.Record{
		ID:    fmt.Sprintf("%d", r.DNSRecord.ID),
		Name:  libdns.RelativeName(r.DNSRecord.Data.Subdomain, zone),
		Type:  r.DNSRecord.Type,
		Value: r.DNSRecord.Data.Value,
	}
}

func (r *RecordResponse) libDNSRecord(zone string) libdns.Record {
	return libdns.Record{
		// ID:       fmt.Sprintf("%d", r.ID),
		Name:  libdns.RelativeName(r.Name, zone),
		Type:  r.Type,
		Value: r.Value,
		TTL:   r.TTL,
		// Priority: r.Priority,
	}
}

func libdnsToRecord(r libdns.Record) Record {
	return Record{
		Type: r.Type,

		Value: r.Value,
		Name:  r.Name,
		TTL:   int(r.TTL),
	}
}
