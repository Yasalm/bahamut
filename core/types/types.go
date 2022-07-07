package types

import "fmt"

type (
	RestartPolicy string
	RunnerGroup   string
	Address       string

	BahamutSettings map[RunnerGroup]int

	Options struct {
		Name          string
		AttachStdin   bool
		AttachStout   bool
		AttachStderr  bool
		Cmd           []string
		Image         string
		Memory        int64
		Disk          int64
		Env           []string
		RestartPolicy RestartPolicy
		// OfType        RunnerGroup
	}

	Network struct {
		Host    string
		Port    int
		Binding Address
	}
)

const (
	OnFailure     RestartPolicy = "on-failure"
	Always        RestartPolicy = "always"
	UnlessStopped RestartPolicy = "unless-stopped"
)

const (
	General       RunnerGroup = "general"
	Database      RunnerGroup = "database"
	DataIngestion RunnerGroup = "data-ingestion"
	Baremetal     RunnerGroup = "baremetal"
	Reporting     RunnerGroup = "reporting"
	ModelServing  RunnerGroup = "model-serving"
	ModelQuality  RunnerGroup = "model-quality"
	Python        RunnerGroup = "python"
	Monitoring    RunnerGroup = "monitoring"
	APIGateway    RunnerGroup = "api-gateway"
)

func (n *Network) Address() Address {
	n.Binding = Address(fmt.Sprintf("%s:%v", n.Host, n.Port))
	return n.Binding
}
func (opt *Options) SetRestartPolicy() {
	// TODO: set the restart policy for docker container
}
