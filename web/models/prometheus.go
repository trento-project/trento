package models

type PrometheusTargetsList []*PrometheusTargets

type PrometheusTargets struct {
	Targets []string
	Labels  map[string]string
}
