package client

import (
	"net/url"
)

type ApiClientConfig interface {
	GetApiKey() string
	GetBaseUri() *url.URL
}
