package web

type HealthContainer struct {
	PassingCount  int
	WarningCount  int
	CriticalCount int
	Layout        string
}
