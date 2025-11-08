# Mijn Host for `libdns`

This package implements the libdns interfaces for the [Mijn Host API](https://mijn.host/api/doc/)

## Authenticating

To authenticate, you need to create an api key [here](https://mijn.host/cp/account/api/).

## Example

Here's a minimal example of how to get all your DNS records using this `libdns` provider

```go
package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	
	"github.com/pbergman/provider"
	"github.com/libdns/mijnhost"
)

func main() { 
	var p = &mijnhost.Provider{
		ApiKey: "***************************",
    }

	zones, err := p.ListZones(context.Background())

	if err != nil {
		panic(err)
	}

	var writer = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	for _, zone := range zones {
		records, err := p.GetRecords(context.Background(), zone.Name)

		if err != nil {
			panic(err)
		}

		for _, record := range provider.RecordIterator(&records) {
			_, _ = fmt.Fprintf(writer, "%s\t%v\t%s\t%s\n", record.Name, record.TTL.Seconds(), record.Type, record.Data)
		}

	}

	_ = writer.Flush()
}
```

## Debugging

This library provides the ability to debug the request/response communication with the API server.

To enable debugging, simply set the `debugging` property to `true`:
```go
	var provider = &mijnhost.Provider{
		ApiKey: "***************************",
		DebugLevel: provider.OutputDebug,
    }

	zones, err := provider.ListZones(context.Background())

	if err != nil {
		panic(err)
	}

	records, err := provider.GetRecords(context.Background(), "example.nl")
```

```shell
........................
[c] GET /api/v2/domains/example.nl/dns HTTP/1.1
[c] Host: mijn.host
[c] Accept: application/json
[c] Api-Key: ***************************
[c] Content-Type: application/json
[c] User-Agent: libdns-client/1.0
[c] 
[s] HTTP/2.0 200 OK
[s] Connection: close
[s] Content-Security-Policy: frame-ancestors 'self';
[s] Content-Type: application/json
[s] Date: Wed, 24 Sep 2025 08:46:46 GMT
[s] Feature-Policy: autoplay 'none'; camera 'none'; encrypted-media 'self'; fullscreen 'self'; geolocation 'none'; microphone 'none'; midi 'none'; payment 'none'; vr 'none'
[s] Permissions-Policy: accelerometer=(), ambient-light-sensor=(), autoplay=(self), battery=(), camera=(), cross-origin-isolated=(self), display-capture=(self), document-domain=(self), encrypted-media=(self), execution-while-not-rendered=(self), execution-while-out-of-viewport=(self), fullscreen=(self), geolocation=(self), gyroscope=(), keyboard-map=(self), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(self), picture-in-picture=(self), publickey-credentials-get=(self), screen-wake-lock=(), sync-xhr=(self), usb=(), web-share=(self), xr-spatial-tracking=(self), clipboard-read=(), clipboard-write=(self), gamepad=(), speaker-selection=(), conversion-measurement=(), focus-without-user-activation=(), hid=(), idle-detection=(), interest-cohort=(), serial=(), sync-script=(), trust-token-redemption=(), unload=(), window-placement=(), vertical-scroll=()
[s] Referrer-Policy: strict-origin-when-cross-origin
[s] Server: Apache/2
[s] Strict-Transport-Security: max-age=31536000
[s] Vary: User-Agent
[s] X-Content-Type-Options: nosniff
[s] X-Frame-Options: SAMEORIGIN
[s] X-Xss-Protection: 1; mode=block
[s] 
[s] {"status":200,"status_description":"Request successful","data":{"domain":"example.nl","records":.....
```

This will by default write to stdout but can set to any `io.Writer` by also setting the `DebugOut` property. 

```go
    var provider = &mijnhost.Provider{
        ApiKey: "***************************",
        DebugLevel: provider.OutputDebug,
        DebugOut: log.Writer(),
    }
```



## Testing

This library comes with a test suite that verifies the interface by creating a few test records, validating them, and then removing those records. To run the tests, you can use:

```shell
API_KEY=<MIJN_HOST_KEY> go test
```

Or run more verbose test to dump all api requests and responses: 

```shell
API_KEY=<MIJN_HOST_KEY> DEBUG=1 go test -v 
```

