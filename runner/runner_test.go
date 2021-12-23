package runner

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestAnsibleFolder string = "../test/ansible_test"
)

// TODO: This test could be improved to check the definitve ansible files structure
// once we have something fixed
func TestCreateAnsibleFiles(t *testing.T) {
	tmpDir, _ := ioutil.TempDir(os.TempDir(), "trentotest")
	err := createAnsibleFiles(tmpDir)

	assert.DirExists(t, path.Join(tmpDir, "ansible"))
	assert.NoError(t, err)

	os.RemoveAll(tmpDir)
}

func TestNewAnsibleMetaRunner(t *testing.T) {

	cfg := &Config{
		ApiHost:       "127.0.0.1",
		ApiPort:       8000,
		AnsibleFolder: TestAnsibleFolder,
	}

	a, err := NewAnsibleMetaRunner(cfg)

	expectedMetaRunner := &AnsibleRunner{
		Playbook: path.Join(TestAnsibleFolder, "ansible/meta.yml"),
		Envs: map[string]string{
			"ANSIBLE_CONFIG":      path.Join(TestAnsibleFolder, "ansible/ansible.cfg"),
			"TRENTO_WEB_API_HOST": "127.0.0.1",
			"TRENTO_WEB_API_PORT": "8000",
		},
		Check: false,
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedMetaRunner, a)
}

func TestNewAnsibleCheckRunner(t *testing.T) {

	cfg := &Config{
		ApiHost:       "127.0.0.1",
		ApiPort:       8000,
		AnsibleFolder: TestAnsibleFolder,
	}

	a, err := NewAnsibleCheckRunner(cfg)

	expectedMetaRunner := &AnsibleRunner{
		Playbook: path.Join(TestAnsibleFolder, "ansible/check.yml"),
		Envs: map[string]string{
			"ANSIBLE_CONFIG":      path.Join(TestAnsibleFolder, "ansible/ansible.cfg"),
			"TRENTO_WEB_API_HOST": "127.0.0.1",
			"TRENTO_WEB_API_PORT": "8000",
		},
		Check: true,
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedMetaRunner, a)
}
