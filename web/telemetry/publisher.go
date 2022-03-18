package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

type TelemetryPublisher struct {
	apiHost    string
	httpClient *http.Client
}

func (tp *TelemetryPublisher) Publish(telemetryName string, installationID uuid.UUID, extractedTelemetry interface{}) error {
	endpoint := fmt.Sprintf("%s/api/collect/hosts", tp.apiHost)
	// here we need to know what we are publishing
	// /api/collect/hosts works ok for hosts, not for the rest
	// What about /api/collect/:installationID/:telemetryName?
	// Or /api/collect and the installationID and telemetryName are in the body?

	requestBody, err := json.Marshal(extractedTelemetry)
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal telemetry %s", telemetryName)
	}

	resp, err := tp.httpClient.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return errors.Wrapf(err, "An error occurred while publishing telemetry %s", telemetryName)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return errors.Errorf("Unexpected response code %d while publishing telemetry %s", resp.StatusCode, telemetryName)
	}

	return nil
}

var telemetryServiceUrl = "https://telemetry.trento.suse.com"

func NewTelemetryPublisher() Publisher {
	return &TelemetryPublisher{
		apiHost:    telemetryServiceUrl,
		httpClient: &http.Client{},
	}
}
