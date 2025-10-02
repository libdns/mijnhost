package mijnhost

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/pbergman/provider/test"
)

func TestProvider_Unmarshall(t *testing.T) {
	var provider *Provider
	var buf = `{
"api_key": "testkey",
"base_uri": "http://127.0.0.1:8080",
"debug": true
}`

	if err := json.Unmarshal([]byte(buf), &provider); err != nil {
		t.Fatal(err)
	}

	if provider.GetApiKey() != "testkey" {
		t.Fatalf("api key = %s; want %s", provider.GetApiKey(), "testkey")
	}

	if provider.GetBaseUri().String() != "http://127.0.0.1:8080" {
		t.Fatalf("base uri = %s; want %s", provider.GetBaseUri().String(), "http://127.0.0.1:8080")
	}

	if provider.GetDebug() != nil {
		t.Fatalf("expected debug to be nil got %+v", provider.GetDebug())
	}

	_ = provider.getClient()

	if provider.GetDebug() != os.Stdout {
		t.Fatalf("expected debug to be os.Stdout got %+v", provider.GetDebug())
	}
}

func TestProvider(t *testing.T) {

	var provider = &Provider{
		ApiKey: os.Getenv("API_KEY"),
	}

	if _, ok := os.LookupEnv("DEBUG"); ok {
		provider.Debug = true
	}

	test.RunProviderTests(t, provider)
}
