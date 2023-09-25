package version

import (
	"fmt"
	"runtime/debug"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

type (
	// BuildInfo holds versioning information
	BuildInfo struct {
		GoVersion  string               `json:"goVersion"`
		Version    string               `json:"version"`
		Commit     string               `json:"commit"`
		Date       string               `json:"date"`
		IsModified bool                 `json:"isModified"`
		Settings   []debug.BuildSetting `json:"-"`
	}

	buildSettings struct {
		Commit     string
		Date       string
		IsModified bool
	}
)

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
		spew.Dump("goModuleInfo: %#v\n", goInfo.Main)
	} else {
		buildInfo.GoVersion = buildGoVersion
	}

	if buildDate == "" {
		settings := lookupBuildSettings(buildInfo.Settings)

		buildInfo.Version = goInfo.Main.Version
		if settings.Commit != "" {
			buildInfo.Commit = settings.Commit
			buildInfo.IsModified = settings.IsModified
		} else {
			buildInfo.Commit = fmt.Sprintf("unknown, mod sum: %q", goInfo.Main.Sum)
		}

		if settings.Date != "" {
			buildInfo.Date = settings.Date
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

func lookupBuildSettings(settings []debug.BuildSetting) buildSettings {
	/*
		Typical settings:
				build	vcs.revision=35798d3d1de0ebb1b241059f74d31486c745c7a4
				build	vcs.time=2023-09-25T13:22:02Z
				build	vcs.modified=true
	*/
	output := buildSettings{}

	for _, setting := range settings {
		switch setting.Key {
		case "vcs.time":
			output.Date = setting.Value
		case "vcs.revision":
			output.Commit = setting.Value
		case "vcs.modified":
			isModified, _ := strconv.ParseBool(setting.Value)
			output.IsModified = isModified
		}
	}

	return output
}
