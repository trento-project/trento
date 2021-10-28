package mocks

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/trento-project/trento/internal/subscription"
)

func NewDiscoveredSubscriptionsMock() subscription.Subscriptions {
	var subs subscription.Subscriptions

	jsonFile, err := os.Open("../../../test/fixtures/discovery/subscriptions/subscriptions_discovery.json")
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &subs)

	return subs
}
