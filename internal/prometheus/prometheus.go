package prometheus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
)

//go:generate mockery --name=PrometheusAPI --inpackage --filename=prometheus_mock.go

func InitPrometheus(ctx context.Context, url string) (PrometheusAPI, error) {
	client, err := api.NewClient(api.Config{
		Address: url,
	})
	if err != nil {
		log.Errorf("Error creating client: %v\n", err)
		return nil, err
	}

	v1api := v1.NewAPI(client)
	promApi := NewPrometheusAPI(v1api)

	healthyURL := fmt.Sprintf("%s/-/healthy", url)
	err = retry.Do(
		func() error {
			resp, err := http.Get(healthyURL)
			if err != nil {
				return err
			}
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to connect to prometheus: %s", resp.Status)
			}
			return nil
		},
		retry.OnRetry(func(_ uint, err error) {
			log.Info("prometheus initialization failed")
			log.Error(err)
		}),
		retry.Delay(2*time.Second),
		retry.MaxJitter(3*time.Second),
		retry.Attempts(8),
		retry.LastErrorOnly(true),
		retry.Context(ctx),
	)

	return promApi, err
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
