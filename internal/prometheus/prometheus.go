package prometheus

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
)

//go:generate mockery --name=PrometheusAPI --inpackage --filename=prometheus_mock.go

func InitPrometheus(address string) (PrometheusAPI, error) {
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		log.Errorf("Error creating client: %v\n", err)
		return nil, err
	}

	v1api := v1.NewAPI(client)
	promApi := NewPrometheusAPI(v1api)

	return promApi, nil
}

type prometheusApi struct {
	v1 v1.API
}

type PrometheusAPI interface {
	Query(ctx context.Context, query string, ts time.Time) (model.Value, v1.Warnings, error)
}

func NewPrometheusAPI(v1 v1.API) *prometheusApi {
	return &prometheusApi{v1}
}

func (p *prometheusApi) Query(ctx context.Context, query string, ts time.Time) (model.Value, v1.Warnings, error) {
	return p.v1.Query(ctx, query, ts)
}
