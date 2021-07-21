package ara

import (
	"encoding/json"
	"fmt"
)

type RecordList struct {
	Count   int                 `json:"count,omitempty"`
	Results []*RecordListResult `json:"results,omitempty"`
}

type RecordListResult struct {
	ID       int    `json:"id,omitempty"`
	Playbook int    `json:"playbook,omitempty"`
	Key      string `json:"key,omitempty"`
	Type     string `json:"type,omitempty"`
}

type Record struct {
	ID    int         `json:"id,omitempty"`
	Value interface{} `json:"value,omitempty"`
	Key   string      `json:"key,omitempty"`
	Type  string      `json:"type,omitempty"`
}

func (a *araService) GetRecordList(filter string) (*RecordList, error) {
	rList := &RecordList{}

	var err error
	resp, err := getJson(a.composeQuery("records", filter))
	if err != nil {
		return rList, err
	}

	err = json.Unmarshal(resp, rList)
	if err != nil {
		return rList, err
	}

	return rList, nil
}

func (a *araService) GetRecord(recordId int) (*Record, error) {
	r := &Record{}

	var err error

	resp, err := getJson(a.composeQuery(fmt.Sprintf("records/%d", recordId), ""))
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(resp, r)
	if err != nil {
		return r, err
	}

	return r, nil
}
