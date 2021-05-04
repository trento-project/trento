package sapsystem

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/pkg/errors"

	"github.com/SUSE/sap_host_exporter/lib/sapcontrol"
	"github.com/spf13/viper"

	"github.com/trento-project/trento/internal"
)

const sapInstallationPath string = "/usr/sap"
const sapIdentifierPattern string = "^[A-Z][A-Z0-9]{2}$" // PRD, HA1, etc
const sapInstancePattern string = "^[A-Z]*([0-9]{2})$"   // HDB00, ASCS00, ERS10, etc

type SAPSystemsList []*SAPSystem

type SAPSystem struct {
	webService sapcontrol.WebService
	Id         string                                  `mapstructure:"id,omitempty"`
	Type       string                                  `mapstructure:"type,omitempty"`
	Processes  map[string]*sapcontrol.OSProcess        `mapstructure:"processes,omitempty"`
	Instances  map[string]*sapcontrol.SAPInstance      `mapstructure:"instances,omitempty"`
	Properties map[string]*sapcontrol.InstanceProperty `mapstructure:"properties,omitempty"`
}

func NewSAPSystemsList() (SAPSystemsList, error) {
	var systems = SAPSystemsList{}

	instances, err := findSystems()
	if err != nil {
		return systems, errors.Wrap(err, "Error walking the path")
	}

	for _, i := range instances {
		s, err := NewSAPSystem(i)
		if err != nil {
			log.Printf("Error discovering a SAP system: %s", err)
			continue
		}
		systems = append(systems, &s)
	}

	return systems, nil
}

// The Id is a unique identifier of a SAP installation.
// It will be used to create totally independent SAP system data
func (s *SAPSystem) getSapSystemId() (string, error) {
	sid := s.Properties["SAPSYSTEMNAME"].Value
	return internal.Md5sum(fmt.Sprintf("/usr/sap/%s/SYS/global/security/rsecssfs/key/SSFS_%s.KEY", sid, sid))
}

// Find the installed SAP instances in the /usr/sap folder
func findSystems() ([]string, error) {
	var instances = []string{}

	_, err := os.Stat(sapInstallationPath)
	if os.IsNotExist(err) {
		log.Print("SAP installation not found")
		return instances, nil
	}

	files, err := ioutil.ReadDir(sapInstallationPath)
	if err != nil {
		return nil, err
	}

	reSAPIdentifier := regexp.MustCompile(sapIdentifierPattern)

	for _, f := range files {
		if reSAPIdentifier.MatchString(f.Name()) {
			log.Printf("New SAP system installation found: %s", f.Name())
			i, err := findInstances(path.Join(sapInstallationPath, f.Name()))
			if err != nil {
				log.Print(err.Error())
			}
			instances = append(instances, i[:]...)
		}
	}

	return instances, nil
}

// Find the installed SAP instances in the /usr/sap/${SID} folder
func findInstances(sapPath string) ([]string, error) {
	var instances = []string{}
	reSAPInstancer := regexp.MustCompile(sapInstancePattern)

	files, err := ioutil.ReadDir(sapPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		for _, matches := range reSAPInstancer.FindAllStringSubmatch(f.Name(), -1) {
			log.Printf("New SAP instance installation found: %s", matches[len(matches)-1])
			instances = append(instances, matches[len(matches)-1])
		}
	}

	return instances, nil
}

func NewSAPSystem(instNumber string) (SAPSystem, error) {
	var sapSystem = SAPSystem{
		Type:       "",
		Processes:  make(map[string]*sapcontrol.OSProcess),
		Instances:  make(map[string]*sapcontrol.SAPInstance),
		Properties: make(map[string]*sapcontrol.InstanceProperty),
	}

	config := viper.New()
	config.SetDefault("sap-control-uds", path.Join("/tmp", fmt.Sprintf(".sapstream5%s13", instNumber)))
	client := sapcontrol.NewSoapClient(config)
	sapSystem.webService = sapcontrol.NewWebService(client)

	properties, err := sapSystem.webService.GetInstanceProperties()
	if err != nil {
		return sapSystem, errors.Wrap(err, "SAPControl web service error")
	}

	for _, prop := range properties.Properties {
		sapSystem.Properties[prop.Property] = prop
	}

	processes, err := sapSystem.webService.GetProcessList()
	if err != nil {
		return sapSystem, errors.Wrap(err, "SAPControl web service error")
	}

	for _, proc := range processes.Processes {
		sapSystem.Processes[proc.Name] = proc
	}

	instances, err := sapSystem.webService.GetSystemInstanceList()
	if err != nil {
		return sapSystem, errors.Wrap(err, "SAPControl web service error")
	}

	for _, inst := range instances.Instances {
		sapSystem.Instances[inst.Hostname] = inst
	}

	_, ok := sapSystem.Properties["HANA Roles"]
	if ok {
		sapSystem.Type = "HANA"
	} else {
		sapSystem.Type = "APP"
	}

	id, err := sapSystem.getSapSystemId()
	if err == nil {
		sapSystem.Id = id
	}

	return sapSystem, nil
}
