package prometheus

import (
  log "github.com/sirupsen/logrus"
  "github.com/prometheus/client_golang/api"
  "github.com/prometheus/client_golang/api/prometheus/v1"
)

func InitPrometheus(address string) (v1.API, error) {
  client, err := api.NewClient(api.Config{
    Address: address,
	})
	if err != nil {
		log.Errorf("Error creating client: %v\n", err)
		return nil, err
	}

	v1api := v1.NewAPI(client)

  return v1api, nil
}
