package steps

import (
	"bjsh/installk8s/cluster"
	"text/template"
	"bjsh/installk8s/utils"
	//log "github.com/sirupsen/logrus"
	"strings"
	"fmt"
)


type PrepareHost struct {
	Config *cluster.Cluster
}

func (p PrepareHost) Run() error{
	p.Config.Hosts.Generatehostname()
    var tasks []func(h *cluster.Host) error
    sethostname := func(h *cluster.Host) error {
		getprevioushostname := "hostname"
		previoushostname,err := h.Execcmd(getprevioushostname)
		if err != nil {
			return err
		}
		
		if previoushostname != h.Hostname {
			sethostnamecmd := fmt.Sprintf("hostnamectl set-hostname %s",h.Hostname)
			if _,err := h.Execcmd(sethostnamecmd); err != nil {
				return err
			}
			
		}
		return nil
	}

	updateetchosts := func(h *cluster.Host) error {
		getprevioushostname := "hostname"
		previoushostname,err := h.Execcmd(getprevioushostname)
		if err != nil {
			return err
		}
		replacehosts := fmt.Sprintf("grep -q '127.0.1.1 %s %s' /etc/hosts || sed -i 's/%s/%s/g' /etc/hosts",
		h.Hostname,h.Hostname,strings.Replace(previoushostname,"\n","",-1),h.Hostname)
		 _,err = h.Execcmd(replacehosts)
        if err != nil {
	    return err
            }
		for _,s := range p.Config.Hosts {
			setcmd := fmt.Sprintf("grep -q '%s %s' /etc/hosts || echo '%s %s' >> /etc/hosts",s.Ip,s.Hostname,s.Ip,s.Hostname)
			//setcmd := "grep -q "+'s.Ip s.Hostname+" /etc/hosts || echo "+ s.Ip+" "+s.Hostname+" >> /etc/hosts"
			_,err := h.Execcmd(setcmd)
			if err != nil {
				return err
			}
		}
		return nil
	}

	disablefirewalld := func(h *cluster.Host) error {
		checkfirewalld := "systemctl status firewalld"
	    checkresult,err := h.Execcmd(checkfirewalld)
		if err != nil {
			return err
		}
		
	    if !strings.Contains(checkresult,"inactive") {
			stopfirewallcmd := "systemctl stop firewalld;systemctl disable firewalld;iptables -F;iptables -X"
			if _,err := h.Execcmd(stopfirewallcmd);err != nil {
				return err
			}

	    }
        return nil
	}
    
	swapoff := func(h *cluster.Host) error {

		swapcheck := "swapon -s"
		swapstatus,err := h.Execcmd(swapcheck)
		if err != nil {
			return err
		}
		if swapstatus != ""{
			swapoffcmd := "swapoff -a && sed -i '/ swap / s/^/#/' /etc/fstab"
			if _,err := h.Execcmd(swapoffcmd);err != nil {
				return err
			}
		}
		return nil
	}
    
	 disableselinux := func(h *cluster.Host) error {
		getselinuxstatus := "getenforce"
		selinuxstatus,err := h.Execcmd(getselinuxstatus)
		if err != nil {
			return err
		}
		if selinuxstatus == "Enforcing\n" {
        disableselinux := "setenforce 0;sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config"
		if _,err := h.Execcmd(disableselinux);err != nil {
			return err
		}
		}
        return nil
	 }

	 modifysysctl := func(h *cluster.Host) error {
		if !utils.FileExists("/etc/sysctl.d/kubernetes.conf") {
			t,err := template.New("sysctlk8s").Parse(utils.Sysctlkubernetes)
			if err != nil {
			return err
			}
			remotefile := "/etc/sysctl.d/kubernetes.conf"
			dstsysctlfile,err := h.Getremotefile(remotefile)
			if err != nil {
			return err
			}
			defer dstsysctlfile.Close()
			err = t.Execute(dstsysctlfile,utils.Sysctlkubernetes)
			if err != nil {
			return err
			}
			}
        return nil
	 }

     settimezone := func(h *cluster.Host) error {
     
		gettimezone := "timedatectl |grep 'Time zone'|awk -F':' '{print $2}'|awk '{print $1}'"
		timezone,err := h.Execcmd(gettimezone)
	   if err != nil {
		   return err
	   }
	   if timezone != "Asia/Shanghai" {
	   settimezonecmd := "timedatectl set-timezone Asia/Shanghai"
	   if _,err := h.Execcmd(settimezonecmd);err != nil {
		return err
	   }
	   }
	   installchronycmd := "rpm -qa|grep chrony || yum install -y chrony"
       if _,err := h.Execcmd(installchronycmd);err != nil {
		return err
	   }

	   type chronyconfig struct {
		   AllHostsIp []string
		   FirstController cluster.Host
	   }
	   var temp chronyconfig
	   temp.AllHostsIp = p.Config.Hosts.GetAllHostsIp()
	   temp.FirstController = *p.Config.Hosts.FirstController()
	   t,err := template.New("chronyconf").Parse(utils.Chronyconf)
	   if err != nil {
		   return err
		   }
	   remotefile := "/etc/chrony.conf"
	   dsttimefile,err := h.Getremotefile(remotefile)
	   if err != nil {
		   return err
	   }
	   defer dsttimefile.Close()
	   err = t.Execute(dsttimefile,temp)
	   if err != nil {
		   return err
	   }
	   restartchronydcmd := "systemctl restart chronyd"
	   if _,err := h.Execcmd(restartchronydcmd);err != nil {
		   return err
	   }
       return nil


	 }
	
    tasks=append(tasks,sethostname,updateetchosts,disablefirewalld,swapoff,disableselinux,modifysysctl,settimezone)
	
	err  := make(chan error) 
	var errors []string
     for _,task := range tasks {
		go func() {
			err <- p.Config.Hosts.ParallelEach(task)
		}()
		go func() {
			for e := range err {
				if e != nil {
					errors = append(errors,   e.Error())
				}
			}

		}()
		if len(errors)>0 {
			return fmt.Errorf("%s", strings.Join(errors, "\n - "))
		}	
		  
	 }
     return nil
	}
    /*
    return p.Config.Hosts.ParallelEach(func(h *cluster.Host) error {
		cmdfirst := cluster.CmdStrings{Cmd: []string{} }
		getprevioushostname := "hostname"
		previoushostname,err := h.Execcmd(getprevioushostname)
		if err != nil {
			return err
		}
        if previoushostname != h.Hostname {
			sethostnamecmd := fmt.Sprintf("hostnamectl set-hostname %s",h.Hostname)
			cmdfirst.Addcmd(sethostnamecmd)
		}
    
		
		replacehosts := fmt.Sprintf("grep -q '127.0.1.1 %s %s' /etc/hosts || sed -i 's/%s/%s/g' /etc/hosts",
		h.Hostname,h.Hostname,strings.Replace(previoushostname,"\n","",-1),h.Hostname)
		
		//replacehosts := "sed -i 's/"+strings.Replace(previoushostname,"\n","",-1)+"/"+h.Hostname+"/g'"+" /etc/hosts"
		for _,s := range p.Config.Hosts {
			setcmd := fmt.Sprintf("grep -q '%s %s' /etc/hosts || echo '%s %s' >> /etc/hosts",s.Ip,s.Hostname,s.Ip,s.Hostname)
			//setcmd := "grep -q "+'s.Ip s.Hostname+" /etc/hosts || echo "+ s.Ip+" "+s.Hostname+" >> /etc/hosts"
			_,err := h.Execcmd(setcmd)
			if err != nil {
				return err
			}
		}
        checkfirewalld := "systemctl status firewalld"
	    checkresult,_ := h.Execcmd(checkfirewalld)
	    if !strings.Contains(checkresult,"inactive") {
			stopfirewallcmd := "systemctl stop firewalld;systemctl disable firewalld;iptables -F;iptables -X"
			cmdfirst.Addcmd(stopfirewallcmd)
	    }
        swapcheck := "swapon -s"
		swapstatus,err := h.Execcmd(swapcheck)
		if err != nil {
			return err
		}
		if swapstatus != ""{
			swapoffcmd := "swapoff -a && sed -i '/ swap / s/^/#/' /etc/fstab"
			cmdfirst.Addcmd(swapoffcmd)
		}
		
		getselinuxstatus := "getenforce"
		selinuxstatus,err := h.Execcmd(getselinuxstatus)
		if err != nil {
			return err
		}
		if selinuxstatus == "Enforcing\n" {
        disableselinux := "setenforce 0;sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config"
		cmdfirst.Addcmd(disableselinux)
		}
		if !utils.FileExists("/etc/sysctl.d/kubernetes.conf") {
        t,err := template.New("sysctlk8s").Parse(utils.Sysctlkubernetes)
	    if err != nil {
	    return err
	    }
	    remotefile := "/etc/sysctl.d/kubernetes.conf"
	    dstsysctlfile,err := h.Getremotefile(remotefile)
	    if err != nil {
		return err
	    }
		defer dstsysctlfile.Close()
	    err = t.Execute(dstsysctlfile,utils.Sysctlkubernetes)
	    if err != nil {
		return err
	    }
	    }
		gettimezone := "timedatectl |grep 'Time zone'|awk -F':' '{print $2}'|awk '{print $1}'"
         timezone,err := h.Execcmd(gettimezone)
        if err != nil {
			return err
		}
		if timezone != "Asia/Shanghai" {
	    settimezonecmd := "timedatectl set-timezone Asia/Shanghai"
        cmdfirst.Addcmd(settimezonecmd)
		}
		installchronycmd := "rpm -qa|grep chrony || yum install -y chrony"
	    cmdfirst.Addcmd(replacehosts,installchronycmd)
		err = cmdfirst.Runcmd(h)
		if err != nil {
			return err
		}
		type chronyconfig struct {
			AllHostsIp []string
			FirstController cluster.Host
		}
		var temp chronyconfig
		temp.AllHostsIp = p.Config.Hosts.GetAllHostsIp()
		temp.FirstController = *p.Config.Hosts.FirstController()
		t,err := template.New("chronyconf").Parse(utils.Chronyconf)
		if err != nil {
			return err
			}
		remotefile := "/etc/chrony.conf"
		dsttimefile,err := h.Getremotefile(remotefile)
		if err != nil {
			return err
		}
		defer dsttimefile.Close()
		err = t.Execute(dsttimefile,temp)
		if err != nil {
			return err
		}
		restartchronydcmd := "systemctl restart chronyd"
		if _,err := h.Execcmd(restartchronydcmd);err != nil {
			return err
		}
	
		return nil
	})
	

}
*/







/*
type DisableFirewall struct {
	Config *cluster.Cluster
}


func (c DisableFirewall) Run() error {
	return c.Config.Hosts.ParallelEach(func(h *cluster.Host) error {
		checkfirewalld := "systemctl status firewalld"
	checkresult,_ := utils.Stringtocmd(h,checkfirewalld)
	if strings.Contains(checkresult,"inactive") {
		return nil
	}
	stopfirewallcmd := "systemctl stop firewalld;systemctl disable firewalld;iptables -F;iptables -X"
	_,err := utils.Stringtocmd(h,stopfirewallcmd)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
	})
}

type SetTimeZone struct {
	Config *cluster.Cluster
}

func (c SetTimeZone) Run() error {
	return c.Config.Hosts.ParallelEach(func(h *cluster.Host) error {
		settimezonecmd := "timedatectl set-timezone Asia/Shanghai"
		if _,err := utils.Stringtocmd(h,settimezonecmd);err != nil {
			return err
		}
		installchronycmd := "yum install -y chrony"
		if _,err := utils.Stringtocmd(h,installchronycmd);err != nil {
			return err
		}
		type chronyconfig struct {
			AllHostsIp []string
			FirstController cluster.Host
		}
		var temp chronyconfig
		temp.AllHostsIp = c.Config.Hosts.GetAllHostsIp()
		temp.FirstController = *c.Config.Hosts.FirstController()
		t,err := template.New("chronyconf").Parse(utils.Chronyconf)
		if err != nil {
			log.Fatal("parse chronyconf failed")
			}
		remotefile := "/etc/chrony.conf"
		dstfile,err := utils.Getremotefile(h,remotefile)
		if err != nil {
			return err
		}
		defer dstfile.Close()
		err = t.Execute(dstfile,temp)
		if err != nil {
			fmt.Println(err)
			return err
		}
		restartchronydcmd := "systemctl restart chronyd"
		if _,err := utils.Stringtocmd(h,restartchronydcmd);err != nil {
			return err
		}
		return nil
		})
}

type UpdateHostName struct {
	Config *cluster.Cluster
}


func (c UpdateHostName) Run() error {
	c.Config.Hosts.Generatehostname()
	return c.Config.Hosts.ParallelEach(func(h *cluster.Host) error {
		getprevioushostname := "hostname"
		previoushostname,err := utils.Stringtocmd(h,getprevioushostname)
		if err != nil {
			return err
		}
		sethostnamecmd := "hostnamectl set-hostname "+ h.Hostname
		_,err1 := utils.Stringtocmd(h,sethostnamecmd)
		if err1 != nil {
			log.Fatal("set hostname failed")
			return err1
		}
		replacehosts := "sed -i 's/"+strings.Replace(previoushostname,"\n","",-1)+"/"+h.Hostname+"/g'"+" /etc/hosts"
		_,err2 := utils.Stringtocmd(h,replacehosts)
		if err2 != nil {
			return err2
		}
		return nil
	})
	}

type DisableSeLinux struct {
	Config *cluster.Cluster
	}

func (c DisableSeLinux) Run() error {
		return c.Config.Hosts.ParallelEach(func(h *cluster.Host) error {
			disableselinux := "setenforce 0;sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config"
			_,err := utils.Stringtocmd(h,disableselinux)
			if err != nil {
				log.Fatal(err)
				return err
			}
			return nil
	})
}

type SetEtcHosts struct {
	Config *cluster.Cluster
	}



func (c SetEtcHosts) Run() error {
	return c.Config.Hosts.ParallelEach(func(h *cluster.Host) error {
		for _,s := range c.Config.Hosts {
			
				ip := s.Ip
				checkcmd := "echo "+ip+" "+s.Hostname+" /etc/hosts"
				checkresult,_ := utils.Stringtocmd(h,checkcmd)
				fmt.Println(checkresult)
				if checkresult == "" {
				   return nil
				}
				cmd := "echo "+ip+" "+s.Hostname+" >> /etc/hosts"
				_,err := utils.Stringtocmd(h,cmd)
				if err != nil {
					log.Fatal(err)
					return err
				}
			}
		
		return nil
})
}

type SwapOff struct {
	Config *cluster.Cluster
	}


func (c SwapOff) Run() error {
	return c.Config.Hosts.ParallelEach(func(h *cluster.Host) error {
		swapoffcmd := "swapoff -a && sed -i '/ swap / s/^/#/' /etc/fstab"
		_,err := utils.Stringtocmd(h,swapoffcmd)
		if err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	})
}

type ModifySysctl struct {
	Config *cluster.Cluster
	}


func (c ModifySysctl) Run() error {
	return c.Config.Hosts.ParallelEach(func(h *cluster.Host) error {
		t,err := template.New("sysctlk8s").Parse(utils.Sysctlkubernetes)
	if err != nil {
	log.Fatal("parse sysctlk8s failed")
	}
	remotefile := "/etc/sysctl.d/kubernetes.conf"
	dstfile,err := utils.Getremotefile(h,remotefile)
	if err != nil {
		return err
	}
	err = t.Execute(dstfile,utils.Sysctlkubernetes)
	if err != nil {
		fmt.Println(err)
		log.Fatal("write to remote file failed")
	}
	defer dstfile.Close()
	return nil
	})
}

*/