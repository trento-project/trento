package cluster

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

func GetConnectionSettings(client consul.Client, clusterId string) (map[string]interface{}, error) {
	err := client.WaitLock(consul.KvClustersPath)
	if err != nil {
		return nil, errors.Wrap(err, "error waiting for the lock for clusters")
	}

	kvPath := fmt.Sprintf(consul.KvClustersConnectionPath, clusterId)

	connectionData, err := client.KV().ListMap(kvPath, kvPath)
	if err != nil {
		return nil, errors.Wrap(err, "error getting the cluster connection settings")
	}

	return connectionData, nil
}

func StoreConnectionSettings(client consul.Client, clusterId string, connectionMap map[string]interface{}) error {

	l, err := client.AcquireLockKey(consul.KvClustersPath)
	if err != nil {
		return errors.Wrap(err, "could not lock the kv for clusters")
	}
	defer l.Unlock()

	kvPath := fmt.Sprintf(consul.KvClustersConnectionPath, clusterId)

	err = client.KV().PutMap(kvPath, connectionMap)
	if err != nil {
		return errors.Wrap(err, "Error storing cluster connection settings")
	}

	return nil
}
