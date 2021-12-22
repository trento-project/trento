package models

const (
	CheckPassing   string = "passing"
	CheckWarning   string = "warning"
	CheckCritical  string = "critical"
	CheckSkipped   string = "skipped"
	CheckUndefined string = "undefined"
)

type ChecksResult struct {
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
