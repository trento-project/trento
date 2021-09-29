package subscription

import (
	"fmt"
	"os"
	"path"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
)

func TestStore(t *testing.T) {
	host, _ := os.Hostname()
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	consulInst.On("KV").Return(kv)

	kvPath := fmt.Sprintf(consul.KvHostsSubscriptionsPath, host)

	expectedPutInterface := []interface{}{
		&Subscription{
			Identifier:         "SLES_SAP",
			Version:            "15.2",
			Arch:               "x86_64",
			Status:             "Registered",
			StartsAt:           "2019-03-20 09:55:32 UTC",
			ExpiresAt:          "2024-03-20 09:55:32 UTC",
			SubscriptionStatus: "ACTIVE",
			Type:               "internal",
		},
		&Subscription{
			Identifier: "sle-module-public-cloud",
			Version:    "15.2",
			Arch:       "x86_64",
			Status:     "Registered",
		},
	}

	kv.On("DeleteTree", kvPath, (*consulApi.WriteOptions)(nil)).Return(nil, nil)
	kv.On("PutInterface", kvPath, expectedPutInterface).Return(nil, nil)
	testLock := consulApi.Lock{}
	consulInst.On("AcquireLockKey", path.Join(consul.KvHostsPath, host)+"/").Return(&testLock, nil)

	s := Subscriptions{
		&Subscription{
			Identifier:         "SLES_SAP",
			Version:            "15.2",
			Arch:               "x86_64",
			Status:             "Registered",
			StartsAt:           "2019-03-20 09:55:32 UTC",
			ExpiresAt:          "2024-03-20 09:55:32 UTC",
			SubscriptionStatus: "ACTIVE",
			Type:               "internal",
		},
		&Subscription{
			Identifier: "sle-module-public-cloud",
			Version:    "15.2",
			Arch:       "x86_64",
			Status:     "Registered",
		},
	}

	result := s.Store(consulInst)
	assert.Equal(t, nil, result)
}

func TestLoad(t *testing.T) {
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

	subs, _ := Load(consulInst, host)

	expectedSubs := Subscriptions{
		&Subscription{
			Identifier:         "SLES_SAP",
			Version:            "15.2",
			Arch:               "x86_64",
			Status:             "Registered",
			StartsAt:           "2019-03-20 09:55:32 UTC",
			ExpiresAt:          "2024-03-20 09:55:32 UTC",
			SubscriptionStatus: "ACTIVE",
			Type:               "internal",
		},
		&Subscription{
			Identifier: "sle-module-public-cloud",
			Version:    "15.2",
			Arch:       "x86_64",
			Status:     "Registered",
		},
	}

	assert.ElementsMatch(t, expectedSubs, subs)
}
