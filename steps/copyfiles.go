package steps

import (
	"bjsh/installk8s/cluster"
	"bjsh/installk8s/utils"
	"fmt"
	"os"
	"path/filepath"
)


type Copyfiles struct {
	Config *cluster.Cluster
}

func(c Copyfiles) Run() error{
  
}

func GetDirs(dataDir string) error {
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
	KubeletVolumePluginDir:     utils.KubeletVolumePluginDir,
	ManifestsDir:               formatPath(dataDir, "manifests"),
	RunDir:                     runDir,
} 
if !utils.DirExist(dataDir) {
	err := os.Mkdir(dataDir,utils.DefaultFileMode)
	if err != nil {
		return err
	}
}
var dirs []string
dirs=append(dirs,k8sdir.BinDir,k8sdir.CertRootDir,k8sdir.EtcdCertDir,k8sdir.EtcdDataDir,k8sdir.KubeletVolumePluginDir)
for _,dir := range dirs {
	if !utils.DirExist(dir) {
		err := os.Mkdir(dir,utils.DefaultFileMode)
		if err != nil {
			return err
		}
	}
}
return nil
}


func formatPath(dir string, file string) string {
	return fmt.Sprintf("%s/%s", dir, file)
}
