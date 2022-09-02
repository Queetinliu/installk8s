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

var (
	etcdhostlist []string
    apiserverhostlist []string
	controllerhostlist []string
	schedulerhostlist []string
)
etcdhostlist = g.Config.Hosts.AppendHostlist("etcd",etcdhostlist)
etcdhostlist=append(etcdhostlist,"127.0.0.1")
etcdReq := utils.Request{
	Name: "etcd",
	CN: "etcd",
	O: "k8s",
	CAKey: caCertKey,
	CACert: caCertPath,
	Hostnames: etcdhostlist,
}
apiserverhostlist = g.Config.Hosts.AppendHostlist("kube-apiserver",apiserverhostlist)
apiserverhostlist = append(apiserverhostlist, "127.0.0.1","10.88.0.1","kubernetes","kubernetes.default","kubernetes.default.svc","kubernetes.default.svc.cluster","kubernetes.default.svc.cluster.local.")
apiserverReq := utils.Request{
	Name: "kube-apiserver",
	CN: "kubernetes-master",
	O: "k8s",
	CAKey: caCertKey,
	CACert: caCertPath,
	Hostnames: apiserverhostlist,
}
metricsserverReq := utils.Request{
	Name: "proxy-client",
	CN: "aggregator",
	O: "k8s",
	CAKey: caCertKey,
	CACert: caCertPath,
	Hostnames: []string{},
}
controllerhostlist = g.Config.Hosts.AppendHostlist("kube-controller-manager",controllerhostlist)
controllerhostlist = append(controllerhostlist, "127.0.0.1")
controllerReq := utils.Request{
	Name: "kube-controller-manager",
	CN: "system:kube-controller-manager",
	O: "system:kube-controller-manager",
	CAKey: caCertKey,
	CACert: caCertPath,
	Hostnames: controllerhostlist,
}
schedulerhostlist = g.Config.Hosts.AppendHostlist("kube-scheduler",schedulerhostlist)
schedulerhostlist = append(schedulerhostlist, "127.0.0.1")
schedulerReq := utils.Request{
	Name: "kube-scheduler",
	CN: "system:kube-scheduler",
	O: "system:kube-scheduler",
	CAKey: caCertKey,
	CACert: caCertPath,
	Hostnames: schedulerhostlist,
}
kubeproxyReq := utils.Request{
	Name: "kube-proxy",
	CN: "system:kube-proxy",
	O: "k8s",
	CAKey: caCertKey,
	CACert: caCertPath,
	Hostnames: []string{},
}
var reqs []utils.Request
reqs=append(reqs,adminReq,etcdReq,apiserverReq,metricsserverReq,controllerReq,schedulerReq,kubeproxyReq)
var wg sync.WaitGroup
var errors []string

ec := make(chan error)
for _,req := range reqs {
	wg.Add(1)
	fmt.Println(req.Hostnames)
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


