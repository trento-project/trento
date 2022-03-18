package models

import (
	"github.com/trento-project/trento/internal/sapsystem/sapcontrol"
)

const (
	SAPSystemTypeApplication = "application"
	SAPSystemTypeDatabase    = "database"

	SAPSystemHealthPassing  = "passing"
	SAPSystemHealthWarning  = "warning"
	SAPSystemHealthCritical = "critical"
	SAPSystemHealthUnknown  = "unknown"
)

type SAPSystem struct {
	ID               string
	SID              string
	Type             string
	Instances        []*SAPSystemInstance
	AttachedDatabase *SAPSystem
	DBName           string
	DBHost           string
	DBAddress        string
	Health           string
	Tags             []string
	// TODO: this is frontend specific, should be removed
	HasDuplicatedSID bool
}

type SAPSystemInstance struct {
	Type                    string
	SID                     string
	Features                string
	InstanceNumber          string
	SystemReplication       string
	SystemReplicationStatus string
	SAPHostname             string
	Status                  string
	StartPriority           string
	HttpPort                int
	HttpsPort               int
	ClusterName             string
	ClusterID               string
	ClusterType             string
	HostID                  string
	Hostname                string
}

type SAPSystemList []*SAPSystem

func (s SAPSystem) GetAllInstances() []*SAPSystemInstance {
	instances := s.Instances

	if s.AttachedDatabase != nil {
		instances = append(instances, s.AttachedDatabase.Instances...)
	}

	return instances
}

func (s SAPSystemInstance) Health() string {
	switch s.Status {
	case string(sapcontrol.STATECOLOR_RED):
		return SAPSystemHealthCritical
	case string(sapcontrol.STATECOLOR_YELLOW):
		return SAPSystemHealthWarning
	case string(sapcontrol.STATECOLOR_GREEN):
		return SAPSystemHealthPassing
	default:
		return SAPSystemHealthUnknown
	}
}
