package cluster

import (
	"path"

	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

func GetCheckSelection(client consul.Client, clusterId string) (string, error) {
	err := client.WaitLock(consul.KvClustersPath)
	if err != nil {
		return "", errors.Wrap(err, "error waiting for the lock for clusters")
	}

	kvPath := path.Join(consul.KvClustersPath, clusterId, "checks")

	pair, _, err := client.KV().Get(kvPath, nil)
	if err != nil {
		return "", errors.Wrap(err, "error getting the cluster checks selection")
	}

	return string(pair.Value), nil
}

func StoreCheckSelection(client consul.Client, clusterId, selected string) error {

	l, err := client.AcquireLockKey(consul.KvClustersPath)
	if err != nil {
		return errors.Wrap(err, "could not lock the kv for clusters")
	}
	defer l.Unlock()

	kvPath := path.Join(consul.KvClustersPath, clusterId, "checks")

	err = client.KV().PutTyped(kvPath, selected)
	if err != nil {
		return errors.Wrap(err, "Error storing cluster checks selection")
	}

	return nil
}
