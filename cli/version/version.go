package version

import (
	"fmt"
	"runtime/debug"

	"github.com/davecgh/go-spew/spew"
)

// BuildInfo holds versioning information
type BuildInfo struct {
	GoVersion string               `json:"goVersion"`
	Version   string               `json:"version"`
	Commit    string               `json:"commit"`
	Date      string               `json:"date"`
	Settings  []debug.BuildSetting `json:"-"`
}

var (
	// Populated by your build system at release build time
	buildGoVersion = "unknown"
	buildVersion   = "master"
	buildCommit    = "?"
	buildDate      = ""
)

func Resolve() BuildInfo {
	var buildInfo BuildInfo

	goInfo, isAvailable := debug.ReadBuildInfo()
	if isAvailable {
		buildInfo.GoVersion = goInfo.GoVersion
		spew.Dump("goInfo: %#v\n", goInfo)
	} else {
		buildInfo.GoVersion = buildGoVersion
	}

	if buildDate == "" {
		buildInfo.Version = goInfo.Main.Version
		buildInfo.Commit = fmt.Sprintf("unknown, mod sum: %q", goInfo.Main.Sum)
		buildInfo.Date = "(unknown)"
	} else {
		buildInfo.Version = buildVersion
		buildInfo.Commit = buildCommit
		buildInfo.Date = buildDate
	}

	return buildInfo
}
