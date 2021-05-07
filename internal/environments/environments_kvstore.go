package environments

import (
	"fmt"

  "github.com/mitchellh/mapstructure"
  "github.com/pkg/errors"

  "github.com/trento-project/trento/internal/consul"
)


func (e *Environment) getKVPath() string {
	return fmt.Sprintf("%s/%s", consul.KvEnvironmentsPath, e.Name)
}

func (e *Environment) Store(client consul.Client) error {
	envMap := make(map[string]interface{})
	mapstructure.Decode(e, &envMap)

	err := client.KV().PutMap(e.getKVPath(), envMap)
	if err != nil {
		return errors.Wrap(err, "Error storing a environment data")
	}

  return nil
}

func Load(client consul.Client) (map[string]*Environment, error) {
	var envs = map[string]*Environment{}

	entries, err := client.KV().ListMap(consul.KvEnvironmentsPath, consul.KvEnvironmentsPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for Environments KV values")
	}

	for env, envValue := range entries {
		environment := &Environment{}
		mapstructure.Decode(envValue, &environment)
		envs[env] = environment
	}

	return envs, nil
}
