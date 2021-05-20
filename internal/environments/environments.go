package environments

import (
	consulApi "github.com/hashicorp/consul/api"
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

	newLand := &Landscape{Name: land, SAPSystems: systemsMap}
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

func NewEnvironmentHealth() EnvironmentHealth {
	return EnvironmentHealth{
		Health:    consulApi.HealthPassing,
		HealthMap: make(map[string]string),
	}
}

func (e *EnvironmentHealth) updateHealth(entry string, health string) EnvironmentHealth {
	e.HealthMap[entry] = health

	if health == consulApi.HealthCritical {
		e.Health = health
	} else if health == consulApi.HealthWarning && e.Health != consulApi.HealthCritical {
		e.Health = health
	}

	return *e
}

func (e *Environment) Health() EnvironmentHealth {
	var health = NewEnvironmentHealth()

	for _, land := range e.Landscapes {
		h := land.Health().Health
		health = health.updateHealth(land.Name, h)
	}

	return health
}

func (l *Landscape) Health() EnvironmentHealth {
	var health = NewEnvironmentHealth()

	for _, system := range l.SAPSystems {
		h := system.Health().Health
		health = health.updateHealth(system.Name, h)
	}

	return health
}

func (s *SAPSystem) Health() EnvironmentHealth {
	var health = NewEnvironmentHealth()

	for _, host := range s.Hosts {
		h := host.Health()
		health = health.updateHealth(host.Name(), h)
	}

	return health
}
