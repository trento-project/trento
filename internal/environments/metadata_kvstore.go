package environments

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

func (m *Metadata) getKVPath() string {
	host, _ := os.Hostname()
	kvPath := fmt.Sprintf(consul.KvHostsMetadataPath, host)

	return kvPath
}

func (m *Metadata) Store(client consul.Client) error {
	kvPath := m.getKVPath()

	metadataMap := make(map[string]interface{})
	mapstructure.Decode(m, &metadataMap)

	err := client.KV().PutMap(kvPath, metadataMap)
	if err != nil {
		return errors.Wrap(err, "Error storing a host metadata")
	}

	return nil
}
