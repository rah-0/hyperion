package config

var (
	Loaded    Config
	Path      string
	ForceHost string
)

type Config struct {
	ClusterName string
	Nodes       []struct {
		Host struct {
			Name string
			IP   string
			Port int
		}
		Path struct {
			Data string // Where data will be stored
		}
		Entities []struct {
			Name string
		}
	}
}
