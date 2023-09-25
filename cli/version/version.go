package version

import (
	"fmt"
	"runtime/debug"
	"strconv"
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

	// buildSettings collect metadata of interest from the
	// []debug.BuildSetting collection
	buildSettings struct {
		Commit     string
		Date       string
		IsModified bool
	}
)

var (
	buildGoVersion = "unknown"
	buildVersion   = "master"
	buildCommit    = "?"
	buildDate      = ""
)

// Resolve build information.
//
// There are 2 possible sources to resolve version information.
//
// 1. It may be populated by your build system at release build time
// e.g. go build -ldflags="-X 'github.com/fredbi/go-cli/cli/version.buildVersion=v1.0.0'"
//
// For more information, see this excellent tutorial:
// https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
//
// 2. It may be populated by the go runtime, collecting metadata from the build.
func Resolve() BuildInfo {
	var buildInfo BuildInfo

	goInfo, isAvailable := debug.ReadBuildInfo()
	if isAvailable {
		buildInfo.GoVersion = goInfo.GoVersion
	} else {
		buildInfo.GoVersion = buildGoVersion
	}

	if buildDate != "" {
		buildInfo.Version = buildVersion
		buildInfo.Commit = buildCommit
		buildInfo.Date = buildDate

		return buildInfo
	}

	settings := lookupBuildSettings(goInfo.Settings)
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

	return buildInfo
}

func lookupBuildSettings(settings []debug.BuildSetting) buildSettings {
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
