package entities

type Settings struct {
	InstallationID string `gorm:"primaryKey"`
	EulaAccepted   bool
}
