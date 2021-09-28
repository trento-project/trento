package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostNames(t *testing.T) {
	cR := &Results{
		Checks: map[string]*ChecksByHost{
			"1.1.1": &ChecksByHost{
				Hosts: map[string]*Check{
					"host1": &Check{
						Result: "passing",
					},
					"host2": &Check{
						Result: "critical",
					},
				},
			},
			"1.1.2": &ChecksByHost{
				Hosts: map[string]*Check{
					"host1": &Check{
						Result: "warning",
					},
					"host2": &Check{
						Result: "skipped",
					},
				},
			},
		},
	}

	expectedHost := []string{"host1", "host2"}

	assert.ElementsMatch(t, expectedHost, cR.GetHostNames())

}
