package ara

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	mocks "github.com/trento-project/trento/web/services/ara/mocksutils"
)

func TestGetRecordList(t *testing.T) {

	mockGetJson := new(mocks.GetJson)

	returnJson := []byte(`{"count":3,"next":null,"previous":null,"results":
    [{"id":3,"playbook":1,"created":"2021-07-20T14:21:28.281795+02:00","updated":
    "2021-07-20T14:21:28.281837+02:00","key":"metadata","type":"json"},
    {"id":2,"playbook":1,"created":"2021-07-20T14:18:44.991988+02:00","updated":
    "2021-07-20T14:18:44.992031+02:00","key":"metadata","type":"json"},
    {"id":1,"playbook":1,"created":"2021-07-20T14:17:02.110605+02:00",
    "updated":"2021-07-20T14:17:02.110644+02:00","key":"metadata","type":"json"}]}`)

	mockGetJson.On("Execute", "http://127.0.0.1:80/api/v1/records?my_filter").Return(
		returnJson, nil,
	)

	getJson = mockGetJson.Execute

	araService := NewAraService("127.0.0.1:80")
	rList, err := araService.GetRecordList("my_filter")

	expectedRecordList := &RecordList{
		Count: 3,
		Results: []*RecordListResult{
			&RecordListResult{
				ID:       3,
				Playbook: 1,
				Key:      "metadata",
				Type:     "json",
			},
			&RecordListResult{
				ID:       2,
				Playbook: 1,
				Key:      "metadata",
				Type:     "json",
			},
			&RecordListResult{
				ID:       1,
				Playbook: 1,
				Key:      "metadata",
				Type:     "json",
			},
		},
	}

	assert.Equal(t, expectedRecordList, rList)
	assert.NoError(t, err)

	mockGetJson.AssertExpectations(t)
}

func TestGetRecordListError(t *testing.T) {

	mockGetJson := new(mocks.GetJson)

	mockGetJson.On("Execute", "http://127.0.0.1:80/api/v1/records?my_filter").Return(
		[]byte(""), fmt.Errorf("Some error"),
	)

	getJson = mockGetJson.Execute

	araService := NewAraService("127.0.0.1:80")
	rList, err := araService.GetRecordList("my_filter")

	expectedRecordList := &RecordList{}

	assert.Equal(t, expectedRecordList, rList)
	assert.EqualError(t, err, "Some error")

	mockGetJson.AssertExpectations(t)
}

func TestGetRecord(t *testing.T) {

	mockGetJson := new(mocks.GetJson)

	returnJson := []byte(`{"id":1,"playbook":{"id":230,"items":{"plays":1,
    "tasks":7,"results":7,"hosts":1,"files":3,"records":1},"arguments":{},
    "labels":[{"id":1,"name":"check:False"},{"id":2,"name":"tags:all"},
    {"id":3,"name":"meta"}],"started":"2021-07-20T14:21:27.198249+02:00",
    "ended":"2021-07-20T14:21:28.587075+02:00","duration":"00:00:01.388826",
    "name":null,"ansible_version":"2.11.2","status":"completed","path":
    "/usr/etc/trento/ansible/meta.yml","controller":"xarbulu-monitoring"},
    "value":[{"id":"1.1.1","name":"This is my test name","description":
    "This is my test description","remediation":"remediation","labels":"generic",
    "implementation":"my impl"},
    {"id":"1.1.2","name":"This is my test name","description":
    "This is my test description","remediation":"This is my test remediation",
    "labels":"generic","implementation":"impl"}],"created":
    "2021-07-20T14:21:28.281795+02:00","updated":"2021-07-20T14:21:28.281837+02:00",
    "key":"metadata","type":"json"}`)

	mockGetJson.On("Execute", "http://127.0.0.1:80/api/v1/records/1?").Return(
		returnJson, nil,
	)

	getJson = mockGetJson.Execute

	araService := NewAraService("127.0.0.1:80")
	r, err := araService.GetRecord(1)

	expectedRecord := &Record{
		ID: 1,
		Value: []interface{}{
			map[string]interface{}{
				"id":             "1.1.1",
				"name":           "This is my test name",
				"description":    "This is my test description",
				"remediation":    "remediation",
				"labels":         "generic",
				"implementation": "my impl",
			},
			map[string]interface{}{
				"id":             "1.1.2",
				"name":           "This is my test name",
				"description":    "This is my test description",
				"remediation":    "This is my test remediation",
				"labels":         "generic",
				"implementation": "impl",
			},
		},
		Key:  "metadata",
		Type: "json",
	}

	assert.Equal(t, expectedRecord, r)
	assert.NoError(t, err)

	mockGetJson.AssertExpectations(t)
}

func TestGetRecordError(t *testing.T) {

	mockGetJson := new(mocks.GetJson)

	mockGetJson.On("Execute", "http://127.0.0.1:80/api/v1/records/1?").Return(
		[]byte(""), fmt.Errorf("Some error"),
	)

	getJson = mockGetJson.Execute

	araService := NewAraService("127.0.0.1:80")
	rList, err := araService.GetRecord(1)

	expectedRecord := &Record{}

	assert.Equal(t, expectedRecord, rList)
	assert.EqualError(t, err, "Some error")

	mockGetJson.AssertExpectations(t)
}
