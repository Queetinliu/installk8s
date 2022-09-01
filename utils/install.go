package utils
import (
	//"log"
	//"fmt"
)
/*
func Install() {
	cfgbool,_ := checkcfgfileexists(string(Cfgfile))
	indirbool,_  := checkpackagesexists(string(Packagesname))
	if cfgbool && indirbool {
		k8scfg := &K8sConfigType{}
		if err := k8scfg.Load(string(Cfgfile));err != nil {
			log.Fatal("load config failed")
		}

		//allhosts,err := getallhost(*k8scfg)
		//if err != nil {
		//	log.Fatal("get all host failed")
		//}
		Master0 := k8scfg.Masters[0]
		hostnames,_ := generatehostname(*k8scfg)
		err := CreateLocalRepo(Master0)
        if err != nil {
			log.Fatal(err)
		}
		for server := range hostnames {
			 ConfigLocalRepo(server)
             UpgradeKernel(server)
             Disablefirewall(server)
			 Swapoff(server)
			 Disableselinux(server)
			 Setk8ssysctl(server)
			 Settimezone(server,hostnames)
		}
        
		for server,hostname := range hostnames {
        Sethostname(server,hostname)
		}
		//fmt.Println(allhost)
		Setetchosts(hostnames)
		//}    
}
}
*/