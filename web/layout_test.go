package web

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_markdownFunc(t *testing.T) {
	markdown := markdownTemplateFunc()
	input := `
# Heading
This is a _test_
`
	output := markdown(input)
	expected := template.HTML("<h1 id=\"heading\">Heading</h1>\n\n<p>This is a <em>test</em></p>\n")
	assert.Equal(t, expected, output)
}
