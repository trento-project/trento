package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	araMocks "github.com/trento-project/trento/web/services/ara/mocks"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services/ara"
)

func TestAggregatedCheckDataString(t *testing.T) {
	aCritical := &AggregatedCheckData{
		PassingCount:  2,
		WarningCount:  1,
		CriticalCount: 1,
	}

	assert.Equal(t, aCritical.String(), "critical")

	aWarning := &AggregatedCheckData{
		PassingCount:  2,
		WarningCount:  1,
		CriticalCount: 0,
	}

	assert.Equal(t, aWarning.String(), "warning")

	aPassing := &AggregatedCheckData{
		PassingCount:  2,
		WarningCount:  0,
		CriticalCount: 0,
	}

	assert.Equal(t, aPassing.String(), "passing")

	aUndefined := &AggregatedCheckData{
		PassingCount:  0,
		WarningCount:  0,
		CriticalCount: 0,
	}

	assert.Equal(t, aUndefined.String(), "undefined")

}

func araResultRecord() *ara.Record {
	return &ara.Record{
		ID: 1,
		Value: map[string]interface{}{
			"results": map[string]interface{}{
				"myClusterId": map[string]interface{}{
					"checks": map[string]interface{}{
						"1.1.1": map[string]interface{}{
							"hosts": map[string]interface{}{
								"host1": map[string]interface{}{
									"result": "passing",
								},
								"host2": map[string]interface{}{
									"result": "passing",
								},
							},
						},
						"1.1.2": map[string]interface{}{
							"hosts": map[string]interface{}{
								"host1": map[string]interface{}{
									"result": "warning",
								},
								"host2": map[string]interface{}{
									"result": "critical",
								},
							},
						},
						"1.1.3": map[string]interface{}{
							"hosts": map[string]interface{}{
								"host1": map[string]interface{}{
									"result": "passing",
								},
								"host2": map[string]interface{}{
									"result": "warning",
								},
							},
						},
						"1.1.4": map[string]interface{}{
							"hosts": map[string]interface{}{
								"host1": map[string]interface{}{
									"result": "skipped",
								},
								"host2": map[string]interface{}{
									"result": "skipped",
								},
							},
						},
					},
				},
			},
		},
		Key:  "results",
		Type: "json",
	}
}

func TestGetChecksCatalog(t *testing.T) {

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{
		Count: 3,
		Results: []*ara.RecordListResult{
			&ara.RecordListResult{
				ID:       3,
				Playbook: 1,
				Key:      "metadata",
				Type:     "json",
			},
			&ara.RecordListResult{
				ID:       2,
				Playbook: 1,
				Key:      "metadata",
				Type:     "json",
			},
			&ara.RecordListResult{
				ID:       1,
				Playbook: 1,
				Key:      "metadata",
				Type:     "json",
			},
		},
	}

	mockAra.On("GetRecordList", "key=trento-metadata&order=-id").Return(
		rList, nil,
	)

	r := &ara.Record{
		ID: 1,
		Value: map[string]interface{}{
			"metadata": map[string]interface{}{
				"checks": []interface{}{
					map[string]interface{}{
						"id":             "ABCDEF",
						"name":           "1.1.1",
						"group":          "group 1",
						"description":    "description 1",
						"remediation":    "remediation 1",
						"labels":         "labels 1",
						"implementation": "implementation 1",
					},
					map[string]interface{}{
						"id":             "123456",
						"name":           "1.1.2",
						"group":          "group 2",
						"description":    "description 2",
						"remediation":    "remediation 2",
						"labels":         "labels 2",
						"implementation": "implementation 2",
					},
				},
			},
		},
		Key:  "metadata",
		Type: "json",
	}

	mockAra.On("GetRecord", 3).Return(
		r, nil,
	)

	checksService := NewChecksService(mockAra)
	c, err := checksService.GetChecksCatalog()

	expectedChecks := models.CheckList{
		&models.Check{
			ID:             "ABCDEF",
			Name:           "1.1.1",
			Group:          "group 1",
			Description:    "description 1",
			Remediation:    "remediation 1",
			Implementation: "implementation 1",
			Labels:         "labels 1",
		},
		&models.Check{
			ID:             "123456",
			Name:           "1.1.2",
			Group:          "group 2",
			Description:    "description 2",
			Remediation:    "remediation 2",
			Implementation: "implementation 2",
			Labels:         "labels 2",
		},
	}

	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedChecks, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksCatalogEmpty(t *testing.T) {

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{}

	mockAra.On("GetRecordList", "key=trento-metadata&order=-id").Return(
		rList, nil,
	)

	checksService := NewChecksService(mockAra)
	c, err := checksService.GetChecksCatalog()

	expectedChecks := models.CheckList(nil)

	assert.EqualError(t, err, "Couldn't find any check catalog record. Check if the runner component is running")
	assert.Equal(t, expectedChecks, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksCatalogListError(t *testing.T) {

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{}

	mockAra.On("GetRecordList", "key=trento-metadata&order=-id").Return(
		rList, fmt.Errorf("Some error"),
	)

	checksService := NewChecksService(mockAra)
	c, err := checksService.GetChecksCatalog()

	expectedChecks := models.CheckList(nil)

	assert.EqualError(t, err, "Some error")
	assert.Equal(t, expectedChecks, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksCatalogRecordError(t *testing.T) {

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{
		Count: 3,
		Results: []*ara.RecordListResult{
			&ara.RecordListResult{
				ID:       1,
				Playbook: 1,
				Key:      "metadata",
				Type:     "json",
			},
		},
	}

	mockAra.On("GetRecordList", "key=trento-metadata&order=-id").Return(
		rList, nil,
	)

	r := &ara.Record{}

	mockAra.On("GetRecord", 1).Return(
		r, fmt.Errorf("Some other error"),
	)

	checksService := NewChecksService(mockAra)
	c, err := checksService.GetChecksCatalog()

	expectedChecks := models.CheckList(nil)

	assert.EqualError(t, err, "Some other error")
	assert.Equal(t, expectedChecks, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksCatalogByGroup(t *testing.T) {

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{
		Count: 3,
		Results: []*ara.RecordListResult{
			&ara.RecordListResult{
				ID:       3,
				Playbook: 1,
				Key:      "metadata",
				Type:     "json",
			},
			&ara.RecordListResult{
				ID:       2,
				Playbook: 1,
				Key:      "metadata",
				Type:     "json",
			},
			&ara.RecordListResult{
				ID:       1,
				Playbook: 1,
				Key:      "metadata",
				Type:     "json",
			},
		},
	}

	mockAra.On("GetRecordList", "key=trento-metadata&order=-id").Return(
		rList, nil,
	)

	r := &ara.Record{
		ID: 1,
		Value: map[string]interface{}{
			"metadata": map[string]interface{}{
				"checks": []interface{}{
					map[string]interface{}{
						"id":             "ABCDEF",
						"name":           "1.1.1",
						"group":          "group 1",
						"description":    "description 1",
						"remediation":    "remediation 1",
						"labels":         "labels 1",
						"implementation": "implementation 1",
					},
					map[string]interface{}{
						"id":             "123456",
						"name":           "1.1.2",
						"group":          "group 1",
						"description":    "description 2",
						"remediation":    "remediation 2",
						"labels":         "labels 2",
						"implementation": "implementation 2",
					},
					map[string]interface{}{
						"id":             "123ABC",
						"name":           "1.2.1",
						"group":          "group 2",
						"description":    "description 3",
						"remediation":    "remediation 3",
						"labels":         "labels 3",
						"implementation": "implementation 3",
					},
				},
			},
		},
		Key:  "metadata",
		Type: "json",
	}

	mockAra.On("GetRecord", 3).Return(
		r, nil,
	)

	checksService := NewChecksService(mockAra)
	c, err := checksService.GetChecksCatalogByGroup()

	expectedChecks := models.GroupedCheckList{
		&models.GroupedChecks{
			Group: "group 1",
			Checks: models.CheckList{
				&models.Check{
					ID:             "ABCDEF",
					Name:           "1.1.1",
					Group:          "group 1",
					Description:    "description 1",
					Remediation:    "remediation 1",
					Implementation: "implementation 1",
					Labels:         "labels 1",
				},
				&models.Check{
					ID:             "123456",
					Name:           "1.1.2",
					Group:          "group 1",
					Description:    "description 2",
					Remediation:    "remediation 2",
					Implementation: "implementation 2",
					Labels:         "labels 2",
				},
			},
		},
		&models.GroupedChecks{
			Group: "group 2",
			Checks: models.CheckList{
				&models.Check{
					ID:             "123ABC",
					Name:           "1.2.1",
					Group:          "group 2",
					Description:    "description 3",
					Remediation:    "remediation 3",
					Implementation: "implementation 3",
					Labels:         "labels 3",
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedChecks, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksResult(t *testing.T) {

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{
		Count: 3,
		Results: []*ara.RecordListResult{
			&ara.RecordListResult{
				ID:       3,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
			&ara.RecordListResult{
				ID:       2,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
			&ara.RecordListResult{
				ID:       1,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
		},
	}

	mockAra.On("GetRecordList", "key=trento-results&order=-id").Return(
		rList, nil,
	)

	mockAra.On("GetRecord", 3).Return(
		araResultRecord(), nil,
	)

	checksService := NewChecksService(mockAra)
	c, err := checksService.GetChecksResult()

	expectedResults := map[string]*models.Results{
		"myClusterId": &models.Results{
			Checks: map[string]*models.ChecksByHost{
				"1.1.1": &models.ChecksByHost{
					Hosts: map[string]*models.Check{
						"host1": &models.Check{
							Result: models.CheckPassing,
						},
						"host2": &models.Check{
							Result: models.CheckPassing,
						},
					},
				},
				"1.1.2": &models.ChecksByHost{
					Hosts: map[string]*models.Check{
						"host1": &models.Check{
							Result: models.CheckWarning,
						},
						"host2": &models.Check{
							Result: models.CheckCritical,
						},
					},
				},
				"1.1.3": &models.ChecksByHost{
					Hosts: map[string]*models.Check{
						"host1": &models.Check{
							Result: models.CheckPassing,
						},
						"host2": &models.Check{
							Result: models.CheckWarning,
						},
					},
				},
				"1.1.4": &models.ChecksByHost{
					Hosts: map[string]*models.Check{
						"host1": &models.Check{
							Result: models.CheckSkipped,
						},
						"host2": &models.Check{
							Result: models.CheckSkipped,
						},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedResults, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksResultEmpty(t *testing.T) {

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{}

	mockAra.On("GetRecordList", "key=trento-results&order=-id").Return(
		rList, nil,
	)

	checksService := NewChecksService(mockAra)
	c, err := checksService.GetChecksResult()

	expectedResults := map[string]*models.Results(nil)

	assert.EqualError(t, err, "Couldn't find any check result record. Check if the runner component is running")
	assert.Equal(t, expectedResults, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksResultListError(t *testing.T) {

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{}

	mockAra.On("GetRecordList", "key=trento-results&order=-id").Return(
		rList, fmt.Errorf("Some error"),
	)

	checksService := NewChecksService(mockAra)
	c, err := checksService.GetChecksResult()

	expectedResults := map[string]*models.Results(nil)

	assert.EqualError(t, err, "Some error")
	assert.Equal(t, expectedResults, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksResultRecordError(t *testing.T) {

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{
		Count: 3,
		Results: []*ara.RecordListResult{
			&ara.RecordListResult{
				ID:       1,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
		},
	}

	mockAra.On("GetRecordList", "key=trento-results&order=-id").Return(
		rList, nil,
	)

	r := &ara.Record{}

	mockAra.On("GetRecord", 1).Return(
		r, fmt.Errorf("Some other error"),
	)

	checksService := NewChecksService(mockAra)
	c, err := checksService.GetChecksResult()

	expectedResults := map[string]*models.Results(nil)

	assert.EqualError(t, err, "Some other error")
	assert.Equal(t, expectedResults, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksResultByCluster(t *testing.T) {

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{
		Count: 3,
		Results: []*ara.RecordListResult{
			&ara.RecordListResult{
				ID:       3,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
			&ara.RecordListResult{
				ID:       2,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
			&ara.RecordListResult{
				ID:       1,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
		},
	}

	mockAra.On("GetRecordList", "key=trento-results&order=-id").Return(
		rList, nil,
	)

	mockAra.On("GetRecord", 3).Return(
		araResultRecord(), nil,
	)

	checksService := NewChecksService(mockAra)
	c, err := checksService.GetChecksResultByCluster("myClusterId")

	expectedResults := &models.Results{
		Checks: map[string]*models.ChecksByHost{
			"1.1.1": &models.ChecksByHost{
				Hosts: map[string]*models.Check{
					"host1": &models.Check{
						Result: models.CheckPassing,
					},
					"host2": &models.Check{
						Result: models.CheckPassing,
					},
				},
			},
			"1.1.2": &models.ChecksByHost{
				Hosts: map[string]*models.Check{
					"host1": &models.Check{
						Result: models.CheckWarning,
					},
					"host2": &models.Check{
						Result: models.CheckCritical,
					},
				},
			},
			"1.1.3": &models.ChecksByHost{
				Hosts: map[string]*models.Check{
					"host1": &models.Check{
						Result: models.CheckPassing,
					},
					"host2": &models.Check{
						Result: models.CheckWarning,
					},
				},
			},
			"1.1.4": &models.ChecksByHost{
				Hosts: map[string]*models.Check{
					"host1": &models.Check{
						Result: models.CheckSkipped,
					},
					"host2": &models.Check{
						Result: models.CheckSkipped,
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedResults, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksResultAndMetadataByCluster(t *testing.T) {
	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{
		Count: 1,
		Results: []*ara.RecordListResult{
			&ara.RecordListResult{
				ID:       1,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
		},
	}

	rMetaList := &ara.RecordList{
		Count: 1,
		Results: []*ara.RecordListResult{
			&ara.RecordListResult{
				ID:       2,
				Playbook: 1,
				Key:      "metadata",
				Type:     "json",
			},
		},
	}

	mockAra.On("GetRecordList", "key=trento-results&order=-id").Return(rList, nil)
	mockAra.On("GetRecordList", "key=trento-metadata&order=-id").Return(rMetaList, nil)

	araResultRecord := &ara.Record{
		ID: 1,
		Value: map[string]interface{}{
			"results": map[string]interface{}{
				"myClusterId": map[string]interface{}{
					"hosts": map[string]interface{}{
						"host1": map[string]interface{}{
							"reachable": true,
							"msg":       "",
						},
						"host2": map[string]interface{}{
							"reachable": false,
							"msg":       "error connecting",
						},
					},
					"checks": map[string]interface{}{
						"ABCDEF": map[string]interface{}{
							"hosts": map[string]interface{}{
								"host1": map[string]interface{}{
									"result": "passing",
								},
								"host2": map[string]interface{}{
									"result": "passing",
								},
							},
						},
						"123456": map[string]interface{}{
							"hosts": map[string]interface{}{
								"host1": map[string]interface{}{
									"result": "warning",
								},
								"host2": map[string]interface{}{
									"result": "critical",
								},
							},
						},
					},
				},
			},
		},
	}

	araMetaRecord := &ara.Record{
		ID: 2,
		Value: map[string]interface{}{
			"metadata": map[string]interface{}{
				"checks": []interface{}{
					map[string]interface{}{
						"id":             "ABCDEF",
						"name":           "1.1.1",
						"group":          "group 1",
						"description":    "description 1",
						"remediation":    "remediation 1",
						"labels":         "labels 1",
						"implementation": "implementation 1",
					},
					map[string]interface{}{
						"id":             "123456",
						"name":           "1.1.2",
						"group":          "group 1",
						"description":    "description 2",
						"remediation":    "remediation 2",
						"labels":         "labels 2",
						"implementation": "implementation 2",
					},
				},
			},
		},
		Key:  "metadata",
		Type: "json",
	}

	mockAra.On("GetRecord", 1).Return(araResultRecord, nil)
	mockAra.On("GetRecord", 2).Return(araMetaRecord, nil)

	checksService := NewChecksService(mockAra)
	results, err := checksService.GetChecksResultAndMetadataByCluster("myClusterId")

	expectedResults := &models.ClusterCheckResults{
		Hosts: map[string]*models.Host{
			"host1": &models.Host{
				Reachable: true,
				Msg:       "",
			},
			"host2": &models.Host{
				Reachable: false,
				Msg:       "error connecting",
			},
		},
		Checks: []models.ClusterCheckResult{
			models.ClusterCheckResult{
				ID:          "ABCDEF",
				Group:       "group 1",
				Description: "description 1",
				Hosts: map[string]*models.Check{
					"host1": &models.Check{
						Result: models.CheckPassing,
					},
					"host2": &models.Check{
						Result: models.CheckPassing,
					},
				},
			},
			models.ClusterCheckResult{
				ID:          "123456",
				Group:       "group 1",
				Description: "description 2",
				Hosts: map[string]*models.Check{
					"host1": &models.Check{
						Result: models.CheckWarning,
					},
					"host2": &models.Check{
						Result: models.CheckCritical,
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedResults, results)
}

func TestGetAggregatedChecksResultByHost(t *testing.T) {

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{
		Count: 3,
		Results: []*ara.RecordListResult{
			&ara.RecordListResult{
				ID:       3,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
			&ara.RecordListResult{
				ID:       2,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
			&ara.RecordListResult{
				ID:       1,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
		},
	}

	mockAra.On("GetRecordList", "key=trento-results&order=-id").Return(
		rList, nil,
	)

	mockAra.On("GetRecord", 3).Return(
		araResultRecord(), nil,
	)

	checksService := NewChecksService(mockAra)
	c, err := checksService.GetAggregatedChecksResultByHost("myClusterId")

	expectedResults := map[string]*AggregatedCheckData{
		"host1": &AggregatedCheckData{
			PassingCount:  2,
			WarningCount:  1,
			CriticalCount: 0,
		},
		"host2": &AggregatedCheckData{
			PassingCount:  1,
			WarningCount:  1,
			CriticalCount: 1,
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedResults, c)

	mockAra.AssertExpectations(t)
}

func TestGetAggregatedChecksResultByCluster(t *testing.T) {

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{
		Count: 3,
		Results: []*ara.RecordListResult{
			&ara.RecordListResult{
				ID:       3,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
			&ara.RecordListResult{
				ID:       2,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
			&ara.RecordListResult{
				ID:       1,
				Playbook: 1,
				Key:      "results",
				Type:     "json",
			},
		},
	}

	mockAra.On("GetRecordList", "key=trento-results&order=-id").Return(
		rList, nil,
	)

	mockAra.On("GetRecord", 3).Return(
		araResultRecord(), nil,
	)

	checksService := NewChecksService(mockAra)
	c, err := checksService.GetAggregatedChecksResultByCluster("myClusterId")

	expectedResults := &AggregatedCheckData{
		PassingCount:  3,
		WarningCount:  2,
		CriticalCount: 1,
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedResults, c)

	mockAra.AssertExpectations(t)
}
