package environments

import (
  "github.com/trento-project/trento/internal/consul"
  "github.com/trento-project/trento/internal/hosts"
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
	Name  string         `mapstructure:"name,omitempty"`
  Type  string         `mapstructure:"type,omitempty"`
  Hosts hosts.HostList `mapstructure:"hosts,omitempty"`
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

func ungrouped(name string) bool {
  return name == consul.KvUngrouped
}

func (e *Environment) Ungrouped() bool {
	return ungrouped(e.Name)
}

func (l *Landscape) Ungrouped() bool {
	return ungrouped(l.Name)
}

func (s *SAPSystem) Ungrouped() bool {
	return ungrouped(s.Name)
}

type EnvironmentHealth struct {
	Health    string
	HealthMap map[string]string
}

func (e *EnvironmentHealth) updateHealth(n string, h string) EnvironmentHealth {
	e.HealthMap[n] = h

	if h == "critical" {
		e.Health = h
	} else if h == "warning" && e.Health != "critical" {
		e.Health = h
	}

	return *e
}

func (e *Environment) Health() EnvironmentHealth {
	var health = EnvironmentHealth{
		Health:    "passing",
		HealthMap: make(map[string]string),
	}

	for _, land := range e.Landscapes {
		h := land.Health().Health
		health = health.updateHealth(land.Name, h)
	}

	return health
}

func (l *Landscape) Health() EnvironmentHealth {
	var health = EnvironmentHealth{
		Health:    "passing",
		HealthMap: make(map[string]string),
	}

	for _, system := range l.SAPSystems {
		h := system.Health().Health
		health = health.updateHealth(system.Name, h)
	}

	return health
}

func (s *SAPSystem) Health() EnvironmentHealth {
	var health = EnvironmentHealth{
		Health:    "passing",
		HealthMap: make(map[string]string),
	}

	for _, host := range s.Hosts {
		h := host.Health()
		health = health.updateHealth(host.Name(), h)
	}

	return health
}
