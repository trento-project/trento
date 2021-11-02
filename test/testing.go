package test

import (
	"os"
	"path"
	"runtime"
)

// importing _ "github.com/trento-project/trento/test" in tests would set the cwd to the root of the repo
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
