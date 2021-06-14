package agent

import (
	"fmt"
	"io"
	"strconv"

	"github.com/aquasecurity/bench-common/check"
	"github.com/pkg/errors"
)

type Checker func() (CheckResult, error)

func NewChecker(rulesetsData []byte) (Checker, error) {
	return func() (CheckResult, error) {
		var result CheckResult

		controls, err := check.NewControls(rulesetsData, nil)
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
	// 1) Print descriptions
	controls := cr.controls
	remedations := ""
	sw.WriteString("[INFO] " + controls.ID + " " + controls.Description + "\n")
	for _, group := range controls.Groups {
		sw.WriteString("[INFO] " + group.ID + " " + group.Description + "\n")
		for _, check := range group.Checks {
			sw.WriteString("[" + string(check.State) + "] " + check.ID + " " + check.Description + "\n")
			if string(check.State) == "FAIL" {
				remedations += check.ID + " " + check.Remediation + "\n"
			}
		}
	}
	// 2) Print remedations
	sw.WriteString("\n== Remedations ==\n" + remedations)

	// 3) Print the summary
	sw.WriteString("\n== Summary ==\n")
	sw.WriteString(strconv.Itoa(controls.Summary.Pass) + " checks PASS\n")
	sw.WriteString(strconv.Itoa(controls.Summary.Fail) + " checks FAIL\n")
	sw.WriteString(strconv.Itoa(controls.Summary.Warn) + " checks WARN\n")
	sw.WriteString(strconv.Itoa(controls.Summary.Info) + " checks INFO\n")
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
