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
		Value: []interface{}{
			map[string]interface{}{
				"id":             "1.1.1",
				"name":           "check 1",
				"description":    "description 1",
				"remediation":    "remediation 1",
				"labels":         "labels 1",
				"implementation": "implementation 1",
			},
			map[string]interface{}{
				"id":             "1.1.2",
				"name":           "check 2",
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

	expectedChecks := []*models.Check{
		&models.Check{
			ID:             "1.1.1",
			Name:           "check 1",
			Description:    "description 1",
			Remediation:    "remediation 1",
			Implementation: "implementation 1",
			Labels:         "labels 1",
		},
		&models.Check{
			ID:             "1.1.2",
			Name:           "check 2",
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

	expectedChecks := []*models.Check{}

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

	expectedChecks := []*models.Check{}

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

	expectedChecks := []*models.Check{}

	assert.EqualError(t, err, "Some other error")
	assert.Equal(t, expectedChecks, c)

	mockAra.AssertExpectations(t)
}
