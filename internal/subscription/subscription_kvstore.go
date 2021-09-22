package subscription

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

func (s *Subscriptions) getKVPath() string {
	host, _ := os.Hostname()
	kvPath := fmt.Sprintf(consul.KvHostsSubscriptionsPath, host)

	return kvPath
}

func (s *Subscriptions) Store(client consul.Client) error {
	host, _ := os.Hostname()
	l, err := client.AcquireLockKey(path.Join(consul.KvHostsPath, host) + "/")
	if err != nil {
		return errors.Wrap(err, "could not lock the kv for cloud data")
	}
	defer l.Unlock()

	kvPath := s.getKVPath()

	// Clean the current data before storing the new values
	_, err = client.KV().DeleteTree(kvPath, nil)
	if err != nil {
		return errors.Wrap(err, "Error deleting subscription data content")
	}

	var subsData []interface{}
	mapstructure.Decode(s, &subsData)

	err = client.KV().PutInterface(kvPath, subsData)
	if err != nil {
		return errors.Wrap(err, "Error storing subscriptions data")
	}

	return nil
}

func Load(client consul.Client, host string) (Subscriptions, error) {
	err := client.WaitLock(path.Join(consul.KvHostsPath, host) + "/")
	if err != nil {
		return nil, errors.Wrap(err, "error waiting for the lock for subscriptions data")
	}

	kvPath := fmt.Sprintf(consul.KvHostsSubscriptionsPath, host)

	data, err := client.KV().ListMap(kvPath, kvPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for subscriptions data KV values")
	}

	subsData := Subscriptions{}
	for _, sub := range data {
		newSub := &Subscription{}
		mapstructure.Decode(sub, &newSub)
		subsData = append(subsData, newSub)
	}

	return subsData, nil
}
