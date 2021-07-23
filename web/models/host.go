package models

type Host struct {
	Name          string `gorm:"primaryKey"`
	Address       string
	Health        string
	CloudProvider string
	Cluster       string
	SAPSystem     string
	Landscape     string
	Environment   string
}
