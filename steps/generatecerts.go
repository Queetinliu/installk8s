package steps
import (
	"bjsh/installk8s/cluster"
	"bjsh/installk8s/utils"
	"path/filepath"
	"github.com/sirupsen/logrus"
	"fmt"
	"os"
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
_,err = utils.EnsureCertificate(CertRootDir,adminReq)
if err != nil {
	return err
}
/*
etcdReq := utils.Request{
	Name: "etcd",
	CN: "etcd",
	O: "k8s",
	CAKey: caCertKey,
	CACert: caCertPath,
	Hostnames: []string{},
} 
*/

return nil
}
