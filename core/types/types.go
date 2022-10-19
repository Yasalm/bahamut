package types

type (
	RestartPolicy string
	MakerType     string

	BahamutSettings map[MakerType]int

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
		// OfType        MakerType
	}
)

const (
	OnFailure     RestartPolicy = "on-failure"
	Always        RestartPolicy = "always"
	UnlessStopped RestartPolicy = "unless-stopped"
)

const (
	General       MakerType = "general"
	Database      MakerType = "database"
	DataIngestion MakerType = "data-ingestion"
	Baremetal     MakerType = "bare-metal"
	Reporting     MakerType = "reporting"
)

func (opt *Options) SetRestartPolicy() {
}
