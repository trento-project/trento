package cloud

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

func (m *CloudInstance) getKVPath() string {
	host, _ := os.Hostname()
	kvPath := fmt.Sprintf(consul.KvHostsClouddataPath, host)

	return kvPath
}

func (m *CloudInstance) Store(client consul.Client) error {
	host, _ := os.Hostname()
	l, err := client.AcquireLockKey(path.Join(consul.KvHostsPath, host) + "/")
	if err != nil {
		return errors.Wrap(err, "could not lock the kv for cloud data")
	}
	defer l.Unlock()

	kvPath := m.getKVPath()

	// Clean the current data before storing the new values
	_, err = client.KV().DeleteTree(kvPath, nil)
	if err != nil {
		return errors.Wrap(err, "Error deleting cloud data content")
	}

	cloudDataMap := make(map[string]interface{})
	mapstructure.Decode(m, &cloudDataMap)

	err = client.KV().PutMap(kvPath, cloudDataMap)
	if err != nil {
		return errors.Wrap(err, "Error storing a cloud data")
	}

	return nil
}

func Load(client consul.Client, host string) (*CloudInstance, error) {
	err := client.WaitLock(path.Join(consul.KvHostsPath, host) + "/")
	if err != nil {
		return nil, errors.Wrap(err, "error waiting for the lock for cloud data")
	}

	kvPath := fmt.Sprintf(consul.KvHostsClouddataPath, host)

	data, err := client.KV().ListMap(kvPath, kvPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for clou data KV values")
	}

	cloudData := &CloudInstance{}
	mapstructure.Decode(data, &cloudData)

	switch cloudData.Provider {
	case azure:
		azureMetadata := &AzureMetadata{}
		mapstructure.Decode(cloudData.Metadata, &azureMetadata)
		cloudData.Metadata = azureMetadata
	}

	return cloudData, nil
}
