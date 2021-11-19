package services

import (
	"github.com/google/uuid"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/gorm"
)

//go:generate mockery --name=SettingsService --inpackage --filename=settings_mock.go

type SettingsService interface {
	InitializeIdentifier() (uuid.UUID, error)
}

type settingsService struct {
	db *gorm.DB
}

func NewSettingsService(db *gorm.DB) SettingsService {
	return &settingsService{db: db}
}

func (s *settingsService) InitializeIdentifier() (uuid.UUID, error) {
	var settings entities.Settings
	s.db.First(&settings)
	if settings.InstallationID != "" {
		return uuid.MustParse(settings.InstallationID), nil
	}

	installationUUID := uuid.New()
	settings.InstallationID = installationUUID.String()
	if err := s.db.Create(&settings).Error; err != nil {
		return uuid.Nil, err
	}

	return installationUUID, nil
}
