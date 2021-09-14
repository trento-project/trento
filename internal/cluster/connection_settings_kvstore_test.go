package cluster

import (
	"fmt"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
)

func TestGetConnectionSettings(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	consulInst.On("KV").Return(kv)

	kvPath := fmt.Sprintf(consul.KvClustersConnectionPath, "foobar")

	expectedConnData := map[string]interface{}{
		"host1": "myuser1",
		"host2": "myuser2",
	}
	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil)
	kv.On("ListMap", kvPath, kvPath).Return(expectedConnData, nil)

	connData, err := GetConnectionSettings(consulInst, "foobar")

	assert.NoError(t, err)
	assert.Equal(t, expectedConnData, connData)
	kv.AssertExpectations(t)
}

func TestStoreConnectionSettings(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	consulInst.On("KV").Return(kv)

	connData := map[string]interface{}{
		"host1": "myuser1",
		"host2": "myuser2",
	}
	kvPath := fmt.Sprintf(consul.KvClustersConnectionPath, "foobar")
	testLock := consulApi.Lock{}
	consulInst.On("AcquireLockKey", consul.KvClustersPath).Return(&testLock, nil)
	kv.On("PutMap", kvPath, connData).Return(nil)

	err := StoreConnectionSettings(consulInst, "foobar", connData)

	assert.NoError(t, err)
	kv.AssertExpectations(t)
}
