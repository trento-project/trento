package services

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/internal/subscription"
)

func TestGetSubscriptionData(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	kv := new(mocks.KV)

	nodes := []*consulApi.Node{
		{
			Node: "node1",
		},
		{
			Node: "node2",
		},
		{
			Node: "node3",
		},
	}

	catalog.On("Nodes", &consulApi.QueryOptions{}).Return(nodes, nil, nil)

	sub1 := map[string]interface{}{
		"0000": map[string]string{
			"identifier": "SLES_SAP",
		},
		"0001": map[string]string{
			"identifier": "sle-module-public-cloud",
		},
	}

	sub2 := map[string]interface{}{
		"0000": map[string]string{
			"identifier": "SLES_SAP",
		},
	}

	sub3 := map[string]interface{}{
		"0000": map[string]string{
			"identifier": "sle-module-public-cloud",
		},
	}

	kvPath := fmt.Sprintf(consul.KvHostsSubscriptionsPath, "node1")
	kv.On("ListMap", kvPath, kvPath).Return(sub1, nil)
	consulInst.On("WaitLock", path.Join(consul.KvHostsPath, "node1")+"/").Return(nil)

	kvPath = fmt.Sprintf(consul.KvHostsSubscriptionsPath, "node2")
	kv.On("ListMap", kvPath, kvPath).Return(sub2, nil)
	consulInst.On("WaitLock", path.Join(consul.KvHostsPath, "node2")+"/").Return(nil)

	kvPath = fmt.Sprintf(consul.KvHostsSubscriptionsPath, "node3")
	kv.On("ListMap", kvPath, kvPath).Return(sub3, nil)
	consulInst.On("WaitLock", path.Join(consul.KvHostsPath, "node3")+"/").Return(nil)

	consulInst.On("Catalog").Return(catalog)
	consulInst.On("KV").Return(kv)

	subsService := NewSubscriptionsService(consulInst)
	subData, err := subsService.GetSubscriptionData()

	expectedSubData := &SubscriptionData{
		Type:            Premium,
		SubscribedCount: 2,
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedSubData, subData)
}

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
	assert.ElementsMatch(t, expectedSubs, subs)
}
