package models

const (
	SAPSystemTypeApplication = "application"
	SAPSystemTypeDatabase    = "database"
)

type SAPSystem struct {
	ID               string
	SID              string
	Type             string
	Instances        []*SAPSystemInstance
	AttachedDatabase *SAPSystem
	DBName           string
	DBHost           string
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
