package cli_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fredbi/go-cli/config"
)

func ExampleLoad() {
	here, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv("CONFIG_DIR", filepath.Join(here, "config", "examples"))

	cfg, err := config.Load("dev", config.WithMute(true))
	if err != nil {
		log.Fatal(fmt.Errorf("loading config: %w", err))
	}

	if cfg == nil {
		log.Fatal("config not found")
	}

	fmt.Println("loaded")

	// Output:
	// loaded
}
