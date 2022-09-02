package steps
import (
	"bjsh/installk8s/cluster"
    //"fmt"
	//"path/filepath"
	//"bjsh/installk8s/utils"
)


type Copyfiles struct {
	Config *cluster.Cluster
}


/*
func CreateDirs(dataDir string) error {
	dataDir, err := filepath.Abs(dataDir)
	if err != nil {
		return err
	}
    certDir := formatPath(dataDir, "pki")
	runDir := formatPath(dataDir, "run")
k8sdir := cluster.K8sdir {
	AdminKubeConfigPath:        formatPath(certDir, "admin.conf"),
	BinDir:                     formatPath(dataDir, "bin"),
	OCIBundleDir:               formatPath(dataDir, "images"),
	CertRootDir:                certDir,
	DataDir:                    dataDir,
	EtcdCertDir:                formatPath(certDir, "etcd"),
	EtcdDataDir:                formatPath(dataDir, "etcd"),
	KubeletAuthConfigPath:      formatPath(dataDir, "kubelet.conf"),
	KubeletBootstrapConfigPath: formatPath(dataDir, "kubelet-bootstrap.conf"),
	KubeletVolumePluginDir:     utils.KubeletVolumePluginDir,//这个常量同样见于上面的定义
	ManifestsDir:               formatPath(dataDir, "manifests"),
	RunDir:                     runDir,
} 



}


func formatPath(dir string, file string) string {
	return fmt.Sprintf("%s/%s", dir, file)
}
*/