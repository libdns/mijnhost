package client

import (
	"net/http"
)

type apiTransport struct {
	http.RoundTripper
	ApiClientConfig
}

func (a *apiTransport) RoundTrip(request *http.Request) (*http.Response, error) {

	if uri := a.GetBaseUri(); uri != nil {
		request.URL = uri.ResolveReference(request.URL)
	}

	request.Header.Set("accept", "application/json")
	request.Header.Set("content-type", "application/json")
	request.Header.Set("api-key", a.GetApiKey())
	request.Header.Set("user-agent", "libdns-client/1.0")

	return a.RoundTripper.RoundTrip(request)
}
