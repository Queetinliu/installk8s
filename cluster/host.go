package cluster

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"fmt"
	"net"
	log "github.com/sirupsen/logrus"
)

type Host struct {
Role string `yaml:"role"`
Ip string `yaml:"ip"`
Port string `yaml:"port"`
Password string `yaml:"password"`
Hostname string `yaml:"hostname"`
}


type CmdStrings struct {
	Cmd []string
}

func (h *Host) IsController() bool {
	return h.Role == "controller" || h.Role == "controller+worker" || h.Role == "single"  || h.Role == "controller+etcd" || h.Role == "all"
}

func (h *Host) IsWorker() bool {
	return h.Role == "worker" || h.Role == "controller+worker" || h.Role == "single" || h.Role == "all"
}

func (h *Host) IsEtcd() bool {
	return h.Role == "single" || h.Role == "controller+etcd" || h.Role == "all"
}


func(c *CmdStrings) Addcmd(cmd ...string)  {
		c.Cmd=append(c.Cmd,cmd...)
}
func(c CmdStrings) Runcmd(h *Host) error {
	for _,cmd := range c.Cmd {
		_,err := h.Execcmd(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}


func(h *Host) Execcmd(cmd string) (string,error) {
	sshclient,err := h.getsshclient()
	if err != nil {
		log.Fatal(err)
        return "",err
	}
	session, err := sshclient.NewSession()
	if err != nil {
		log.Fatal("create ssh session failed")
	    return "",err
	}
	defer session.Close()
	//var buffstdout bytes.Buffer
	//session.Stdout = &buffstdout
	//if err := session.Run(cmd); err != nil {
		//用CombinedOutput会同时获取Stdout和Stderr，目前开发阶段比较有用
	fmt.Println(h,cmd)	
	stdoutandstderr,err := session.CombinedOutput(cmd)
	if err != nil {
		//var buffstderr bytes.Buffer
		//session.Stderr = &buffstderr
		//fmt.Println(buffstdout.String())
		//输出是[]bytes，必须这样转换
	//fmt.Println(string(stdoutandstderr))
	return string(stdoutandstderr),err
	//log.Fatal("run "+cmd+" failed")
	//return err
	}
	//fmt.Println(buffstdout.String())
	return string(stdoutandstderr),nil
	}
	

func(h *Host) getsshclient() (*ssh.Client, error) { 
	clientconfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{ssh.Password(h.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
		},
}	   
	// connet to ssh
	addr := fmt.Sprintf("%s:%s", h.Ip, h.Port)
	 sshClient, err := ssh.Dial("tcp", addr, clientconfig)
	 if err != nil {
	   return nil, err
	}
     return sshClient,nil
    }


func(h *Host) getsftpclient() (*sftp.Client, error) {
	sshclient,err := h.getsshclient()
	if err != nil {
		return nil,err
	}
	sftpClient, err := sftp.NewClient(sshclient)
	if err != nil {
	return nil, err
	}
		return sftpClient, nil
   }


func(h *Host) Getremotefile(file string) (*sftp.File,error) {
    sftpclient,err := h.getsftpclient()
	if err != nil {
		return nil,err
	}
	dstFile, err := sftpclient.Create(file)
	if err != nil {
	log.Fatal(err)
	return nil,err
	}
	//defer dstFile.Close()
	return dstFile,nil
}

