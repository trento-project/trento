package services

import (
	log "github.com/sirupsen/logrus"
)

const (
	premium   = "Premium"
	community = "Community"
)

//go:generate mockery --name=PremiumDetectionService --inpackage --filename=premium_detection_mock.go

type PremiumDetectionService interface {
	RequiresEulaAcceptance() (bool, error)
	CanPublishTelemetry() (bool, error)
	IsPremiumActive() (bool, error)
}

type premiumDetectionService struct {
	flavor        string
	subscriptions SubscriptionsService
	settings      SettingsService
}

func NewPremiumDetectionService(flavor string, subscriptions SubscriptionsService, settings SettingsService) *premiumDetectionService {
	return &premiumDetectionService{
		flavor,
		subscriptions,
		settings,
	}
}

func (premiumDetection *premiumDetectionService) RequiresEulaAcceptance() (bool, error) {
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

func (premiumDetection *premiumDetectionService) CanPublishTelemetry() (bool, error) {
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

func (premiumDetection *premiumDetectionService) IsPremiumActive() (bool, error) {
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

func (premiumDetection *premiumDetectionService) isPremiumFlavor() bool {
	return premium == premiumDetection.flavor
}
