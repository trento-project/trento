package services

import (
	"encoding/json"
	"testing"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"

	"github.com/trento-project/trento/web/models"
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

type ChecksServiceTestSuite struct {
	suite.Suite
	db            *gorm.DB
	tx            *gorm.DB
	checksService ChecksService
}

func TestChecksServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ChecksServiceTestSuite))
}

func (suite *ChecksServiceTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(
		models.CheckRaw{}, models.CheckResultsRaw{}, models.SelectedChecks{}, models.ConnectionSettings{})
	loadChecksCatalogFixtures(suite.db)
	loadChecksResultsFixtures(suite.db)
	loadSelectedChecksFixtures(suite.db)
	loadConnectionSettingsFixtures(suite.db)
}

func (suite *ChecksServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(models.CheckRaw{})
	suite.db.Migrator().DropTable(models.CheckResultsRaw{})
	suite.db.Migrator().DropTable(models.SelectedChecks{})
	suite.db.Migrator().DropTable(models.ConnectionSettings{})
}

func (suite *ChecksServiceTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.checksService = NewChecksService(suite.tx)
}

func (suite *ChecksServiceTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func loadChecksCatalogFixtures(db *gorm.DB) {
	check1payload := `{"id":"check1","name":"name1","group":"group1","description":"description1"}`
	db.Create(&models.CheckRaw{
		ID:      "check1",
		Payload: datatypes.JSON([]byte(check1payload)),
	})
	check3payload := `{"id":"check3","name":"name3","group":"group2","description":"description3"}`
	db.Create(&models.CheckRaw{
		ID:      "check3",
		Payload: datatypes.JSON([]byte(check3payload)),
	})
	check2payload := `{"id":"check2","name":"name2","group":"group1","description":"description2"}`
	db.Create(&models.CheckRaw{
		ID:      "check2",
		Payload: datatypes.JSON([]byte(check2payload)),
	})
}

func loadChecksResultsFixtures(db *gorm.DB) {
	group1payload := `{"hosts":{"host1":{"reachable": true, "msg":""},"host2":{"reachable":false,"msg":"error connecting"}},
	"checks":{"check1":{"hosts":{"host1":{"result":"critical"},"host2":{"result":"critical"}}},
	"check2":{"hosts":{"host1":{"result":"critical"}, "host2":{"result":"critical"}}}}}`
	db.Create(&models.CheckResultsRaw{
		GroupID: "group1",
		Payload: datatypes.JSON([]byte(group1payload)),
	})
	group1payloadLast := `{"hosts":{"host1":{"reachable": true, "msg":""},"host2":{"reachable":false,"msg":"error connecting"}},
	"checks":{"check1":{"hosts":{"host1":{"result":"passing"},"host2":{"result":"passing"}}},
	"check2":{"hosts":{"host1":{"result":"warning"}, "host2":{"result":"critical"}}}}}`
	db.Create(&models.CheckResultsRaw{
		GroupID: "group1",
		Payload: datatypes.JSON([]byte(group1payloadLast)),
	})
	group2payload := `{"hosts":{"host3":{"reachable":true "msg":""},"host4":{"reachable":true,"msg":""}},
	"checks":{"check1":{"hosts":{"host3":{"result":"critical"},"host4":{"result":"critical"}}},
	"check2":{"hosts":{"host3":{"result":"passing"},"host4":{"result":"warning"}}}}}`
	db.Create(&models.CheckResultsRaw{
		GroupID: "group2",
		Payload: datatypes.JSON([]byte(group2payload)),
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

	var newCheckRaw models.CheckRaw
	var newCheck models.Check

	suite.tx.Where("id", "checkNew").First(&newCheckRaw)

	json.Unmarshal(newCheckRaw.Payload, &newCheck)

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

	suite.tx.Where("id", "checkNew").First(&newCheckRaw)

	json.Unmarshal(newCheckRaw.Payload, &newCheck)

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

	var checkRaw1 models.CheckRaw
	var checkStored1 models.Check

	suite.tx.Where("id", "checkOther").First(&checkRaw1)

	json.Unmarshal(checkRaw1.Payload, &checkStored1)
	suite.Equal(check1, &checkStored1)

	var checkRaw2 models.CheckRaw
	var checkStored2 models.Check

	suite.tx.Where("id", "checkYetAnother").First(&checkRaw2)

	json.Unmarshal(checkRaw2.Payload, &checkStored2)
	suite.Equal(check2, &checkStored2)
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetChecksResultById() {
	results, err := suite.checksService.GetChecksResultById("group1")

	var resultsRaw models.CheckResultsRaw
	var resultsStored models.Results

	suite.tx.Where("group_id", "group1").Last(&resultsRaw)

	json.Unmarshal(resultsRaw.Payload, &resultsStored)
	suite.NoError(err)
	suite.Equal(&resultsStored, results)
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetChecksResultByIdError() {
	_, err := suite.checksService.GetChecksResultById("other")

	suite.EqualError(err, "record not found")
}

func (suite *ChecksServiceTestSuite) TestChecksService_CreateChecksResultsById() {
	results := &models.Results{
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
		Checks: map[string]*models.ChecksByHost{
			"check1": &models.ChecksByHost{
				Hosts: map[string]*models.Check{
					"host1": &models.Check{
						Result: models.CheckCritical,
					},
					"host2": &models.Check{
						Result: models.CheckCritical,
					},
				},
			},
			"check2": &models.ChecksByHost{
				Hosts: map[string]*models.Check{
					"host1": &models.Check{
						Result: models.CheckWarning,
					},
					"host2": &models.Check{
						Result: models.CheckPassing,
					},
				},
			},
		},
	}

	err := suite.checksService.CreateChecksResultsById("group1", results)

	var resultsRaw models.CheckResultsRaw
	var resultsStored models.Results

	suite.tx.Where("group_id", "group1").Last(&resultsRaw)

	json.Unmarshal(resultsRaw.Payload, &resultsStored)
	suite.NoError(err)
	suite.Equal(results, &resultsStored)
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetAggregatedChecksResultByHost() {
	results, err := suite.checksService.GetAggregatedChecksResultByHost("group1")

	expectedResults := map[string]*AggregatedCheckData{
		"host1": &AggregatedCheckData{
			PassingCount:  1,
			WarningCount:  1,
			CriticalCount: 0,
		},
		"host2": &AggregatedCheckData{
			PassingCount:  1,
			WarningCount:  0,
			CriticalCount: 1,
		},
	}

	suite.NoError(err)
	suite.Equal(expectedResults, results)
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetAggregatedChecksResultById() {
	results, err := suite.checksService.GetAggregatedChecksResultById("group1")

	expectedResults := &AggregatedCheckData{
		PassingCount:  2,
		WarningCount:  1,
		CriticalCount: 1,
	}

	suite.NoError(err)
	suite.Equal(expectedResults, results)
}

func (suite *ChecksServiceTestSuite) TestChecksService_GetChecksResultAndMetadataById() {
	results, err := suite.checksService.GetChecksResultAndMetadataById("group1")

	expectedResults := &models.ResultsAsList{
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
		Checks: []*models.ChecksByHost{
			&models.ChecksByHost{
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
			&models.ChecksByHost{
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
	_, err := suite.checksService.GetSelectedChecksById("other")

	suite.EqualError(err, "record not found")
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
