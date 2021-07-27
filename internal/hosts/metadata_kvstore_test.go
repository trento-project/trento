package hosts

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
)

func TestMetadataStore(t *testing.T) {
	host, _ := os.Hostname()
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	consulInst.On("KV").Return(kv)

	m := Metadata{
		Cluster:     "test-cluster",
		Environment: "env1",
		Landscape:   "land1",
		SAPSystems:  "sys1",
	}

	expectedPutMap := map[string]interface{}{
		"ha-cluster":      "test-cluster",
		"sap-environment": "env1",
		"sap-landscape":   "land1",
		"sap-systems":     "sys1",
	}

	kvPath := fmt.Sprintf(consul.KvHostsMetadataPath, host)
	kv.On("PutMap", kvPath, expectedPutMap).Return(nil, nil)

	result := m.Store(consulInst)

	assert.Equal(t, result, nil)
}
