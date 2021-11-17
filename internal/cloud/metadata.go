package cloud

import (
	"os/exec"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	Azure = "azure"
	Aws   = "aws"
	Gcp   = "gcp"
	// DMI chassis asset tag for Azure machines, needed to identify wether or not we are running on Azure
	// This is actually ASCII-encoded, the decoding into a string results in "MSFT AZURE VM"
	azureDmiTag = "7783-7084-3265-9085-8269-3286-77"
)

type CloudInstance struct {
	Provider string      `mapstructure:"provider,omitempty"`
	Metadata interface{} `mapstructure:"metadata,omitempty"`
}

type CustomCommand func(name string, arg ...string) *exec.Cmd

var customExecCommand CustomCommand = exec.Command

// All these detection methods are based in crmsh code, which has been refined over the years
// https://github.com/ClusterLabs/crmsh/blob/master/crmsh/utils.py#L2009

func identifyAzure() (bool, error) {
	log.Debug("Checking if the VM is running on Azure...")
	output, err := customExecCommand("dmidecode", "-s", "chassis-asset-tag").Output()
	if err != nil {
		return false, err
	}

	provider := strings.TrimSpace(string(output))
	log.Debugf("dmidecode output: %s", provider)

	return provider == azureDmiTag, nil
}

func identifyAws() (bool, error) {
	log.Debug("Checking if the VM is running on Aws...")
	output, err := customExecCommand("dmidecode", "-s", "system-version").Output()
	if err != nil {
		return false, err
	}

	provider := strings.TrimSpace(string(output))
	log.Debugf("dmidecode output: %s", provider)

	return regexp.MatchString(".*amazon.*", provider)
}

func identifyGcp() (bool, error) {
	log.Debug("Checking if the VM is running on Gcp...")
	output, err := customExecCommand("dmidecode", "-s", "bios-vendor").Output()
	if err != nil {
		return false, err
	}

	provider := strings.TrimSpace(string(output))
	log.Debugf("dmidecode output: %s", provider)

	return regexp.MatchString(".*Google.*", provider)
}

func IdentifyCloudProvider() (string, error) {
	log.Info("Identifying if the VM is running in a cloud environment...")

	if result, err := identifyAzure(); err != nil {
		return "", err
	} else if result {
		log.Infof("VM is running on %s", Azure)
		return Azure, nil
	}

	if result, err := identifyAws(); err != nil {
		return "", err
	} else if result {
		log.Infof("VM is running on %s", Aws)
		return Aws, nil
	}

	if result, err := identifyGcp(); err != nil {
		return "", err
	} else if result {
		log.Infof("VM is running on %s", Gcp)
		return Gcp, nil
	}

	log.Info("VM is not running in any recognized cloud provider")
	return "", nil
}

func NewCloudInstance() (*CloudInstance, error) {
	var err error
	var cloudMetadata interface{}

	provider, err := IdentifyCloudProvider()
	if err != nil {
		return nil, err
	}

	cInst := &CloudInstance{
		Provider: provider,
	}

	switch provider {
	case Azure:
		cloudMetadata, err = NewAzureMetadata()
		if err != nil {
			return nil, err
		}
	}

	cInst.Metadata = cloudMetadata

	return cInst, nil

}
