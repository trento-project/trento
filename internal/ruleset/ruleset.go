package ruleset

import (
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/pkg/errors"
)

//go:embed rulesets
var ruleSetsFS embed.FS

const rulesetsFolder = "rulesets"

const (
	Embedded = iota
	User
)

const defaultRuleset = `
controls:
version: 0.0.0
id: 0
description: "Default ruleset. Only used if other rulesets are not selected"
type: "master"
groups:
  - id: 0.1
    description: "Default"
    checks:
      - id: 0.1.1
        description: "Default"
        audit: 'ps -ef | grep trento'
        tests:
          test_items:
            - flag: trento
        remediation: |
          ## Remediation
          Nothing to be fixed. Select other rulesets for advanced checking
        scored: true`

type RuleSets []*RuleSet

type RuleSet struct {
	Path    string `mapstructure:"path,omitempty"`
	Enabled bool   `mapstructure:"enabled"`
	Type    int    `mapstructure:"type"`
}

func NewRuleSets(userFiles []string) (RuleSets, error) {
	var rsets = RuleSets{}

	embeddedRulesets, err := loadEmbeddedFiles()
	if err != nil {
		return rsets, err
	}

	rsets = append(rsets, embeddedRulesets...)

	for _, d := range userFiles {
		userRuleset := &RuleSet{
			Path:    path.Join(d),
			Enabled: false,
			Type:    User,
		}
		rsets = append(rsets, userRuleset)
	}

	return rsets, nil
}

func loadEmbeddedFiles() (RuleSets, error) {
	var embeddedRulesets = RuleSets{}

	dirEntries, err := ruleSetsFS.ReadDir(rulesetsFolder)
	if err != nil {
		return embeddedRulesets, err
	}

	for _, d := range dirEntries {
		r := &RuleSet{
			Path:    path.Join(rulesetsFolder, d.Name()),
			Enabled: false,
			Type:    Embedded,
		}
		embeddedRulesets = append(embeddedRulesets, r)
	}

	return embeddedRulesets, nil
}

func (r RuleSets) GetEnabled() []string {
	var files []string

	for _, ruleSet := range r {
		if ruleSet.Enabled {
			files = append(files, ruleSet.Path)
		}
	}

	return files
}

func (r RuleSets) Enable(paths []string) error {
	var found bool = false

	for _, path := range paths {
		for _, ruleSet := range r {
			if ruleSet.Path == path {
				ruleSet.Enabled = true
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("no ruleset found with the given path: %s", path)
		}
	}

	return nil
}

/* Remove
  '---
   controls:
	version: ...
	id: ...
	description: ...
	type: ...
	groups:
'
that comes before '    - id:...'
You need to do it to append the yaml to another yaml file
*/
func trim_yaml_header(yaml_text string) string {
	pattern := "\ngroups:"
	index := strings.Index(yaml_text, pattern)
	if index == -1 {
		log.Fatal("the yaml file doesn't contain the 'id:' field and is possibly broken")
	}
	index += len(pattern)
	return yaml_text[index:]
}

func (r RuleSets) GetRulesetsYaml(onlyEnabled bool) ([]byte, error) {
	var data [][]byte

	for _, rset := range r {
		if onlyEnabled && !rset.Enabled {
			continue
		}
		if rset.Type == Embedded {
			datum, err := ruleSetsFS.ReadFile(rset.Path)
			if err != nil {
				return nil, errors.Wrap(err, "could not read embedded rulesets file")
			}
			data = append(data, datum)
		} else if rset.Type == User {
			datum, err := ioutil.ReadFile(rset.Path)
			if err != nil {
				return nil, errors.Wrap(err, "could not read user rulesets file")
			}
			data = append(data, datum)
		}
	}

	if len(data) == 0 {
		return []byte(""), nil
	}

	first_yaml_file := string(data[0])
	// Appent all next files to the first
	for _, datum := range data[1:] {
		first_yaml_file += trim_yaml_header(string(datum))
	}

	return []byte(first_yaml_file), nil
}

func GetDefaultYaml() []byte {
	return []byte(defaultRuleset)
}
