package cluster

import (
	"fmt"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
)

func TestGetCheckSelection(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	consulInst.On("KV").Return(kv)

	kvPath := fmt.Sprintf(consul.KvClustersChecksPath, "foobar")
	pair := &consulApi.KVPair{
		Value: []byte("1.1.1,1.2.3"),
	}

	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil)
	kv.On("Get", kvPath, (*consulApi.QueryOptions)(nil)).Return(pair, nil, nil)

	selectedChecks, err := GetCheckSelection(consulInst, "foobar")

	assert.NoError(t, err)
	assert.Equal(t, "1.1.1,1.2.3", selectedChecks)
	kv.AssertExpectations(t)
}

func TestStoreCheckSelection(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	consulInst.On("KV").Return(kv)

	kvPath := fmt.Sprintf(consul.KvClustersChecksPath, "foobar")
	testLock := consulApi.Lock{}
	consulInst.On("AcquireLockKey", consul.KvClustersPath).Return(&testLock, nil)
	kv.On("PutTyped", kvPath, "1.1.1,1.2.3").Return(nil)

	err := StoreCheckSelection(consulInst, "foobar", "1.1.1,1.2.3")

	assert.NoError(t, err)
	kv.AssertExpectations(t)
}
