package cloud

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

	kvPath := fmt.Sprintf(consul.KvHostsClouddataPath, host)

	expectedPutMap := map[string]interface{}{
		"provider": "azure",
		"metadata": &AzureMetadata{
			Compute: Compute{
				Name: "test",
			},
		},
	}

	kv.On("DeleteTree", kvPath, (*consulApi.WriteOptions)(nil)).Return(nil, nil)
	kv.On("PutMap", kvPath, expectedPutMap).Return(nil, nil)
	testLock := consulApi.Lock{}
	consulInst.On("AcquireLockKey", path.Join(consul.KvHostsPath, host)+"/").Return(&testLock, nil)

	m := CloudInstance{
		Provider: "azure",
		Metadata: &AzureMetadata{
			Compute: Compute{
				Name: "test",
			},
		},
	}

	result := m.Store(consulInst)
	assert.Equal(t, nil, result)
}

func TestLoadAzure(t *testing.T) {
	host, _ := os.Hostname()
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	listMap := map[string]interface{}{
		"provider": "azure",
		"metadata": map[string]interface{}{
			"compute": map[string]interface{}{
				"name": "test",
			},
		},
	}

	kvPath := fmt.Sprintf(consul.KvHostsClouddataPath, host)

	kv.On("ListMap", kvPath, kvPath).Return(listMap, nil)
	consulInst.On("WaitLock", path.Join(consul.KvHostsPath, host)+"/").Return(nil)

	consulInst.On("KV").Return(kv)

	m, _ := Load(consulInst, host)

	assert.Equal(t, "azure", m.Provider)
	meta := m.Metadata.(*AzureMetadata)
	assert.Equal(t, "test", meta.Compute.Name)
}
