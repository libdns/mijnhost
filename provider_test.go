package mijnhost

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/pbergman/provider"
	"github.com/pbergman/provider/test"
)

func TestProvider_Unmarshall(t *testing.T) {
	var object *Provider
	var buf = `{
"api_key": "testkey",
"base_uri": "http://127.0.0.1:8080",
"debug_level": 2
}`

	if err := json.Unmarshal([]byte(buf), &object); err != nil {
		t.Fatal(err)
	}

	if object.GetApiKey() != "testkey" {
		t.Fatalf("api key = %s; want %s", object.GetApiKey(), "testkey")
	}

	if object.GetBaseUri().String() != "http://127.0.0.1:8080" {
		t.Fatalf("base uri = %s; want %s", object.GetBaseUri().String(), "http://127.0.0.1:8080")
	}

	if object.DebugOutputLevel() != provider.OutputVeryVerbose {
		t.Fatalf("expected debug to be 2 got %d", object.DebugOutputLevel())
	}
}

func TestProvider(t *testing.T) {

	var object = &Provider{
		ApiKey: os.Getenv("API_KEY"),
	}

	if _, ok := os.LookupEnv("DEBUG"); ok {
		if x, ok := strconv.Atoi(os.Getenv("DEBUG")); ok == nil {
			object.DebugLevel = provider.OutputLevel(x)
		}
	}

	test.RunProviderTests(t, object, test.TestAll)
}
