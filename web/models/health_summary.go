package models

const (
	HealthSummaryHealthPassing  = "passing"
	HealthSummaryHealthWarning  = "warning"
	HealthSummaryHealthCritical = "critical"
	HealthSummaryHealthUnknown  = "unknown"
)

type HealthSummary []SAPSystemHealthSummary

type SAPSystemHealthSummary struct {
	ID              string `json:"id"`
	SID             string `json:"sid"`
	SAPSystemHealth string `json:"sapsystem_health"`
	ClustersHealth  string `json:"clusters_health"`
	DatabaseHealth  string `json:"database_health"`
	HostsHealth     string `json:"hosts_health"`
}
