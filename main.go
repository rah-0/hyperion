package main

import (
	"flag"
	"os"

	"github.com/rah-0/nabu"
)

var pathConfig string

func main() {
	if err := checkPathConfig(); err != nil {
		nabu.FromError(err).WithLevelFatal().Log()
		os.Exit(1)
	}
}

func checkPathConfig() error {
	flag.StringVar(&pathConfig, "pathConfig", "", "")
	flag.Parse()

	if pathConfig == "" {
		pathConfig = getEnvKeyValue("HyperionPathConfig")
	}
	if pathConfig == "" {
		return ErrPathConfigNotSpecified
	}

	exists, err := pathExists(pathConfig)
	if err != nil {
		return err
	}
	if !exists {
		return ErrPathConfigNotFound
	}

	return nil
}
