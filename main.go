package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/rah-0/nabu"

	. "github.com/rah-0/hyperion/util"
)

type Config struct {
	ClusterName string
	Nodes       []*Node
}

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

	return nil
}

func checkPathConfig() error {
	if pathConfig == "" {
		pathConfig = GetEnvKeyValue("HyperionPathConfig")
	}
	if pathConfig == "" {
		return ErrPathConfigNotSpecified
	}

	exists, err := PathExists(pathConfig)
	if err != nil {
		return err
	}
	if !exists {
		return ErrPathConfigNotFound
	}
	if !FileIsEditable(pathConfig) {
		return ErrPathConfigNotEditable
	}

	return nil
}

func checkConfig() error {
	content, err := FileRead(pathConfig)
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
		return ErrConfigNodesNotFound
	}

	hostName, err := os.Hostname()
	if err != nil {
		return err
	}

	for _, node := range config.Nodes {
		if node.Host.Name == hostName {
			nodes = append(nodes, NewNode(node.Host, node.Path))
		}
	}

	if len(nodes) == 0 {
		return ErrConfigNodesNotFoundForHost
	}

	return nil
}

func startNodes() {
	for _, node := range nodes {
		go node.Start()
	}
	waitNodesToBeReady()
}

func waitNodesToBeReady() {
	for {
		allReady := true
		for _, node := range nodes {
			node.Mu.Lock()
			status := node.Status
			node.Mu.Unlock()
			if status != NodeStatusReady {
				allReady = false
				break
			}
		}
		if allReady {
			break
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}
