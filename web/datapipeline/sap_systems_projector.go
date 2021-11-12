package datapipeline

import (
	"fmt"

	"gorm.io/gorm"
)

func NewSAPSystemsProjector(db *gorm.DB) *projector {
	SAPSystemsProjector := NewProjector("hosts", db)

	SAPSystemsProjector.AddHandler(SAPsystemDiscovery, SAPSystemsProjector_SAPSystemsDiscoveryHandler)

	return SAPSystemsProjector
}

func SAPSystemsProjector_SAPSystemsDiscoveryHandler(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
	fmt.Println(dataCollectedEvent)
	return nil
}
