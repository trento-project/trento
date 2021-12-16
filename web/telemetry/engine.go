package telemetry

import (
	"context"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/web/services"
)

var telemetryCollectionInterval = 24 * time.Hour

// Engine is the entrypoint for the telemetry extraction and publishing system.
type Engine struct {
	installationID    uuid.UUID
	publisher         Publisher
	telemetryRegistry *TelemetryRegistry
	premiumDetection  services.PremiumDetectionService
}

//go:generate mockery --name=Extractor --inpackage --filename=extractor_mock.go

// Extractor extracts telemetry data to be published.
type Extractor interface {
	Extract() (interface{}, error)
}

//go:generate mockery --name=InstallationIdAwareExtractor --inpackage --filename=installation_id_aware_extractor_mock.go

// InstallationIdAwareExtractor is an Extractor that can be identified by an installation ID.
type InstallationIdAwareExtractor interface {
	Extractor
	WithInstallationID(uuid.UUID)
}

// TelemetryRegistry is a map of enabled/supported extractors.
type TelemetryRegistry map[string]Extractor

//go:generate mockery --name=Publisher --inpackage --filename=publisher_mock.go

// Publisher publishes the extracted telemetry data to a collection service.
type Publisher interface {
	Publish(telemetryName string, installationID uuid.UUID, extractedTelemetry interface{}) error
}

func (e *Engine) Start(ctx context.Context) {
	log.Infof("Starting Telemetry Engine")

	canPublishTelemetry, err := e.premiumDetection.CanPublishTelemetry()
	if err != nil {
		log.Errorf("Unable to start Telemetry Engine. Error: %s", err)
		return
	}
	if !canPublishTelemetry {
		log.Infof("Telemetry publishing is not supported by this installation")
		return
	}

	extractAndPublishFn := func() {
		for telemetryName, extractor := range *e.telemetryRegistry {
			if identifiedExtractor, ok := extractor.(InstallationIdAwareExtractor); ok {
				identifiedExtractor.WithInstallationID(e.installationID)
			}
			extractedTelemetry, err := extractor.Extract()
			if err != nil {
				log.Errorf("Error while extracting telemetry %s: %s", telemetryName, err)
				continue
			}
			if err := e.publisher.Publish(telemetryName, e.installationID, extractedTelemetry); err != nil {
				log.Errorf("Error while publishing telemetry %s: %s", telemetryName, err)
			}
		}
	}

	internal.Repeat(
		"telemetry.extraction_and_publishing",
		extractAndPublishFn,
		telemetryCollectionInterval,
		ctx,
	)
}

func NewEngine(
	installationID uuid.UUID,
	publisher Publisher,
	telemetries *TelemetryRegistry,
	premium services.PremiumDetectionService,
) *Engine {
	return &Engine{
		installationID:    installationID,
		publisher:         publisher,
		telemetryRegistry: telemetries,
		premiumDetection:  premium,
	}
}

// installationIdAwareExtractor is an Extractor knowledgable of Trento's installation ID.
// It can be embedded in other extractors to support this information.
type installationIdAwareExtractor struct {
	installationID uuid.UUID
}

func (ex *installationIdAwareExtractor) WithInstallationID(ID uuid.UUID) {
	ex.installationID = ID
}
