package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/rah-0/nabu"
)

var (
	nodes      []*Node
	config     Config
	pathConfig string
)

func main() {
	flag.StringVar(&pathConfig, "pathConfig", "", "")
	flag.Parse()

	nabu.SetLogLevel(nabu.LevelDebug)

	if err := run(); err != nil {
		nabu.FromError(err).WithLevelFatal().Log()
		os.Exit(1)
	}
}

func run() error {
	if err := checkPathConfig(); err != nil {
		return nabu.FromError(err).Log()
	}
	if err := checkConfig(); err != nil {
		return nabu.FromError(err).Log()
	}
	if err := checkCurrentNodes(); err != nil {
		return nabu.FromError(err).Log()
	}

	startNodes()
	go func() {
		for _, node := range nodes {
			nabu.FromError(<-node.errCh).Log()
			os.Exit(1)
		}
	}()

	return nil
}

func checkPathConfig() error {
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
	if !fileIsEditable(pathConfig) {
		return ErrPathConfigNotEditable
	}

	return nil
}

func checkConfig() error {
	content, err := fileRead(pathConfig)
	if err != nil {
		return err
	}

	if len(content) == 0 {
		return ErrPathConfigNoContent
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	return nil
}

func checkCurrentNodes() error {
	if len(config.Nodes) == 0 {
		return ErrConfigNoNodes
	}

	hostName, err := os.Hostname()
	if err != nil {
		return err
	}

	for _, node := range config.Nodes {
		if node.Host.Name == hostName {
			nodes = append(nodes, &node)
		}
	}

	if len(nodes) == 0 {
		return ErrConfigCurrentNodesNotFound
	}

	return nil
}
