package services

import (
	log "github.com/sirupsen/logrus"
)

const (
	Premium   = "Premium"
	Community = "Community"
)

//go:generate mockery --name=PremiumDetection --inpackage --filename=premium_detection_mock.go

type PremiumDetectionService interface {
	RequiresEulaAcceptance() (bool, error)
	CanPublishTelemetry() (bool, error)
	IsPremiumActive() (bool, error)
}

type premiumDetection struct {
	flavor        string
	subscriptions SubscriptionsService
	settings      SettingsService
}

func NewPremiumDetection(flavor string, subscriptions SubscriptionsService, settings SettingsService) *premiumDetection {
	return &premiumDetection{
		flavor,
		subscriptions,
		settings,
	}
}

func (premiumDetection *premiumDetection) RequiresEulaAcceptance() (bool, error) {
	if !premiumDetection.isPremiumFlavor() {
		return false, nil
	}
	isEulaAccepted, err := premiumDetection.settings.IsEulaAccepted()
	if err != nil {
		log.Errorf("Unable to determine whether the EULA has been accepted. Error: %s", err)
		return false, err
	}
	return !isEulaAccepted, err
}

func (premiumDetection *premiumDetection) CanPublishTelemetry() (bool, error) {
	if !premiumDetection.isPremiumFlavor() {
		return false, nil
	}
	isEulaAccepted, err := premiumDetection.settings.IsEulaAccepted()
	if err != nil {
		log.Errorf("Unable to determine whether telemetry can be published. Error: %s", err)
		return false, err
	}
	return isEulaAccepted, nil
}

func (premiumDetection *premiumDetection) IsPremiumActive() (bool, error) {
	if !premiumDetection.isPremiumFlavor() {
		return false, nil
	}
	isPremiumActive, err := premiumDetection.subscriptions.IsTrentoPremium()
	if err != nil {
		log.Errorf("Unable to determine whether the Trento Premium installation is active. Error: %s", err)
		return false, err
	}
	return isPremiumActive, nil
}

func (premiumDetection *premiumDetection) isPremiumFlavor() bool {
	return Premium == premiumDetection.flavor
}
