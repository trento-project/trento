package mocks

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/trento-project/trento/internal/cloud"
)

func NewDiscoveredCloudMock() cloud.CloudInstance {
	metadata := &cloud.AzureMetadata{}

	jsonFile, err := os.Open("../../test/fixtures/discovery/azure/azure_discovery.json")
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, metadata)

	return cloud.CloudInstance{
		Provider: "7783-7084-3265-9085-8269-3286-77",
		Metadata: metadata,
	}
}
