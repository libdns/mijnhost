package client

import (
	"context"
	"net/http"

	"github.com/pbergman/provider"
)

type domainsData struct {
	Domains []*Domain `json:"domains"`
}

type Domain struct {
	Id          int      `json:"id"`
	Domain      string   `json:"domain"`
	RenewalDate string   `json:"renewal_date"`
	Status      string   `json:"status"`
	StatusId    int      `json:"status_id"`
	Tags        []string `json:"tags"`
}

func (d *Domain) Name() string {
	return d.Domain
}

func (a *ApiClient) Domains(ctx context.Context) ([]provider.Domain, error) {

	var object struct {
		status
		Data *domainsData `json:"data"`
	}

	if err := a.fetch(ctx, "domains", http.MethodGet, nil, &object); err != nil {
		return nil, err
	}

	var domains = make([]provider.Domain, len(object.Data.Domains))

	for i, domain := 0, len(object.Data.Domains); i < domain; i++ {
		domains[i] = object.Data.Domains[i]
	}

	return domains, nil
}
