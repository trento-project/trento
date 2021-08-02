package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	araMocks "github.com/trento-project/trento/web/services/ara/mocks"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services/ara"
)

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
			"1.1.1": map[string]interface{}{
				"id":             "1.1.1",
				"name":           "check 1",
				"group":          "group 1",
				"description":    "description 1",
				"remediation":    "remediation 1",
				"labels":         "labels 1",
				"implementation": "implementation 1",
			},
			"1.1.2": map[string]interface{}{
				"id":             "1.1.2",
				"name":           "check 2",
				"group":          "group 2",
				"description":    "description 2",
				"remediation":    "remediation 2",
				"labels":         "labels 2",
				"implementation": "implementation 2",
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

	expectedChecks := map[string]*models.Check{
		"1.1.1": &models.Check{
			ID:             "1.1.1",
			Name:           "check 1",
			Group:          "group 1",
			Description:    "description 1",
			Remediation:    "remediation 1",
			Implementation: "implementation 1",
			Labels:         "labels 1",
		},
		"1.1.2": &models.Check{
			ID:             "1.1.2",
			Name:           "check 2",
			Group:          "group 2",
			Description:    "description 2",
			Remediation:    "remediation 2",
			Implementation: "implementation 2",
			Labels:         "labels 2",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedChecks, c)

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

	expectedChecks := map[string]*models.Check{}

	assert.NoError(t, err)
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

	expectedChecks := map[string]*models.Check{}

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

	expectedChecks := map[string]*models.Check{}

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
			"1.1.1": map[string]interface{}{
				"id":             "1.1.1",
				"name":           "check 1",
				"group":          "group 1",
				"description":    "description 1",
				"remediation":    "remediation 1",
				"labels":         "labels 1",
				"implementation": "implementation 1",
			},
			"1.1.2": map[string]interface{}{
				"id":             "1.1.2",
				"name":           "check 2",
				"group":          "group 1",
				"description":    "description 2",
				"remediation":    "remediation 2",
				"labels":         "labels 2",
				"implementation": "implementation 2",
			},
			"1.2.3": map[string]interface{}{
				"id":             "1.2.3",
				"name":           "check 3",
				"group":          "group 2",
				"description":    "description 3",
				"remediation":    "remediation 3",
				"labels":         "labels 3",
				"implementation": "implementation 3",
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

	expectedChecks := map[string]map[string]*models.Check{
		"1.1 - group 1": {
			"1.1.1": &models.Check{
				ID:             "1.1.1",
				Name:           "check 1",
				Group:          "group 1",
				Description:    "description 1",
				Remediation:    "remediation 1",
				Implementation: "implementation 1",
				Labels:         "labels 1",
			},
			"1.1.2": &models.Check{
				ID:             "1.1.2",
				Name:           "check 2",
				Group:          "group 1",
				Description:    "description 2",
				Remediation:    "remediation 2",
				Implementation: "implementation 2",
				Labels:         "labels 2",
			},
		},
		"1.2 - group 2": {
			"1.2.3": &models.Check{
				ID:             "1.2.3",
				Name:           "check 3",
				Group:          "group 2",
				Description:    "description 3",
				Remediation:    "remediation 3",
				Implementation: "implementation 3",
				Labels:         "labels 3",
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedChecks, c)

	mockAra.AssertExpectations(t)

}
