package environments

import (
  "github.com/trento-project/trento/internal/consul"
)

type Environment struct {
	Name       string                `mapstructure:"name,omitempty"`
	Landscapes map[string]*Landscape `mapstructure:"landscapes,omitempty"`
}

type Landscape struct {
	Name       string                `mapstructure:"name,omitempty"`
	SAPSystems map[string]*SAPSystem `mapstructure:"sapsystems,omitempty"`
}

type SAPSystem struct {
	Name  string `mapstructure:"name,omitempty"`
}

func NewEnvironment(env, land, sys string) Environment {
  newSys := &SAPSystem{Name: sys}
  systemsMap := make(map[string]*SAPSystem)
  systemsMap[sys] = newSys

  newLand := &Landscape{Name: consul.KvUngrouped, SAPSystems: systemsMap}
  landsMap := make(map[string]*Landscape)
  landsMap[land] = newLand

  newEnv := Environment{Name: env, Landscapes: landsMap}
  return newEnv
}
