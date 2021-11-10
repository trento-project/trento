package mocks

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/trento-project/trento/internal/sapsystem"
)

func NewDiscoveredSAPSystemMock() sapsystem.SAPSystemsList {
	var s sapsystem.SAPSystemsList

	jsonFile, err := os.Open("./test/fixtures/discovery/sap_system/sap_system_discovery.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &s)

	return s
}
