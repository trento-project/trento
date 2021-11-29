package models

type SAPSystem struct {
	ID        string
	SID       string
	Type      string
	Instances []*SAPSystemInstance
}

type SAPSystemInstance struct {
	InstanceNumber string
	Features       string
}
