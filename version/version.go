package version

import "strings"

var Version string

func GetShortVersion() string {
	return strings.Split(Version, "+")[0]
}
