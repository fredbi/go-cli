package config_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fredbi/go-cli/config"
)

func ExampleLoadWithSecrets() {
	// loads a config, merge clear-text secrets, then save the result file.

	var err error
	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	folder, err := os.MkdirTemp("", "")
	if err != nil {
		err = fmt.Errorf("creating temp dir: %w", err)

		return
	}

	defer func() {
		_ = os.RemoveAll(folder)
	}()

	here, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("get current working dir: %w", err)

		return
	}

	os.Setenv("CONFIG_DIR", filepath.Join(here, "examples"))

	// load and merge configuration files
	cfg, err := config.LoadWithSecrets("dev", config.WithMute(true))
	if err != nil {
		err = fmt.Errorf("loading config: %w", err)

		return
	}

	// writes down the merged config
	configmap := filepath.Join(folder, "configmap.yaml")
	err = cfg.WriteConfigAs(configmap)
	if err != nil {
		err = fmt.Errorf("writing config: %w", err)

		return
	}

	result, err := os.ReadFile(configmap)
	if err != nil {
		err = fmt.Errorf("reading result: %w", err)

		return
	}

	_, _ = os.Stdout.Write(result)

	// Output:
	// app:
	//     threads: 10
	//     url: https://example.dev.co
	// log:
	//     level: info
	// metrics:
	//     enabled: true
	//     exporter: prometheus
	// secrets:
	//     token: xyz
	// trace:
	//     enabled: true
	//     exporter: jaeger
}
