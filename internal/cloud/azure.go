/*
Based on https://docs.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service?tabs=linux#instance-metadata
*/

package cloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	log "github.com/sirupsen/logrus"
)

const (
	azureApiVersion = "2021-02-01"
	azureApiAddress = "169.254.169.254"
	azurePortalUrl  = "https://portal.azure.com/#@SUSERDBillingsuse.onmicrosoft.com/resource"
)

type AzureMetadata struct {
	Compute Compute `json:"compute,omitempty" mapstructure:"compute,omitempty"`
	Network Network `json:"network,omitempty" mapstructure:"network,omitempty"`
}

type Compute struct {
	AzEnvironment              string              `json:"azEnvironment,omitempty" mapstructure:"azenvironment,omitempty"`
	EvictionPolicy             string              `json:"evictionPolicy,omitempty" mapstructure:"evictionpolicy,omitempty"`
	IsHostCompatibilityLayerVm string              `json:"isHostCompatibilityLayerVm,omitempty" mapstructure:"ishostcompatibilitylayervm,omitempty"`
	LicenseType                string              `json:"licenseType,omitempty" mapstructure:"licensetype,omitempty"`
	Location                   string              `json:"location,omitempty" mapstructure:"location,omitempty"`
	Name                       string              `json:"name,omitempty" mapstructure:"name,omitempty"`
	Offer                      string              `json:"offer,omitempty" mapstructure:"offer,omitempty"`
	OsProfile                  OsProfile           `json:"osProfile,omitempty" mapstructure:"ofprofile,omitempty"`
	OsType                     string              `json:"osType,omitempty" mapstructure:"ostype,omitempty"`
	PlacementGroupId           string              `json:"placementGroupId,omitempty" mapstructure:"placementgroupid,omitempty"`
	Plan                       Plan                `json:"plan,omitempty" mapstructure:"plan,omitempty"`
	PlatformFaultDomain        string              `json:"platformFaultDomain,omitempty" mapstructure:"platformfaultdomain,omitempty"`
	PlatformUpdateDomain       string              `json:"platformUpdateDomain,omitempty" mapstructure:"platformupdatedomain,omitempty"`
	Priority                   string              `json:"priority,omitempty" mapstructure:"priority,omitempty"`
	Provider                   string              `json:"provider,omitempty" mapstructure:"provider,omitempty"`
	PublicKeys                 []*PublicKey        `json:"publicKeys,omitempty" mapstructure:"publickeys,omitempty"`
	Publisher                  string              `json:"publisher,omitempty" mapstructure:"publisher,omitempty"`
	ResourceGroupName          string              `json:"resourceGroupName,omitempty" mapstructure:"resourcegroupname,omitempty"`
	ResourceId                 string              `json:"resourceId,omitempty" mapstructure:"resourceid,omitempty"`
	SecurityProfile            SecurityProfile     `json:"securityProfile,omitempty" mapstructure:"securityprofile,omitempty"`
	Sku                        string              `json:"sku,omitempty" mapstructure:"sku,omitempty"`
	StorageProfile             StorageProfile      `json:"storageProfile,omitempty" mapstructure:"storageprofile,omitempty"`
	SubscriptionId             string              `json:"subscriptionId,omitempty" mapstructure:"subscriptionid,omitempty"`
	Tags                       string              `json:"tags,omitempty" mapstructure:"tags,omitempty"`
	TagsList                   []map[string]string `json:"tagsList,omitempty" mapstructure:"tagslist,omitempty"`
	UserData                   string              `json:"userData,omitempty" mapstructure:"userdata,omitempty"`
	Version                    string              `json:"version,omitempty" mapstructure:"version,omitempty"`
	VmId                       string              `json:"vmId,omitempty" mapstructure:"vmid,omitempty"`
	VmScaleSetName             string              `json:"vmScaleSetName,omitempty" mapstructure:"vmscalesetname,omitempty"`
	VmSize                     string              `json:"vmSize,omitempty" mapstructure:"vmsize,omitempty"`
	Zone                       string              `json:"zone,omitempty" mapstructure:"zone,omitempty"`
}

type OsProfile struct {
	AdminUserName                 string `json:"adminUsername,omitempty" mapstructure:"adminusername,omitempty"`
	ComputerName                  string `json:"computerName,omitempty" mapstructure:"computername,omitempty"`
	DisablePasswordAuthentication string `json:"disablePasswordAuthentication,omitempty" mapstructure:"disablepasswordauthentication,omitempty"`
}

type Plan struct {
	Name      string `json:"name,omitempty" mapstructure:"name,omitempty"`
	Product   string `json:"product,omitempty" mapstructure:"product,omitempty"`
	Publisher string `json:"publisher,omitempty" mapstructure:"publisher,omitempty"`
}

type PublicKey struct {
	KeyData string `json:"keyData,omitempty" mapstructure:"keydata,omitempty"`
	Path    string `json:"path,omitempty" mapstructure:"path,omitempty"`
}

type SecurityProfile struct {
	SecureBootEnabled string `json:"secureBootEnabled,omitempty" mapstructure:"securebootenabled,omitempty"`
	VirtualTpmEnabled string `json:"virtualTpmEnabled,omitempty" mapstructure:"virtualtpmenabled,omitempty"`
}

type StorageProfile struct {
	DataDisks      []*Disk        `json:"dataDisks,omitempty" mapstructure:"datadisks,omitempty"`
	ImageReference ImageReference `json:"imageReference,omitempty" mapstructure:"imagereference,omitempty"`
	OsDisk         Disk           `json:"osDisk,omitempty" mapstructure:"osDisk,omitempty"`
}

type Disk struct {
	Caching                 string            `json:"caching,omitempty" mapstructure:"caching,omitempty"`
	CreateOption            string            `json:"createOption,omitempty" mapstructure:"createoption,omitempty"`
	DiffDiskSettings        map[string]string `json:"diffDiskSettings,omitempty" mapstructure:"diskdiffsettings,omitempty"`
	DiskSizeGB              string            `json:"diskSizeGB,omitempty" mapstructure:"disksizegb,omitempty"`
	EncryptionSettings      map[string]string `json:"encryptionSettings,omitempty" mapstructure:"encryptionsettings,omitempty"`
	Image                   map[string]string `json:"image,omitempty" mapstructure:"image,omitempty"`
	Lun                     string            `json:"lun,omitempty" mapstructure:"lun,omitempty"`
	ManagedDisk             ManagedDisk       `json:"managedDisk,omitempty" mapstructure:"manageddisk,omitempty"`
	Name                    string            `json:"name,omitempty" mapstructure:"name,omitempty"`
	OsType                  string            `json:"osType,omitempty" mapstructure:"ostype,omitempty"`
	Vhd                     map[string]string `json:"vhd,omitempty" mapstructure:"vhd,omitempty"`
	WriteAcceleratorEnabled string            `json:"writeAcceleratorEnabled,omitempty" mapstructure:"writeacceleratorenabled,omitempty"`
}

type ManagedDisk struct {
	Id                 string `json:"id,omitempty" mapstructure:"id,omitempty"`
	StorageAccountType string `json:"storageAccountType,omitempty" mapstructure:"storageaccounttype,omitempty"`
}

type ImageReference struct {
	Id        string `json:"id,omitempty" mapstructure:"id,omitempty"`
	Offer     string `json:"offer,omitempty" mapstructure:"offer,omitempty"`
	Publisher string `json:"publisher,omitempty" mapstructure:"publisher,omitempty"`
	Sku       string `json:"sku,omitempty" mapstructure:"sku,omitempty"`
	Version   string `json:"version,omitempty" mapstructure:"version,omitempty"`
}

type Network struct {
	Interfaces []*Interface `json:"interface,omitempty" mapstructure:"interfaces,omitempty"`
}

type Interface struct {
	Ipv4       Ip     `json:"ipv4,omitempty" mapstructure:"ipv4,omitempty"`
	Ipv6       Ip     `json:"ipv6,omitempty" mapstructure:"ipv6,omitempty"`
	MacAddress string `json:"macAddress,omitempty" mapstructure:"macaddress,omitempty"`
}

type Ip struct {
	Addresses []*Address `json:"ipAddress,omitempty" mapstructure:"ipaddress,omitempty"`
	Subnets   []*Subnet  `json:"subnet,omitempty" mapstructure:"subbet,omitempty"`
}

type Address struct {
	PrivateIp string `json:"privateIpAddress,omitempty" mapstructure:"privateip,omitempty"`
	PublicIp  string `json:"publicIpAddress,omitempty" mapstructure:"publicip,omitempty"`
}

type Subnet struct {
	Address string `json:"address,omitempty" mapstructure:"address,omitempty"`
	Prefix  string `json:"prefix,omitempty" mapstructure:"prefix,omitempty"`
}

func NewAzureMetadata() (*AzureMetadata, error) {
	var err error
	var transport = &http.Transport{Proxy: nil}

	m := &AzureMetadata{}

	client := http.Client{Transport: transport}

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s/metadata/instance", azureApiAddress), nil)
	req.Header.Add("Metadata", "True")

	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("api-version", azureApiVersion)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	log.Print(string(body))

	err = json.Unmarshal(body, m)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return m, nil
}

func (m *AzureMetadata) GetVmUrl() string {
	return path.Join(azurePortalUrl, m.Compute.ResourceId)
}

func (m *AzureMetadata) GetResourceGroupUrl() string {
	return path.Join(
		azurePortalUrl,
		"subscriptions",
		m.Compute.SubscriptionId,
		"resourceGroups",
		m.Compute.ResourceGroupName,
		"overview",
	)
}
