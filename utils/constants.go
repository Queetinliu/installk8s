package utils

const (
HttpRepotpl string = `
[remote]
name=RHEL Apache
baseurl=http://{{ . }}:18080
enabled=1
gpgcheck=0
`
Allhosts string = `
{{ range $key,$value :=  . }} {{ $key }}{{ $value }}{{ "\n" }}{{ end }}
`
Sysctlkubernetes string = `
net.bridge.bridge-nf-call-iptables=1
net.bridge.bridge-nf-call-ip6tables=1
net.ipv4.ip_forward=1
net.ipv4.tcp_tw_recycle=0
net.ipv4.neigh.default.gc_thresh1=1024
net.ipv4.neigh.default.gc_thresh2=2048
net.ipv4.neigh.default.gc_thresh3=4096
vm.swappiness=0
vm.overcommit_memory=1
vm.panic_on_oom=0
fs.inotify.max_user_instances=8192
fs.inotify.max_user_watches=1048576
fs.file-max=52706963
fs.nr_open=52706963
net.ipv6.conf.all.disable_ipv6=1
net.netfilter.nf_conntrack_max=2310720
`
Chronyconf string = `
server {{.FirstController.Ip }} iburst
driftfile /var/lib/chrony/drift
makestep 1.0 3
rtcsync
{{ range .AllHostsIp  -}} 
allow  {{ . }} 
{{ end }}
local stratum 10
logdir /var/log/chrony
`
CertSecureMode = 0640
CertMode = 0644

Etcdhost = `
'{% range .Ip %}
"{{ . }}",
{% endfor %} "127.0.0.1"'
`
)

type CfgVars struct {
	CertRootDir string
}

