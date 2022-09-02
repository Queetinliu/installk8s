package steps
import (
	"bjsh/installk8s/cluster"
)

type Disconnect struct {
	Config *cluster.Cluster
}

func(d Disconnect) Run() error {
	return d.Config.Hosts.ParallelEach(func(h *cluster.Host) error {
		h.DisSshconnect()
        h.DisSftpconnect()
		return nil
	})




}