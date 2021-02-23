//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=webapp/templates

package main

import "github.com/SUSE/console-for-sap/cmd"

func main() {
	cmd.Execute()
}
