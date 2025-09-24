package mijnhost

import (
	"io"
	"net/url"
	"os"
	"sync"

	"github.com/libdns/libdns"
	"github.com/libdns/mijnhost/client"
)

type Provider struct {
	// ApiKey used for authenticating the mijn.host api see:
	// https://mijn.host/api/doc/doc-343216#obtaining-your-api-key
	ApiKey string `json:"api_key"`
	// Debug when true it will dump the http.Client request/response to os.Stdout
	// or you can change that by setting `DebugOut`
	Debug    bool      `json:"debug"`
	DebugOut io.Writer `json:"-"`
	// BaseUri used for the api calls and will default to https://mijn.host/api/v2/
	BaseUri *ApiBaseUri `json:"base_uri"`

	client *client.ApiClient
	mutex  sync.RWMutex
}

func (p *Provider) getClient() *client.ApiClient {
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

func (p *Provider) GetApiKey() string {
	return p.ApiKey
}

func (p *Provider) GetDebug() io.Writer {
	if p.Debug && p.DebugOut != nil {
		return p.DebugOut
	}
	return nil
}

func (p *Provider) GetBaseUri() *url.URL {

	if nil == p.BaseUri {
		return nil
	}

	return (*url.URL)(p.BaseUri)
}

func fqdn(name string) string {

	if name[len(name)-1] != '.' {
		return name + "."
	}

	return name
}

// Interface guards
var (
	_ client.ApiClientConfig = (*Provider)(nil)

	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
	_ libdns.ZoneLister     = (*Provider)(nil)
)
