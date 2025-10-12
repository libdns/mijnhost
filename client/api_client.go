package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pbergman/provider"
)

func NewApiClient(x ApiClientConfig) *ApiClient {

	var transporter http.RoundTripper = &apiTransport{
		RoundTripper:    http.DefaultTransport,
		ApiClientConfig: x,
	}

	if v, ok := x.(provider.DebugConfig); ok {
		transporter = &provider.DebugTransport{
			RoundTripper: transporter,
			Config:       v,
		}
	}

	return &ApiClient{
		client: &http.Client{
			Transport: transporter,
		},
	}
}

type ApiClient struct {
	client *http.Client
}

func (a *ApiClient) toDnsPath(domain string) string {
	return fmt.Sprintf("domains/%s/dns", url.PathEscape(strings.TrimSuffix(domain, ".")))
}

func (a *ApiClient) fetch(ctx context.Context, path string, method string, body io.Reader, object any) error {

	request, err := http.NewRequestWithContext(ctx, method, path, body)

	if err != nil {
		return err
	}

	response, err := a.client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if !strings.HasPrefix(response.Header.Get("content-type"), "application/json") {
		return fmt.Errorf("unexpected response type: %s", response.Header.Get("content-type"))
	}

	if nil != object {

		if err := json.NewDecoder(response.Body).Decode(object); err != nil {
			return err
		}

		if v, ok := object.(StatusResponse); ok {
			if err := v.Error(); err != nil {
				return err
			}
		}
	}

	return nil
}
