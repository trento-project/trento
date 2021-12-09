package services

import (
	"encoding/json"
	"fmt"
	"testing"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	araMocks "github.com/trento-project/trento/web/services/ara/mocks"

	"github.com/trento-project/trento/web/entities"
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
									"msg":    "some random message",
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

func TestGetChecksResult(t *testing.T) {
	db := helpers.SetupTestDatabase(t)

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

	checksService := NewChecksService(mockAra, db)
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
							Msg:    "some random message",
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
	db := helpers.SetupTestDatabase(t)

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{}

	mockAra.On("GetRecordList", "key=trento-results&order=-id").Return(
		rList, nil,
	)

	checksService := NewChecksService(mockAra, db)
	c, err := checksService.GetChecksResult()

	expectedResults := map[string]*models.Results(nil)

	assert.EqualError(t, err, "Couldn't find any check result record. Check if the runner component is running")
	assert.Equal(t, expectedResults, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksResultListError(t *testing.T) {
	db := helpers.SetupTestDatabase(t)

	mockAra := new(araMocks.AraService)

	rList := &ara.RecordList{}

	mockAra.On("GetRecordList", "key=trento-results&order=-id").Return(
		rList, fmt.Errorf("Some error"),
	)

	checksService := NewChecksService(mockAra, db)
	c, err := checksService.GetChecksResult()

	expectedResults := map[string]*models.Results(nil)

	assert.EqualError(t, err, "Some error")
	assert.Equal(t, expectedResults, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksResultRecordError(t *testing.T) {
	db := helpers.SetupTestDatabase(t)

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

	checksService := NewChecksService(mockAra, db)
	c, err := checksService.GetChecksResult()

	expectedResults := map[string]*models.Results(nil)

	assert.EqualError(t, err, "Some other error")
	assert.Equal(t, expectedResults, c)

	mockAra.AssertExpectations(t)
}

func TestGetChecksResultByCluster(t *testing.T) {
	db := helpers.SetupTestDatabase(t)

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

	checksService := NewChecksService(mockAra, db)
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
						Msg:    "some random message",
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

func TestGetAggregatedChecksResultByHost(t *testing.T) {
	db := helpers.SetupTestDatabase(t)

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

	checksService := NewChecksService(mockAra, db)
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
	db := helpers.SetupTestDatabase(t)

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

	checksService := NewChecksService(mockAra, db)
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

type ChecksServiceTestSuite struct {
	suite.Suite
	db            *gorm.DB
	tx            *gorm.DB
	ara           *araMocks.AraService
	checksService ChecksService
}

func TestChecksServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ChecksServiceTestSuite))
}

func (suite *ChecksServiceTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())
	suite.ara = new(araMocks.AraService)

	suite.db.AutoMigrate(
		entities.Check{}, models.SelectedChecks{}, models.ConnectionSettings{})
	loadAraFixtures(suite.ara)
	loadChecksCatalogFixtures(suite.db)
	loadSelectedChecksFixtures(suite.db)
	loadConnectionSettingsFixtures(suite.db)
}

func (suite *ChecksServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(entities.Check{})
	suite.db.Migrator().DropTable(models.SelectedChecks{})
	suite.db.Migrator().DropTable(models.ConnectionSettings{})
}

func (suite *ChecksServiceTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.checksService = NewChecksService(suite.ara, suite.tx)
}

func (suite *ChecksServiceTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func loadAraFixtures(a *araMocks.AraService) {
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

	a.On("GetRecordList", "key=trento-results&order=-id").Return(rList, nil)

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
						"check1": map[string]interface{}{
							"hosts": map[string]interface{}{
								"host1": map[string]interface{}{
									"result": "passing",
								},
								"host2": map[string]interface{}{
									"result": "passing",
								},
							},
						},
						"check2": map[string]interface{}{
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

	a.On("GetRecord", 1).Return(araResultRecord, nil)
}

func loadChecksCatalogFixtures(db *gorm.DB) {
	check1payload := `{"id":"check1","name":"name1","group":"group1","description":"description1"}`
	db.Create(&entities.Check{
		ID:      "check1",
		Payload: datatypes.JSON([]byte(check1payload)),
	})
	check3payload := `{"id":"check3","name":"name3","group":"group2","description":"description3"}`
	db.Create(&entities.Check{
		ID:      "check3",
		Payload: datatypes.JSON([]byte(check3payload)),
	})
	check2payload := `{"id":"check2","name":"name2","group":"group1","description":"description2"}`
	db.Create(&entities.Check{
		ID:      "check2",
		Payload: datatypes.JSON([]byte(check2payload)),
	})
}

func loadSelectedChecksFixtures(db *gorm.DB) {
	db.Create(&models.SelectedChecks{
		ID:             "group1",
		SelectedChecks: []string{"ABCDEF", "123456"},
	})
	db.Create(&models.SelectedChecks{
		ID:             "group2",
		SelectedChecks: []string{"ABC123", "123ABC"},
	})
	db.Create(&models.SelectedChecks{
		ID:             "group3",
		SelectedChecks: []string{"DEF456", "456DEF"},
	})
}

func loadConnectionSettingsFixtures(db *gorm.DB) {
	db.Create(&models.ConnectionSettings{
		ID:   "group1",
		Node: "node1",
		User: "user1",
	})
	db.Create(&models.ConnectionSettings{
		ID:   "group1",
		Node: "node2",
		User: "user2",
	})
	db.Create(&models.ConnectionSettings{
		ID:   "group2",
		Node: "node3",
		User: "user3",
	})
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetChecksCatalog() {
	catalog, err := suite.checksService.GetChecksCatalog()
	expectedCatalog := models.ChecksCatalog{
		&models.Check{
			ID:          "check1",
			Name:        "name1",
			Group:       "group1",
			Description: "description1",
		},
		&models.Check{
			ID:          "check2",
			Name:        "name2",
			Group:       "group1",
			Description: "description2",
		},
		&models.Check{
			ID:          "check3",
			Name:        "name3",
			Group:       "group2",
			Description: "description3",
		},
	}

	suite.NoError(err)
	suite.ElementsMatch(expectedCatalog, catalog)
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetChecksCatalogByGroup() {
	catalog, err := suite.checksService.GetChecksCatalogByGroup()
	expectedCatalog := models.GroupedCheckList{
		&models.GroupedChecks{
			Group: "group1",
			Checks: models.ChecksCatalog{
				&models.Check{
					ID:          "check1",
					Name:        "name1",
					Group:       "group1",
					Description: "description1",
				},
				&models.Check{
					ID:          "check2",
					Name:        "name2",
					Group:       "group1",
					Description: "description2",
				},
			},
		},
		&models.GroupedChecks{
			Group: "group2",
			Checks: models.ChecksCatalog{
				&models.Check{
					ID:          "check3",
					Name:        "name3",
					Group:       "group2",
					Description: "description3",
				},
			},
		},
	}

	suite.NoError(err)
	suite.ElementsMatch(expectedCatalog, catalog)
}

func (suite *ChecksServiceTestSuite) TestChecksService_CreateChecksCatalogEntry() {
	check := &models.Check{
		ID:          "checkNew",
		Name:        "nameNew",
		Group:       "groupNew",
		Description: "descriptionNew",
	}
	err := suite.checksService.CreateChecksCatalogEntry(check)

	var newCheckEntity entities.Check
	var newCheck models.Check

	suite.tx.Where("id", "checkNew").First(&newCheckEntity)

	json.Unmarshal(newCheckEntity.Payload, &newCheck)

	suite.NoError(err)
	suite.Equal(check, &newCheck)

	// Check update works
	check = &models.Check{
		ID:          "checkNew",
		Name:        "nameNewUpdated",
		Group:       "groupNewUpdated",
		Description: "descriptionNewUpdated",
	}
	err = suite.checksService.CreateChecksCatalogEntry(check)

	suite.tx.Where("id", "checkNew").First(&newCheckEntity)

	json.Unmarshal(newCheckEntity.Payload, &newCheck)

	suite.NoError(err)
	suite.Equal(check, &newCheck)
}

func (suite *ChecksServiceTestSuite) TestChecksService_CreateChecksCatalog() {
	check1 := &models.Check{
		ID:          "checkOther",
		Name:        "nameNew",
		Group:       "groupNew",
		Description: "descriptionNew",
	}

	check2 := &models.Check{
		ID:          "checkYetAnother",
		Name:        "nameNew",
		Group:       "groupNew",
		Description: "descriptionNew",
	}

	catalog := models.ChecksCatalog{check1, check2}

	err := suite.checksService.CreateChecksCatalog(catalog)

	suite.NoError(err)

	var checkEntity1 entities.Check
	var checkStored1 models.Check

	suite.tx.Where("id", "checkOther").First(&checkEntity1)

	json.Unmarshal(checkEntity1.Payload, &checkStored1)
	suite.Equal(check1, &checkStored1)

	var checkEntity2 entities.Check
	var checkStored2 models.Check

	suite.tx.Where("id", "checkYetAnother").First(&checkEntity2)

	json.Unmarshal(checkEntity2.Payload, &checkStored2)
	suite.Equal(check2, &checkStored2)

	// Count the number of checks is correct, and old entries have been removed
	var count int64
	suite.tx.Table("checks").Count(&count)
	suite.Equal(int64(2), count)
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetChecksResultAndMetadataByCluster() {
	results, err := suite.checksService.GetChecksResultAndMetadataByCluster("myClusterId")

	expectedResults := &models.ClusterCheckResults{
		Hosts: map[string]*models.CheckHost{
			"host1": &models.CheckHost{
				Reachable: true,
				Msg:       "",
			},
			"host2": &models.CheckHost{
				Reachable: false,
				Msg:       "error connecting",
			},
		},
		Checks: []models.ClusterCheckResult{
			models.ClusterCheckResult{
				ID:          "check1",
				Group:       "group1",
				Description: "description1",
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
				ID:          "check2",
				Group:       "group1",
				Description: "description2",
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

	suite.NoError(err)
	suite.Equal(expectedResults, results)
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetSelectedChecksById() {
	selectedChecks, err := suite.checksService.GetSelectedChecksById("group1")

	suite.NoError(err)
	suite.ElementsMatch([]string{"ABCDEF", "123456"}, selectedChecks.SelectedChecks)

	selectedChecks, err = suite.checksService.GetSelectedChecksById("group2")

	suite.NoError(err)
	suite.ElementsMatch([]string{"ABC123", "123ABC"}, selectedChecks.SelectedChecks)
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetSelectedChecksByIdError() {
	selectedChecks, err := suite.checksService.GetSelectedChecksById("other")

	suite.NoError(err)
	suite.EqualValues([]string{}, selectedChecks.SelectedChecks)
}

func (suite *ChecksServiceTestSuite) TestChecksService_CreateSelectedChecks() {
	err := suite.checksService.CreateSelectedChecks("group4", []string{"FEDCBA", "ABCDEF"})

	var selectedChecks models.SelectedChecks

	suite.tx.Where("id", "group4").First(&selectedChecks)
	expectedValue := models.SelectedChecks{
		ID:             "group4",
		SelectedChecks: []string{"FEDCBA", "ABCDEF"},
	}

	suite.NoError(err)
	suite.Equal(expectedValue, selectedChecks)

	// Check if an update works
	err = suite.checksService.CreateSelectedChecks("group4", []string{"ABCDEF", "FEDCBA"})

	suite.tx.Where("id", "group4").First(&selectedChecks)
	expectedValue = models.SelectedChecks{
		ID:             "group4",
		SelectedChecks: []string{"ABCDEF", "FEDCBA"},
	}

	suite.NoError(err)
	suite.Equal(expectedValue, selectedChecks)
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetConnectionSettingsByNode() {
	data, err := suite.checksService.GetConnectionSettingsByNode("node1")

	suite.NoError(err)
	suite.Equal(models.ConnectionSettings{ID: "group1", Node: "node1", User: "user1"}, data)

	data, err = suite.checksService.GetConnectionSettingsByNode("node2")

	suite.NoError(err)
	suite.Equal(models.ConnectionSettings{ID: "group1", Node: "node2", User: "user2"}, data)
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetConnectionSettingsByNodeError() {
	_, err := suite.checksService.GetConnectionSettingsByNode("other")

	suite.EqualError(err, "record not found")
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetConnectionSettingsById() {
	data, err := suite.checksService.GetConnectionSettingsById("group1")

	expectedData := map[string]models.ConnectionSettings{
		"node1": models.ConnectionSettings{ID: "group1", Node: "node1", User: "user1"},
		"node2": models.ConnectionSettings{ID: "group1", Node: "node2", User: "user2"},
	}
	suite.NoError(err)
	suite.Equal(expectedData, data)

	data, err = suite.checksService.GetConnectionSettingsById("group2")

	expectedData = map[string]models.ConnectionSettings{
		"node3": models.ConnectionSettings{ID: "group2", Node: "node3", User: "user3"},
	}
	suite.NoError(err)
	suite.Equal(expectedData, data)

	data, err = suite.checksService.GetConnectionSettingsById("other")

	expectedData = map[string]models.ConnectionSettings{}
	suite.NoError(err)
	suite.Equal(expectedData, data)
}

func (suite *ChecksServiceTestSuite) TestChecksService_CreateConnectionSettings() {
	err := suite.checksService.CreateConnectionSettings("group4", "node4", "user4")

	var data models.ConnectionSettings

	suite.tx.Where("id", "group4").First(&data)
	expectedValue := models.ConnectionSettings{ID: "group4", Node: "node4", User: "user4"}

	suite.NoError(err)
	suite.Equal(expectedValue, data)

	// Check if an update works
	err = suite.checksService.CreateConnectionSettings("group4", "node4", "user5")

	suite.tx.Where("id", "group4").First(&data)
	expectedValue = models.ConnectionSettings{ID: "group4", Node: "node4", User: "user5"}

	suite.NoError(err)
	suite.Equal(expectedValue, data)
}
