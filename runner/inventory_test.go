package runner

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"

	apiMocks "github.com/trento-project/trento/api/mocks"
	webApi "github.com/trento-project/trento/web"
	"github.com/trento-project/trento/web/models"
)

type InventoryTestSuite struct {
	suite.Suite
}

func TestInventoryTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryTestSuite))
}

func (suite *InventoryTestSuite) Test_CreateInventory() {
	tmpDir, _ := ioutil.TempDir(os.TempDir(), "trentotest")
	destination := path.Join(tmpDir, "ansible_hosts")

	content := &InventoryContent{
		Nodes: []*Node{
			&Node{
				Name:        "node1",
				AnsibleHost: "192.168.10.1",
				AnsibleUser: "trento",
				Variables: map[string]interface{}{
					"key1": "value1",
					"key2": []string{"value2", "value3"},
				},
			},
			&Node{
				Name: "node2",
			},
		},
		Groups: []*Group{
			&Group{
				Name: "group1",
				Nodes: []*Node{
					{
						Name:        "node3",
						AnsibleHost: "192.168.11.1",
						AnsibleUser: "trento",
						Variables: map[string]interface{}{
							"key1": 1,
							"key2": []string{"value2", "value3"},
						},
					},
					&Node{
						Name: "node4",
					},
				},
			},
			&Group{
				Name: "group2",
				Nodes: []*Node{
					{
						Name: "node5",
					},
					&Node{
						Name: "node6",
					},
				},
			},
		},
	}

	err := CreateInventory(destination, content)

	suite.NoError(err)
	suite.FileExists(destination)

	// Cannot use backticks as the lines have a final space in many lines
	expectedContent := "\n" +
		"node1 ansible_host=192.168.10.1 ansible_user=trento key1=value1 key2=[value2 value3] \n" +
		"node2 ansible_host= ansible_user= \n" +
		"[group1]\n" +
		"node3 ansible_host=192.168.11.1 ansible_user=trento key1=1 key2=[value2 value3] \n" +
		"node4 ansible_host= ansible_user= \n" +
		"[group2]\n" +
		"node5 ansible_host= ansible_user= \n" +
		"node6 ansible_host= ansible_user= \n"

	data, err := ioutil.ReadFile(destination)
	if err == nil {
		suite.Equal(expectedContent, string(data))
	}
}

func (suite *InventoryTestSuite) Test_NewClusterInventoryContent() {
	apiInst := new(apiMocks.TrentoApiService)

	apiInst.On("GetClustersSettings").Return(mockedClustersSettings(), nil)

	content, err := NewClusterInventoryContent(apiInst)

	expectedContent := &InventoryContent{
		Groups: []*Group{
			&Group{
				Name: "cluster1",
				Nodes: []*Node{
					&Node{
						Name: "node1",
						Variables: map[string]interface{}{
							"cluster_selected_checks": "[\"check1\",\"check2\"]",
						},
						AnsibleHost: "192.168.10.1",
						AnsibleUser: "user1",
					},
					&Node{
						Name: "node2",
						Variables: map[string]interface{}{
							"cluster_selected_checks": "[\"check1\",\"check2\"]",
						},
						AnsibleHost: "192.168.10.2",
						AnsibleUser: "user2",
					},
				},
			},
			&Group{
				Name: "cluster2",
				Nodes: []*Node{
					&Node{
						Name: "node3",
						Variables: map[string]interface{}{
							"cluster_selected_checks": "[\"check3\",\"check4\"]",
						},
						AnsibleHost: "192.168.10.3",
						AnsibleUser: "clouduser",
					},
					&Node{
						Name: "node4",
						Variables: map[string]interface{}{
							"cluster_selected_checks": "[\"check3\",\"check4\"]",
						},
						AnsibleHost: "",
						AnsibleUser: "root",
					},
				},
			},
		},
	}

	suite.NoError(err)
	suite.ElementsMatch(expectedContent.Groups, content.Groups)
	apiInst.AssertExpectations(suite.T())
}

func mockedClustersSettings() webApi.ClustersSettingsResponse {
	return webApi.ClustersSettingsResponse{
		{
			ID:             "cluster1",
			SelectedChecks: []string{"check1", "check2"},
			Hosts: []*models.HostConnection{
				{
					Name:    "node1",
					Address: "192.168.10.1",
					User:    "user1",
				},
				{
					Name:    "node2",
					Address: "192.168.10.2",
					User:    "user2",
				},
			},
		},
		{
			ID:             "cluster2",
			SelectedChecks: []string{"check3", "check4"},
			Hosts: []*models.HostConnection{
				{
					Name:    "node3",
					Address: "192.168.10.3",
					User:    "clouduser",
				},
				{
					Name:    "node4",
					Address: "",
					User:    "root",
				},
			},
		},
	}
}
