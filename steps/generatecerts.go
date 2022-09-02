package steps

import (
	"bjsh/installk8s/cluster"
	"bjsh/installk8s/utils"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"github.com/sirupsen/logrus"
)


type GenerateCerts struct {
	Config *cluster.Cluster
}

func(g GenerateCerts) Run() error{
CertRootDir := utils.GetCertRootDir()
err := utils.EnsureCA(CertRootDir,"ca","kubernetes-ca")
if err != nil {
	return err
}
caCertPath := filepath.Join(CertRootDir, "ca.crt")
caCertKey := filepath.Join(CertRootDir, "ca.key")
logrus.Debugf("CA key and cert exists, loading")
_, err = os.ReadFile(caCertPath)
if err != nil {
	return fmt.Errorf("failed to read ca cert: %w", err)
	}
adminReq := utils.Request{
	Name: "admin",
	CN: "admin",
	O: "system:masters",
	CAKey: caCertKey,
	CACert: caCertPath,
	Hostnames: []string{},
} 
/*
_,err = utils.EnsureCertificate(CertRootDir,adminReq)
if err != nil {
	return err
}

etcdReq := utils.Request{
	Name: "etcd",
	CN: "etcd",
	O: "k8s",
	CAKey: caCertKey,
	CACert: caCertPath,
	Hostnames: []string{},
} 
*/
var etcdhostlist []string
for _,host := range g.Config.Hosts.Etcds() {
	etcdhostlist=append(etcdhostlist,host.Ip)
}
etcdhostlist=append(etcdhostlist,"127.0.0.1")
etcdReq := utils.Request{
	Name: "etcd",
	CN: "etcd",
	O: "k8s",
	CAKey: caCertKey,
	CACert: caCertPath,
	Hostnames: etcdhostlist,
}
var reqs []utils.Request
reqs=append(reqs,adminReq,etcdReq)
var wg sync.WaitGroup
var errors []string
wg.Add(len(reqs))
ec := make(chan error)
for _,req := range reqs {
	go func(r utils.Request) {
		_,err := utils.EnsureCertificate(CertRootDir,r)
		ec <- err
	}(req)

    go func() {
		for err := range ec {
			if err != nil {
				errors=append(errors,err.Error())
			}
			wg.Done()
		}
	}()
    wg.Wait()

}
if len(errors) > 0 {
	return fmt.Errorf("%d certs create failed",len(errors))
}

return nil
}


