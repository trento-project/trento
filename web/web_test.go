package web

import (
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/trento-project/trento/internal/grafana"
	"github.com/trento-project/trento/web/services"
)

func setupTestDependencies() Dependencies {
	return Dependencies{
		webEngine:               gin.Default(),
		collectorEngine:         gin.Default(),
		store:                   cookie.NewStore([]byte("secret")),
		settingsService:         newMockedSettingsService(),
		subscriptionsService:    newMockedSubscriptionsService(),
		premiumDetectionService: newMockedPremiumDetectionService(),
	}
}

func setupTestConfig() *Config {
	return &Config{
		Host: "",
		Port: 80,
		GrafanaConfig: &grafana.Config{
			PublicURL: "localhost",
			ApiURL:    "localhost",
			User:      "admin",
			Password:  "admin",
		},
	}
}

func newMockedSettingsService() services.SettingsService {
	settingsService := new(services.MockSettingsService)

	settingsService.On("InitializeIdentifier").Return(uuid.MustParse("59fd8017-b7fd-477b-9ebe-b658c558f3e9"), nil)
	settingsService.On("AcceptEula").Return(nil)
	settingsService.On("IsEulaAccepted").Return(true, nil)

	return settingsService
}

func newMockedSubscriptionsService() services.SubscriptionsService {
	subscriptionsService := new(services.MockSubscriptionsService)
	subscriptionsService.On("IsTrentoPremium").Return(true, nil)

	return subscriptionsService
}

func newMockedPremiumDetectionService() services.PremiumDetectionService {
	premiumDetection := new(services.MockPremiumDetectionService)
	premiumDetection.On("RequiresEulaAcceptance").Return(false, nil)

	return premiumDetection
}
