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


type RuleSets []*RuleSet

type readFile func(filename string) ([]byte, error)

type RuleSet struct {
	Path     string `mapstructure:"path,omitempty"`
	Enabled  bool   `mapstructure:"enabled"`
	Type     int    `mapstructure:"type"`
	readFile readFile
}



func NewRuleSets(userRuleSetFiles []string) (RuleSets, error) {
	var rsets = RuleSets{}

	embeddedRulesets, err := loadEmbeddedFiles()
	if err != nil {
		return rsets, err
	}

	rsets = append(rsets, embeddedRulesets...)
	rsets = append(rsets, loadUserFiles(userRuleSetFiles)...)

	return rsets, nil
}

func loadUserFiles(userRuleSetFiles []string) RuleSets {
	var userRulesets = RuleSets{}

	for _, d := range userRuleSetFiles {
		userRuleset := &RuleSet{
			Path:     path.Join(d),
			Enabled:  false,
			Type:     User,
			readFile: ioutil.ReadFile,
		}
		userRulesets = append(userRulesets, userRuleset)
	}

	return userRulesets
}

func loadEmbeddedFiles() (RuleSets, error) {
	var embeddedRulesets = RuleSets{}

	dirEntries, err := ruleSetsFS.ReadDir(rulesetsFolder)
	if err != nil {
		return embeddedRulesets, err
	}

	for _, d := range dirEntries {
		r := &RuleSet{
			Path:     path.Join(rulesetsFolder, d.Name()),
			Enabled:  false,
			Type:     Embedded,
			readFile: ruleSetsFS.ReadFile,
		}
		embeddedRulesets = append(embeddedRulesets, r)
	}

	return embeddedRulesets, nil
}

func (r RuleSets) GetEnabled() RuleSets {
	var enabledRuleSets = RuleSets{}

	for _, ruleSet := range r {
		if ruleSet.Enabled {
			enabledRuleSets = append(enabledRuleSets, ruleSet)
		}
	}

	return enabledRuleSets
}

func (r RuleSets) GetPaths() []string {
	var rPaths []string

	for _, ruleSet := range r {
		rPaths = append(rPaths, ruleSet.Path)
	}

	return rPaths
}

func (r RuleSets) Enable(rPaths []string) error {
	var found bool = false

	for _, rPath := range rPaths {
		for _, ruleSet := range r {
			if ruleSet.Path == rPath {
				ruleSet.Enabled = true
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("no ruleset found with the given path: %s", rPath)
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

func (r RuleSets) GetRulesetsYaml() ([]byte, error) {
	var data [][]byte

	for _, rset := range r {
		datum, err := rset.readFile(rset.Path)
		if err != nil {
			return nil, errors.Wrap(err, "could not read rulesets file")
		}
		data = append(data, datum)
	}

	if len(data) == 0 {
		return []byte(""), nil
	}

	first_yaml_file := string(data[0])
	// Append all next files to the first
	for _, datum := range data[1:] {
		first_yaml_file += trim_yaml_header(string(datum))
	}

	return []byte(first_yaml_file), nil
}
