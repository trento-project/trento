package agent

import (
	"io/ioutil"

	"github.com/aquasecurity/bench-common/check"
	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

type Check func() (*CheckResult, error)

func NewCheck(definitionsPath string) (Check, error) {
	data, err := ioutil.ReadFile(definitionsPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read definitions file")
	}

	return func() (*CheckResult, error) {
		controls, err := check.NewControls(data, nil)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse definitions file")
		}

		summary := controls.RunGroup()

		out, err := controls.JSON()
		if err != nil {
			return nil, errors.Wrap(err, "could not convert check results to JSON")
		}

		result := &CheckResult{
			Output: out,
		}
		switch true {
		case summary.Fail > 0:
			result.Status = consul.HealthCritical
		case summary.Warn > 0:
			result.Status = consul.HealthWarning
		default:
			result.Status = consul.HealthPassing
		}

		return result, nil
	}, nil
}

type CheckResult struct {
	Output []byte
	Status string
}

func (r *CheckResult) String() string {
	return string(r.Output)
}
