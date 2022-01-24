package models

const (
	CheckPassing   string = "passing"
	CheckWarning   string = "warning"
	CheckCritical  string = "critical"
	CheckSkipped   string = "skipped"
	CheckUndefined string = "undefined"
)

type ChecksResult struct {
	ID     string                   `json:"-"`
	Hosts  map[string]*HostState    `json:"hosts,omitempty"`
	Checks map[string]*ChecksByHost `json:"checks,omitempty"`
}

// Simplifed stuct consumed by the frontend
type ChecksResultAsList struct {
	Hosts  map[string]*HostState `json:"hosts,omitempty"`
	Checks []*ChecksByHost       `json:"checks,omitempty"`
}

// The ChecksByHost struct stores the checks list, but the results are grouped by hosts
type ChecksByHost struct {
	Hosts       map[string]*Check `json:"hosts,omitempty"`
	ID          string            `json:"id,omitempty"`
	Group       string            `json:"group,omitempty"`
	Description string            `json:"description,omitempty"`
}

type HostState struct {
	Reachable bool   `json:"reachable"`
	Msg       string `json:"msg"`
}

type AggregatedCheckData struct {
	PassingCount  int
	WarningCount  int
	CriticalCount int
}

func (c *ChecksResult) GetAggregatedChecksResultByHost() map[string]*AggregatedCheckData {
	aCheckDataByHost := make(map[string]*AggregatedCheckData)

	for _, check := range c.Checks {
		for hostName, host := range check.Hosts {
			if _, ok := aCheckDataByHost[hostName]; !ok {
				aCheckDataByHost[hostName] = &AggregatedCheckData{}
			}
			switch host.Result {
			case CheckCritical:
				aCheckDataByHost[hostName].CriticalCount += 1
			case CheckWarning:
				aCheckDataByHost[hostName].WarningCount += 1
			case CheckPassing:
				aCheckDataByHost[hostName].PassingCount += 1
			}
		}
	}

	return aCheckDataByHost
}

func (c *ChecksResult) GetAggregatedChecksResultByCluster() *AggregatedCheckData {
	aCheckData := &AggregatedCheckData{}
	aCheckDataByHost := c.GetAggregatedChecksResultByHost()

	for _, aData := range aCheckDataByHost {
		aCheckData.CriticalCount += aData.CriticalCount
		aCheckData.WarningCount += aData.WarningCount
		aCheckData.PassingCount += aData.PassingCount
	}

	return aCheckData
}

func (a *AggregatedCheckData) String() string {
	if a.CriticalCount > 0 {
		return CheckCritical
	} else if a.WarningCount > 0 {
		return CheckWarning
	} else if a.PassingCount > 0 {
		return CheckPassing
	}

	return CheckUndefined
}
