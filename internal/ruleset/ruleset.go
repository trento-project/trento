package ruleset

import (
	"embed"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/pkg/errors"
)

//go:embed rulesets
var ruleSetsFS embed.FS

const rulesetsFolder = "rulesets"

type RuleSet struct {
	EmbeddedFiles []string
	UserFiles     []string
}

func NewRuleSet(userFiles []string) (*RuleSet, error) {
	var r = &RuleSet{}

	embeddedFiles, err := loadEmbeddedFiles()
	if err != nil {
		return r, err
	}

	r.EmbeddedFiles = embeddedFiles
	r.UserFiles = userFiles

	return r, nil
}

func loadEmbeddedFiles() ([]string, error) {
	var embeddedFiles []string

	dirEntries, err := ruleSetsFS.ReadDir(rulesetsFolder)
	if err != nil {
		return embeddedFiles, err
	}

	for _, d := range dirEntries {
		embeddedFiles = append(embeddedFiles, path.Join(rulesetsFolder, d.Name()))
	}

	return embeddedFiles, nil
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

func (r *RuleSet) GetRulesetFiles() []string {
	return append(r.EmbeddedFiles, r.UserFiles...)
}

func (r *RuleSet) GetRulesets() ([]byte, error) {
	var data [][]byte

	for _, path := range r.EmbeddedFiles {
		datum, err := ruleSetsFS.ReadFile(path)
		if err != nil {
			return nil, errors.Wrap(err, "could not read embedded rulesets file")
		}
		data = append(data, datum)
	}

	for _, path := range r.UserFiles {
		datum, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrap(err, "could not read user rulesets file")
		}
		data = append(data, datum)
	}

	first_yaml_file := string(data[0])
	// Appent all next files to the first
	for _, datum := range data[1:] {
		first_yaml_file += trim_yaml_header(string(datum))
	}

	return []byte(first_yaml_file), nil
}
