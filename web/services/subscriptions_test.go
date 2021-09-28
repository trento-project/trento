package services

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/internal/subscription"
)

func TestGetHostSubscriptions(t *testing.T) {
	host, _ := os.Hostname()
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	listMap := map[string]interface{}{
		"0000": map[string]string{
			"identifier":          "SLES_SAP",
			"version":             "15.2",
			"arch":                "x86_64",
			"status":              "Registered",
			"starts_at":           "2019-03-20 09:55:32 UTC",
			"expires_at":          "2024-03-20 09:55:32 UTC",
			"subscription_status": "ACTIVE",
			"type":                "internal",
		},
		"0001": map[string]string{
			"identifier": "sle-module-public-cloud",
			"version":    "15.2",
			"arch":       "x86_64",
			"status":     "Registered",
		},
	}

	kvPath := fmt.Sprintf(consul.KvHostsSubscriptionsPath, host)

	kv.On("ListMap", kvPath, kvPath).Return(listMap, nil)
	consulInst.On("WaitLock", path.Join(consul.KvHostsPath, host)+"/").Return(nil)

	consulInst.On("KV").Return(kv)

	subsService := NewSubscriptionsService(consulInst)
	subs, err := subsService.GetHostSubscriptions(host)

	expectedSubs := subscription.Subscriptions{
		&subscription.Subscription{
			Identifier:         "SLES_SAP",
			Version:            "15.2",
			Arch:               "x86_64",
			Status:             "Registered",
			StartsAt:           "2019-03-20 09:55:32 UTC",
			ExpiresAt:          "2024-03-20 09:55:32 UTC",
			SubscriptionStatus: "ACTIVE",
			Type:               "internal",
		},
		&subscription.Subscription{
			Identifier: "sle-module-public-cloud",
			Version:    "15.2",
			Arch:       "x86_64",
			Status:     "Registered",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedSubs, subs)
}
