package cloud

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trento-project/trento/internal/cloud/mocks"
)

func TestNewAzureMetadata(t *testing.T) {
	clientMock := new(mocks.HTTPClient)

	aFile, _ := os.Open("../../test/azure_metadata")
	bodyText, _ := ioutil.ReadAll(aFile)
	body := ioutil.NopCloser(bytes.NewReader([]byte(bodyText)))

	response := &http.Response{
		StatusCode: 200,
		Body:       body,
	}

	clientMock.On("Do", mock.AnythingOfType("*http.Request")).Return(
		response, nil,
	)

	client = clientMock

	m, err := NewAzureMetadata()

	expectedMeta := &AzureMetadata{
		Compute: Compute{
			AzEnvironment:              "AzurePublicCloud",
			EvictionPolicy:             "",
			IsHostCompatibilityLayerVm: "false",
			LicenseType:                "",
			Location:                   "westeurope",
			Name:                       "vmhana01",
			Offer:                      "sles-sap-15-sp2-byos",
			OsProfile: OsProfile{
				AdminUserName:                 "cloudadmin",
				ComputerName:                  "vmhana01",
				DisablePasswordAuthentication: "true",
			},
			OsType:           "Linux",
			PlacementGroupId: "",
			Plan: Plan{
				Name:      "",
				Product:   "",
				Publisher: "",
			},
			PlatformFaultDomain:  "1",
			PlatformUpdateDomain: "1",
			Priority:             "",
			Provider:             "Microsoft.Compute",

			PublicKeys: []*PublicKey{
				{
					KeyData: "ssh-rsa content\n",
					Path:    "/home/cloudadmin/.ssh/authorized_keys",
				},
			},
			Publisher:         "SUSE",
			ResourceGroupName: "test",
			ResourceId:        "/subscriptions/xxxxx/resourceGroups/test/providers/Microsoft.Compute/virtualMachines/vmhana01",
			SecurityProfile: SecurityProfile{
				SecureBootEnabled: "false",
				VirtualTpmEnabled: "false",
			},
			Sku: "gen2",
			StorageProfile: StorageProfile{
				DataDisks: []*Disk{
					{
						Caching:      "None",
						CreateOption: "Empty",
						DiskSizeGB:   "128",
						Image: map[string]string{
							"uri": "",
						},
						Lun: "0",
						ManagedDisk: ManagedDisk{
							Id:                 "/subscriptions/xxxxx/resourceGroups/test/providers/Microsoft.Compute/disks/disk-hana01-Data01",
							StorageAccountType: "Premium_LRS",
						},
						Name: "disk-hana01-Data01",
						Vhd: map[string]string{
							"uri": "",
						},
						WriteAcceleratorEnabled: "false",
					},
					{
						Caching:      "None",
						CreateOption: "Empty",
						DiskSizeGB:   "128",
						Image: map[string]string{
							"uri": "",
						},
						Lun: "1",
						ManagedDisk: ManagedDisk{
							Id:                 "/subscriptions/xxxxx/resourceGroups/test/providers/Microsoft.Compute/disks/disk-hana01-Data02",
							StorageAccountType: "Premium_LRS",
						},
						Name: "disk-hana01-Data02",
						Vhd: map[string]string{
							"uri": "",
						},
						WriteAcceleratorEnabled: "false",
					},
				},
				ImageReference: ImageReference{
					Id:        "",
					Offer:     "sles-sap-15-sp2-byos",
					Publisher: "SUSE",
					Sku:       "gen2",
					Version:   "latest",
				},
				OsDisk: Disk{
					Caching:      "ReadWrite",
					CreateOption: "FromImage",
					DiffDiskSettings: map[string]string{
						"option": "",
					},
					DiskSizeGB: "30",
					EncryptionSettings: map[string]string{
						"enabled": "false",
					},
					Image: map[string]string{
						"uri": "",
					},
					Lun: "",
					ManagedDisk: ManagedDisk{
						Id:                 "/subscriptions/xxxxx/resourceGroups/test/providers/Microsoft.Compute/disks/disk-hana01-Os",
						StorageAccountType: "Premium_LRS",
					},
					Name:   "disk-hana01-Os",
					OsType: "Linux",
					Vhd: map[string]string{
						"uri": "",
					},
					WriteAcceleratorEnabled: "false",
				},
			},
			SubscriptionId: "xxxxx",
			Tags:           "workspace:xdemo",
			TagsList: []map[string]string{
				map[string]string{
					"name":  "workspace",
					"value": "xdemo",
				},
			},
			UserData:       "",
			Version:        "2021.06.05",
			VmId:           "data",
			VmScaleSetName: "",
			VmSize:         "Standard_E4s_v3",
			Zone:           "",
		},
		Network: Network{
			Interfaces: []*Interface{
				{
					Ipv4: Ip{
						Addresses: []*Address{
							{
								PrivateIp: "10.74.1.10",
								PublicIp:  "1.2.3.4",
							},
						},
						Subnets: []*Subnet{
							{
								Address: "10.74.1.0",
								Prefix:  "24",
							},
						},
					},
					Ipv6: Ip{
						Addresses: []*Address{},
						Subnets:   []*Subnet(nil),
					},
					MacAddress: "000D3A2267C3",
				},
			},
		},
	}

	assert.Equal(t, expectedMeta, m)
	assert.NoError(t, err)
}

func TestGetVmUrl(t *testing.T) {
	meta := &AzureMetadata{
		Compute: Compute{
			ResourceId: "myresourceid",
		},
	}

	assert.Equal(t, "https:/portal.azure.com/#@SUSERDBillingsuse.onmicrosoft.com/resource/myresourceid", meta.GetVmUrl())
}

func TestGetResourceGroupUrl(t *testing.T) {
	meta := &AzureMetadata{
		Compute: Compute{
			SubscriptionId:    "xxx",
			ResourceGroupName: "myresourcegroupname",
		},
	}

	assert.Equal(t, "https:/portal.azure.com/#@SUSERDBillingsuse.onmicrosoft.com/resource/subscriptions/xxx/resourceGroups/myresourcegroupname/overview", meta.GetResourceGroupUrl())
}
