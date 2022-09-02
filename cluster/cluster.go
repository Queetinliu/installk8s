package cluster


type Cluster struct {
	Hosts Hosts `yaml:"hosts"`
	}

type K8sdir struct {
	AdminKubeConfigPath        string // The cluster admin kubeconfig location
	BinDir                     string // location for all pki related binaries
	CertRootDir                string // CertRootDir defines the root location for all pki related artifacts
	WindowsCertRootDir         string // WindowsCertRootDir defines the root location for all pki related artifacts
	EtcdCertDir                string // EtcdCertDir contains etcd certificates
	EtcdDataDir                string // EtcdDataDir contains etcd state
	KubeletAuthConfigPath      string // KubeletAuthConfigPath defines the default kubelet auth config path
	KubeletBootstrapConfigPath string // KubeletBootstrapConfigPath defines the default path for kubelet bootstrap auth config
	KubeletVolumePluginDir     string // location for kubelet plugins volume executables
	ManifestsDir               string // location for all stack manifests
	RunDir                     string // location of supervised pid files and sockets
	OCIBundleDir               string // location for OCI bundles
	DataDir                    string // Data directory containing k0s state
}