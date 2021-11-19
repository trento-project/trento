package web

import (
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/trento-project/trento/web/services"
)

func setupTestDependencies() Dependencies {
	return Dependencies{
		webEngine:       gin.Default(),
		collectorEngine: gin.Default(),
		store:           cookie.NewStore([]byte("secret")),
		settingsService: newMockedSettingsService(),
	}
}

func setupTestConfig() *Config {
	return &Config{
		Host: "",
		Port: 80,
	}
}

func newMockedSettingsService() services.SettingsService {
	settingsService := new(services.MockSettingsService)

	settingsService.On("InitializeIdentifier").Return(uuid.MustParse("59fd8017-b7fd-477b-9ebe-b658c558f3e9"), nil)

	return settingsService
}
