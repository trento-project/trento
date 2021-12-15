package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=SettingsService --inpackage --filename=settings_mock.go

type SettingsService interface {
	InitializeIdentifier() (uuid.UUID, error)
	IsEulaAccepted() (bool, error)
	AcceptEula() error
}

type settingsService struct {
	db *gorm.DB
}

func NewSettingsService(db *gorm.DB) SettingsService {
	return &settingsService{db: db}
}

func (s *settingsService) InitializeIdentifier() (uuid.UUID, error) {
	var settings entities.Settings
	err := s.db.First(&settings).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, err
	}
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

func (s *settingsService) IsEulaAccepted() (bool, error) {
	var settings entities.Settings
	err := s.db.First(&settings).Error
	if err != nil {
		return false, err
	}

	return settings.EulaAccepted, nil
}

func (s *settingsService) AcceptEula() error {
	var settings entities.Settings
	s.db.First(&settings)
	settings.EulaAccepted = true

	return s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "installation_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"eula_accepted"}),
	}).Create(&settings).Error
}
