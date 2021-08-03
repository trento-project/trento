package environments

import (
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
)

func (e *Environment) getKVPath() string {
	return path.Join(consul.KvEnvironmentsPath, e.Name)
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
		err := loadHosts(client, environment)
		if err != nil {
			return nil, err
		}
		envs[env] = environment
	}

	return envs, nil
}

func loadHosts(client consul.Client, env *Environment) error {
	for landKey, landValue := range env.Landscapes {
		for sysKey, sysValue := range landValue.SAPSystems {
			query := hosts.CreateFilterMetaQuery(map[string][]string{
				"trento-sap-environment": []string{env.Name},
				"trento-sap-landscape":   []string{landKey},
				"trento-sap-systems":     []string{sysKey},
			})
			h, err := hosts.Load(client, query, []string{})
			if err != nil {
				return err
			}
			sysValue.Hosts = h
		}
	}
	return nil
}
