package mocks

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/trento-project/trento/internal/sapsystem"
)

func NewDiscoveredSAPSystemDatabaseMock() sapsystem.SAPSystemsList {
	var s sapsystem.SAPSystemsList

	jsonFile, err := os.Open("./test/fixtures/discovery/sap_system/sap_system_discovery_database.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &s)

	return s
}

func NewDiscoveredSAPSystemApplicationMock() sapsystem.SAPSystemsList {
	var s sapsystem.SAPSystemsList

	jsonFile, err := os.Open("./test/fixtures/discovery/sap_system/sap_system_discovery_application.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &s)

	return s
}

func NewDiscoveredSAPSystemDiagnosticsMock() sapsystem.SAPSystemsList {
	var s sapsystem.SAPSystemsList

	jsonFile, err := os.Open("./test/fixtures/discovery/sap_system/sap_system_discovery_diagnostics.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &s)

	return s
}
