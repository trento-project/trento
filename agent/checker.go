package agent

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/aquasecurity/bench-common/check"
	"github.com/pkg/errors"
)

type Checker func() (CheckResult, error)

func NewChecker(definitionsPath string) (Checker, error) {
	data, err := ioutil.ReadFile(definitionsPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read definitions file")
	}

	return func() (CheckResult, error) {
		result := CheckResult{}
		controls, err := check.NewControls(data, nil)
		if err != nil {
			return result, errors.Wrap(err, "could not parse definitions file")
		}

		controls.RunGroup()

		result.controls = controls

		return result, nil
	}, nil
}

type CheckResult struct {
	controls *check.Controls
}

func (cr *CheckResult) CheckPrettyPrint(sw io.StringWriter) {
	sw.WriteString("== Summary ==\n")
	sw.WriteString(strconv.Itoa(cr.controls.Summary.Pass) + " checks PASS\n")
	sw.WriteString(strconv.Itoa(cr.controls.Summary.Fail) + " checks FAIL\n")
	sw.WriteString(strconv.Itoa(cr.controls.Summary.Warn) + " checks WARN\n")
	sw.WriteString(strconv.Itoa(cr.controls.Summary.Info) + " checks INFO\n")
}

func (r CheckResult) Summary() check.Summary {
	return r.controls.Summary
}

func (r CheckResult) MarshalJSON() ([]byte, error) {
	out, err := r.controls.JSON()
	if err != nil {
		return out, errors.Wrap(err, "could not convert check results to JSON")
	}
	return out, nil
}

func (r CheckResult) String() string {
	summary := r.controls.Summary
	return fmt.Sprintf("== Summary ==\n%d checks PASS\n%d checks FAIL\n%d checks WARN\n%d checks INFO\n",
		summary.Pass, summary.Fail, summary.Warn, summary.Info,
	)
}
