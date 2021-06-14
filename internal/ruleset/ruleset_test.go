package ruleset

import (
	"io/ioutil"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRuleSet(t *testing.T) {
	r, err := NewRuleSet([]string{})

	assert.ElementsMatch(t, []string{"rulesets/1-azure-rules.yaml", "rulesets/2-azure-rules-perf-optimized.yaml"}, r.EmbeddedFiles)
	assert.NoError(t, err)
}

func TestNewRuleSetUserFiles(t *testing.T) {
	r, err := NewRuleSet([]string{"ruleset1", "ruleset2"})

	assert.ElementsMatch(t, []string{"rulesets/1-azure-rules.yaml", "rulesets/2-azure-rules-perf-optimized.yaml"}, r.EmbeddedFiles)
	assert.ElementsMatch(t, []string{"ruleset1", "ruleset2"}, r.UserFiles)

	assert.NoError(t, err)
}

func TestGetRulesets(t *testing.T) {
	rulesetData := `
controls:
version: 1.0.0
id: 9999
description: "My dummy Test description"
type: "master"
groups:
  - id: 9999.1
    description: "My dummy Test"
    checks:
      - id: 9999.1.1
        description: "My test"
        audit: 'do something'
        tests:
          test_items:
            - flag: 30000
        remediation: |
          ## Remediation
          It always passes
        scored: true`

	// create test file
	f, err := ioutil.TempFile("", "rulesets.yaml")
	defer syscall.Unlink(f.Name())

	ioutil.WriteFile(f.Name(), []byte(rulesetData), 0644)

	r, _ := NewRuleSet([]string{f.Name()})
	data, err := r.GetRulesets()

	assert.NoError(t, err)
	assert.Contains(t, string(data), "HA Configuration checks for SAP on MS Azure (generic scenario)")
	assert.Contains(t, string(data), "id: 1.1")
	assert.NotContains(t, string(data), "HA Configuration checks for SAP on MS Azure (Scale-up Performance-Optimized scenario)")
	assert.Contains(t, string(data), "id: 2.1")
	assert.NotContains(t, string(data), "My dummy Test description")
	assert.Contains(t, string(data), "id: 9999.1")
}
