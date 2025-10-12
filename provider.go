package mijnhost

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/libdns/libdns"
	"github.com/libdns/mijnhost/client"
	"github.com/pbergman/provider"
)

type Client interface {
	provider.Client
	provider.ZoneAwareClient
}

type Provider struct {
	// ApiKey used for authenticating the mijn.host api see:
	// https://mijn.host/api/doc/doc-343216#obtaining-your-api-key
	ApiKey string `json:"api_key"`
	// BaseUri used for the api calls and will default to https://mijn.host/api/v2/
	BaseUri *ApiBaseUri `json:"base_uri"`

	// DebugLevel sets the verbosity for logging API requests and responses.
	DebugLevel provider.OutputLevel `json:"debug_level"`
	// DebugOut defines the output destination for debug logs.
	// Defaults to standard output (STDOUT).
	DebugOut io.Writer `json:"-"`

	client Client
	mutex  sync.RWMutex
}

func (p *Provider) getClient() Client {
	if nil == p.client {

		if nil == p.BaseUri {
			p.BaseUri = DefaultApiBaseUri()
		}

		if nil == p.DebugOut {
			p.DebugOut = os.Stdout
		}

		p.client = client.NewApiClient(p)
	}

	return p.client
}

func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	return provider.GetRecords(ctx, &p.mutex, p.getClient(), zone)
}

func (p *Provider) AppendRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	return provider.AppendRecords(ctx, &p.mutex, p.getClient(), zone, recs)
}

func (p *Provider) SetRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	return provider.SetRecords(ctx, &p.mutex, p.getClient(), zone, recs)
}

func (p *Provider) DeleteRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	return provider.DeleteRecords(ctx, &p.mutex, p.getClient(), zone, recs)
}

func (p *Provider) ListZones(ctx context.Context) ([]libdns.Zone, error) {
	return provider.ListZones(ctx, &p.mutex, p.getClient())
}

// Interface guards
var (
	_ client.ApiClientConfig = (*Provider)(nil)
	_ libdns.RecordGetter    = (*Provider)(nil)
	_ libdns.RecordAppender  = (*Provider)(nil)
	_ libdns.RecordSetter    = (*Provider)(nil)
	_ libdns.RecordDeleter   = (*Provider)(nil)
	_ libdns.ZoneLister      = (*Provider)(nil)
)
