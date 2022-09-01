package steps

import "bjsh/installk8s/cluster"

//"fmt"
//"os"
//"strings"
//"text/template"
//log "github.com/sirupsen/logrus"

type Step interface {
	Run() error
}


type Manager struct {
	Config *cluster.Cluster
	Steps []Step
}

func (m *Manager) Addstep(s ...Step) {
m.Steps=append(m.Steps,s...)
}

func (m *Manager) Run() error {
	for _,step := range m.Steps {
		result := step.Run()
		if result != nil {
			return result
		}
	}
	return nil
}



//type StepsConfig *Cluster

/*
func  CreateLocalRepo(s *Host) error {
	queryhttpd := "rpm -qa|grep httpd"
	curpath, err := os.Getwd()
	if err != nil {
		log.Fatal("can't get current dir")
		return err
	}
	masterbasecmd :="bash "+curpath+"/basefirst.sh"
	result,_ := stringtocmd(s,queryhttpd)
	fmt.Println(result)
	if strings.Contains(result,"httpd") {
		return nil
	}
	if _,err := stringtocmd(s,masterbasecmd);err != nil {
		log.Fatal(masterbasecmd+" failed")
		return err
	}
    return nil
}


func ConfigLocalRepo(s *Host) error {
// _,err := os.Stat("/etc/yum.repos.d/httplocal.repo") 不能这么写，这查看的是本地的
checkhttplocalcmd := "file /etc/yum.repos.d/httplocal.repo"
_,err := stringtocmd(s,checkhttplocalcmd)
if err == nil {
	fmt.Println("skip configlocalrepo")
	return nil
}
deleterepo := "cd /etc/yum.repos.d/;rm -f *.repo"
_,err = stringtocmd(s,deleterepo)
if err != nil {
	log.Fatal(deleterepo+" failed")
	return err
}
masterip0 := MasterIPs[0]
t,err := template.New("httprepo").Parse(HttpRepotpl)
if err != nil {
log.Fatal("parse masterip0 failed")
}
remotefile := "/etc/yum.repos.d/httplocal.repo"
dstfile,err := getremotefile(s,remotefile)
if err != nil {
	return err
}
err = t.Execute(dstfile,masterip0)
if err != nil {
	fmt.Println(err)
	log.Fatal("write to remote file failed")
}
defer dstfile.Close()
return nil
}


/*
t,_ := template.New("allhoststmp").Parse(Allhosts)
//if err != nil {
//	log.Fatal(err)
//}
for _,server := range a {
sftpfile,err := getremotefile(server,"/etc/hosts")
if err != nil {
	return err
}
//err = t.Execute(sftpfile,iph)
err = t.Execute(os.Stdout,s)
defer sftpfile.Close()
if err != nil {
return err
}
}
return nil
}




func Swapoff(s *Host) error {
	swapoffcmd := "swapoff -a && sed -i '/ swap / s/^/#/' /etc/fstab"
	_,err := stringtocmd(s,swapoffcmd)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}



func Setk8ssysctl(s *Host) error {
	t,err := template.New("sysctlk8s").Parse(Sysctlkubernetes)
	if err != nil {
	log.Fatal("parse sysctlk8s failed")
	}
	remotefile := "/etc/sysctl.d/kubernetes.conf"
	dstfile,err := getremotefile(s,remotefile)
	if err != nil {
		return err
	}
	err = t.Execute(dstfile,Sysctlkubernetes)
	if err != nil {
		fmt.Println(err)
		log.Fatal("write to remote file failed")
	}
	defer dstfile.Close()
	return nil

/*
writek8ssystl := "echo "+Sysctlkubernetes+"> /etc/sysctl.d/kubernetes.conf"	
_,err := stringtocmd(s,writek8ssystl)
if err != nil {
	log.Fatal(err)
	return err
}
return nil

}

*/