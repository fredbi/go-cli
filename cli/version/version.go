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
		if buildRevision := lookupBuildSettings(buildInfo.Settings, "vcs.revision"); buildRevision != "" {
			buildInfo.Commit = buildRevision
		} else {
			buildInfo.Commit = fmt.Sprintf("unknown, mod sum: %q", goInfo.Main.Sum)
		}

		if buildTime := lookupBuildSettings(buildInfo.Settings, "vcs.time"); buildTime != "" {
			buildInfo.Date = buildTime
		} else {
			buildInfo.Date = "(unknown)"
		}
	} else {
		buildInfo.Version = buildVersion
		buildInfo.Commit = buildCommit
		buildInfo.Date = buildDate
	}

	return buildInfo
}

func lookupBuildSettings(settings []debug.BuildSetting, key string) string {
	/*
		Typical settings:
				build	vcs.revision=35798d3d1de0ebb1b241059f74d31486c745c7a4
				build	vcs.time=2023-09-25T13:22:02Z
				build	vcs.modified=true
	*/
	for _, setting := range settings {
		if setting.Key == key {
			return setting.Value
		}
	}

	return ""
}
