package ruleset

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

func (r RuleSets) getKVPath(host string) string {
	kvPath := fmt.Sprintf(consul.KvHostsRuleSetsPath, host)

	return kvPath
}

func (r RuleSets) Store(client consul.Client, host string) error {
	kvPath := r.getKVPath(host)

	_, err := client.KV().DeleteTree(kvPath, nil)
	if err != nil {
		return errors.Wrap(err, "Error deleting rulesets content")
	}

	rulesetsSlice := make([]interface{}, 0)

	for _, ruleSet := range r {
		rulesetMap := make(map[string]interface{})
		mapstructure.Decode(ruleSet, &rulesetMap)

		rulesetsSlice = append(rulesetsSlice, rulesetMap)

	}

	err = client.KV().PutSlice(kvPath, rulesetsSlice)
	if err != nil {
		return errors.Wrap(err, "Error storing a host rulesets")
	}

	return nil
}

func Load(client consul.Client, host string) (RuleSets, error) {
	var ruleSets = RuleSets{}
	kvPath := ruleSets.getKVPath(host)

	entries, err := client.KV().ListMap(kvPath, kvPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for rulesets KV values")
	}

	for _, ruleSet := range entries {
		r := &RuleSet{}
		mapstructure.Decode(ruleSet, &r)
		ruleSets = append(ruleSets, r)
	}

	return ruleSets, nil
}
