package main

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"sync"
	"syscall"

	"github.com/rah-0/nabu"
	"github.com/rah-0/parsort"

	"github.com/rah-0/hyperion/config"
	"github.com/rah-0/hyperion/node"
	"github.com/rah-0/hyperion/profiler"
	"github.com/rah-0/hyperion/template"

	"github.com/rah-0/hyperion/model"
	"github.com/rah-0/hyperion/util"
)

func main() {
	// Basic configuration flags
	flag.StringVar(&config.Path, "pathConfig", "", "Path to configuration file")
	flag.StringVar(&config.ForceHost, "forceHost", "", "Force specific hostname")
	// Profiler flags
	flag.BoolVar(&config.ProfilerEnabled, "profiler", false, "Enable profiler")
	flag.StringVar(&config.ProfilerIP, "profiler-ip", "0.0.0.0", "IP to bind profiler (default: 0.0.0.0 if profiler enabled)")
	flag.IntVar(&config.ProfilerPort, "profiler-port", 6060, "Port for profiler (default: 6060 if profiler enabled)")
	flag.Parse()

	nabu.SetLogLevel(nabu.LevelDebug)
	if 1 == 2 {
		parsort.TuneSpecific(1000, 1000, 2000, -25, false)
	}

	n, err := checkConfigs()
	if err != nil {
		nabu.FromError(err).WithLevelFatal().Log()
		os.Exit(1)
	}

	if err = template.RegisterEntities(); err != nil {
		nabu.FromError(err).WithLevelFatal().Log()
		os.Exit(1)
	}

	startProfilerIfEnabled()
	run(n)
}

func checkConfigs() (*node.Node, error) {
	if err := checkPathConfig(); err != nil {
		return nil, nabu.FromError(err).Log()
	}
	if err := checkConfig(); err != nil {
		return nil, nabu.FromError(err).Log()
	}
	checkForceHost()

	n, err := checkCurrentNode()
	if err != nil {
		return nil, nabu.FromError(err).Log()
	}

	return n, nil
}

func run(n *node.Node) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := n.Start(); err != nil {
			if errors.Is(err, syscall.EADDRINUSE) {
				nabu.FromError(err).WithLevelFatal().Log()
				os.Exit(1)
			} else {
				nabu.FromError(err).Log()
			}
		}
	}()

	wg.Wait()
}

func checkPathConfig() error {
	if config.Path == "" {
		config.Path = util.GetEnvKeyValue("HyperionPathConfig")
	}
	if config.Path == "" {
		return model.ErrPathConfigNotSpecified
	}

	exists, err := util.PathExists(config.Path)
	if err != nil {
		return err
	}
	if !exists {
		return model.ErrPathConfigNotFound
	}
	if !util.FileIsEditable(config.Path) {
		return model.ErrPathConfigNotEditable
	}

	return nil
}

func checkForceHost() {
	if config.ForceHost == "" {
		config.ForceHost = util.GetEnvKeyValue("HyperionForceHost")
	}
	return
}

func checkConfig() error {
	content, err := util.FileRead(config.Path)
	if err != nil {
		return err
	}

	if len(content) == 0 {
		return model.ErrPathConfigNoContent
	}

	err = json.Unmarshal(content, &config.Loaded)
	if err != nil {
		return err
	}

	return nil
}

func checkCurrentNode() (*node.Node, error) {
	if len(config.Loaded.Nodes) == 0 {
		return nil, model.ErrConfigNodesNotFound
	}

	hostName := config.ForceHost
	if hostName == "" {
		h, err := os.Hostname()
		if err != nil {
			return nil, err
		}
		hostName = h
	}

	for _, nodeConfig := range config.Loaded.Nodes {
		if nodeConfig.Host.Name == hostName {
			n := node.NewNode().
				WithHost(nodeConfig.Host.Name, nodeConfig.Host.IP, nodeConfig.Host.Port).
				WithPath(nodeConfig.Path.Data)

			for _, e := range nodeConfig.Entities {
				n.AddEntity(e.Name)
			}

			addNodePeers(n, config.Loaded)
			return n, nil
		}
	}

	return nil, model.ErrConfigNodeNotFoundForHost
}

func addNodePeers(n *node.Node, c config.Config) {
	for _, nc := range c.Nodes {
		// Skip self
		if n.Host.Name == nc.Host.Name &&
			n.Host.IP == nc.Host.IP &&
			n.Host.Port == nc.Host.Port {
			continue
		}

		peer := node.NewNode().
			WithHost(nc.Host.Name, nc.Host.IP, nc.Host.Port).
			WithPath(nc.Path.Data)

		for _, e := range nc.Entities {
			peer.AddEntity(e.Name)
		}

		n.AddPeer(peer)
	}
}

func startProfilerIfEnabled() {
	if config.ProfilerEnabled {
		profiler.Start(config.ProfilerIP, config.ProfilerPort)
	}
}
