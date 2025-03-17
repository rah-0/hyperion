package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"sync"

	"github.com/rah-0/nabu"

	. "github.com/rah-0/hyperion/model"
	"github.com/rah-0/hyperion/register"
	. "github.com/rah-0/hyperion/util"
)

type Config struct {
	ClusterName string
	Nodes       []*Node
}

var (
	GlobalNode   *Node
	GlobalConfig Config
	pathConfig   string
	forceHost    string
)

func main() {
	flag.StringVar(&pathConfig, "pathConfig", "", "")
	flag.StringVar(&forceHost, "forceHost", "", "")
	flag.Parse()

	nabu.SetLogLevel(nabu.LevelDebug)

	if err := checkConfigs(); err != nil {
		nabu.FromError(err).WithLevelFatal().Log()
		os.Exit(1)
	}

	run()
}

func checkConfigs() error {
	if err := checkPathConfig(); err != nil {
		return nabu.FromError(err).Log()
	}
	if err := checkConfig(); err != nil {
		return nabu.FromError(err).Log()
	}
	checkForceHost()
	if err := checkCurrentNodes(); err != nil {
		return nabu.FromError(err).Log()
	}
	return nil
}

func run() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		GlobalNode.Start()
	}()

	wg.Wait()
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

func checkForceHost() {
	if forceHost == "" {
		forceHost = GetEnvKeyValue("HyperionForceHost")
	}
	return
}

func checkConfig() error {
	content, err := FileRead(pathConfig)
	if err != nil {
		return err
	}

	if len(content) == 0 {
		return ErrPathConfigNoContent
	}

	err = json.Unmarshal(content, &GlobalConfig)
	if err != nil {
		return err
	}

	return nil
}

func checkCurrentNodes() error {
	if len(GlobalConfig.Nodes) == 0 {
		return ErrConfigNodesNotFound
	}

	var hostName string
	if forceHost == "" {
		var err error
		hostName, err = os.Hostname()
		if err != nil {
			return err
		}
	} else {
		hostName = forceHost
	}

	for _, node := range GlobalConfig.Nodes {
		if node.Host.Name == hostName {
			GlobalNode = NewNode(node.Host, node.Path, node.Entities)
		}
	}

	if GlobalNode == nil {
		return ErrConfigNodeNotFoundForHost
	}

	// Config per node targets an entity by name but here we find all versions for that entity
	for _, e := range GlobalNode.Entities {
		for _, re := range register.Entities {
			if e.Name == re.Name {
				GlobalNode.EntitiesStorage = append(GlobalNode.EntitiesStorage, &EntityStorage{
					Disk:   NewDisk().WithPath(filepath.Join(GlobalNode.Path.Data, re.DbFileName)).WithEntity(re),
					Memory: re,
				})
			}
		}
	}

	return nil
}
