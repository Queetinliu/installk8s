package utils

import (
	"bjsh/installk8s/cluster"
	"os"
	"path/filepath"
)


func GetCertRootDir() string {
	pwd,err := os.Getwd()
	if err != nil {
		return ""
	}
	CertPath := filepath.Join(pwd,"/cert")
	if !DirExist(CertPath) {
		err := os.Mkdir(CertPath,0644)
		if err != nil {
			return ""
		}
	}
    return CertPath
}


   /*
func parseip(s interface{}) ([]Host,error) {
   var servers []Host
	for _,ips := range strings.Split(s.(string),",") {
    ipconfig := strings.Split(ips, ":")
	switch (len(ipconfig)) {
	case 2:
		//ip := ipconfig[0]
		//password := ipconfig[1]
		//servers = append(servers,ServerInfo{ip,"22",password})
	
    case 3:
		//ip := ipconfig[0]
		//port := ipconfig[1]
		//password := ipconfig[2]
		//servers = append(servers,ServerInfo{ip,port,password})

	default:
		return nil,errors.New("invalid server config")
	}
}
return servers,nil
   }

/*
func getallhost(k K8sConfigType) ([]ServerInfo,error) {
allhosts := make(map[ServerInfo]bool)
if k.Masters != nil {
	for _,masters := range k.Masters {
		if _,ok := allhosts[masters]; !ok {
			allhosts[masters] = true
		}
	}
}
if k.Nodes != nil {
	for _,nodes := range k.Nodes {
		if _,ok := allhosts[nodes];!ok {
			allhosts[nodes] = true
		}
	}
}
if k.Etcds != nil {
	for _,etcds := range k.Etcds {
     if _,ok := allhosts[etcds];!ok {
		 allhosts[etcds] = true
	 }
	}
}
var s []ServerInfo
for server := range allhosts {
	s = append(s,server)
}
return s,nil
}
*/

type  ServerHostName map[cluster.Host]string
/*

*/
