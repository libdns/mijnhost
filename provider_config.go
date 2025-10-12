package mijnhost

import (
	"io"
	"net/url"
	"strconv"

	"github.com/pbergman/provider"
)

func DefaultApiBaseUri() *ApiBaseUri {
	return &ApiBaseUri{
		Scheme: "https",
		Host:   "mijn.host",
		Path:   "/api/v2/",
	}
}

type ApiBaseUri url.URL

func (a *ApiBaseUri) MarshalJSON() ([]byte, error) {
	return []byte((*url.URL)(a).String()), nil
}

func (a *ApiBaseUri) UnmarshalJSON(data []byte) error {

	if out, err := strconv.Unquote(string(data)); err == nil {
		data = []byte(out)
	}

	b, err := url.Parse(string(data))

	if err != nil {
		return err
	}

	*a = ApiBaseUri(*b)

	return nil
}

func (p *Provider) GetApiKey() string {
	return p.ApiKey
}

func (p *Provider) GetBaseUri() *url.URL {
	if nil == p.BaseUri {
		return nil
	}
	return (*url.URL)(p.BaseUri)
}

func (p *Provider) DebugOutputLevel() provider.OutputLevel {
	return p.DebugLevel
}

func (p *Provider) DebugOutput() io.Writer {
	return p.DebugOut
}

func (p *Provider) SetDebug(level provider.OutputLevel, writer io.Writer) {
	p.DebugLevel = level
	p.DebugOut = writer
}
