package environments

import (
	consulApi "github.com/hashicorp/consul/api"
	"github.com/trento-project/trento/internal/hosts"
)

const defaultName = "default"

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
	Type  int            `mapstructure:"type,omitempty"`
	Hosts hosts.HostList `mapstructure:"hosts,omitempty"`
}

func NewEnvironment(name string, landscapes ...*Landscape) *Environment {
	landsMap := make(map[string]*Landscape)
	for _, l := range landscapes {
		landsMap[l.Name] = l
	}

	return &Environment{Name: name, Landscapes: landsMap}
}

func NewLandscape(name string, systems ...*SAPSystem) *Landscape {
	systemsMap := make(map[string]*SAPSystem)

	for _, s := range systems {
		systemsMap[s.Name] = s
	}
	return &Landscape{Name: name, SAPSystems: systemsMap}
}

func NewSystem(sysName string, sysType int) *SAPSystem {
	return &SAPSystem{Name: sysName, Type: sysType}
}

func NewDefaultEnvironment() *Environment {
	return NewEnvironment(defaultName)
}

func NewDefaultLandscape() *Landscape {
	return NewLandscape(defaultName)
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

func (e *Environment) AddLandscape(land *Landscape) {
	e.Landscapes[land.Name] = land
}

func (l *Landscape) Health() EnvironmentHealth {
	var health = NewEnvironmentHealth()

	for _, system := range l.SAPSystems {
		h := system.Health().Health
		health = health.updateHealth(system.Name, h)
	}

	return health
}

func (l *Landscape) AddSystem(system *SAPSystem) {
	l.SAPSystems[system.Name] = system
}

func (s *SAPSystem) Health() EnvironmentHealth {
	var health = NewEnvironmentHealth()

	for _, host := range s.Hosts {
		h := host.Health()
		health = health.updateHealth(host.Name(), h)
	}

	return health
}
