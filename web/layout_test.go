package web

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_markdownToHTML(t *testing.T) {
	input := `
# Heading
This is a _test_
`
	output := markdownToHTML(input)
	expected := template.HTML("<h1 id=\"heading\">Heading</h1>\n\n<p>This is a <em>test</em></p>\n")
	assert.Equal(t, expected, output)
}
