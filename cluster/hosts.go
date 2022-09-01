package cluster
import (
"sync"
"fmt"
"strings"
"strconv"
)

type Hosts []*Host

func (hosts Hosts) ParallelEach(filter ...func(h *Host) error) error {  //如注释所说，是并行地运行多个方法
	var wg sync.WaitGroup
	var errors []string
	type erritem struct {
		address string
		err     error
	}
	ec := make(chan erritem, 1)

	for _, f := range filter {
		wg.Add(len(hosts))  //将主机数量作为任务添加进去

		for _, h := range hosts {
			go func(h *Host) {
				ec <- erritem{h.Ip, f(h)}  //每个主机执行参数中方法，将错误拼接起来
			}(h)
		}

		go func() {
			for e := range ec {
				if e.err != nil {
					errors = append(errors, fmt.Sprintf("%s: %s", e.address, e.err.Error()))
				}
				wg.Done()
			}
		}()

		wg.Wait()
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed on %d hosts:\n - %s", len(errors), strings.Join(errors, "\n - "))
	}

	return nil
}

func (hosts Hosts) WithRole(s string) Hosts {
	return hosts.Filter(func(h *Host) bool {
		return h.Role == s
	})
}

// Controllers returns hosts with the role "controller"
func (hosts Hosts) Controllers() Hosts {
	return hosts.Filter(func(h *Host) bool { return h.IsController() })
}

// Workers returns hosts with the role "worker"
func (hosts Hosts) Workers() Hosts {
	return hosts.Filter(func(h *Host) bool { return h.IsWorker() })
}


func (hosts Hosts) Etcds() Hosts {
	return hosts.Filter(func(h *Host) bool { return h.IsEtcd() })
}

func (hosts Hosts) Filter(filter func(h *Host) bool) Hosts {
	result := make(Hosts, 0, len(hosts))

	for _, h := range hosts {
		if filter(h) {
			result = append(result, h)
		}
	}

	return result
}

func (hosts Hosts) FirstController() *Host {
	h := Hosts{}

	if len(hosts) == 0 {
		return nil
	}
	for _,host := range hosts {
		if host.IsController() {
			h = append(h, host)
		}
	}
	return (h)[0]
}

func (hosts Hosts) GetAllHostsIp() []string {
   IPS := []string{}
   for _,h := range hosts {
	IPS=append(IPS,h.Ip)
   }
   return IPS
}

//generate auto name if doesn't specify
func (hosts *Hosts) Generatehostname() error {
    var (
		Controllers Hosts
        Workers Hosts
        Etcds Hosts
	)
	Controllers = hosts.Controllers()
		
	for index,s := range Controllers {
	if s.Hostname == "" {
	   if index < 9 {
		s.Hostname= "k8smaster" + "0"+strconv.Itoa(index+1)
	   } else {
		s.Hostname = "k8smaster" + strconv.Itoa(index+1)
	   }
	   
	}
	}
     
    Workers = hosts.Workers()

	for index,s := range Workers {
	if s.Hostname == "" {
		if index < 9 {
			s.Hostname = "k8snode" + "0"+strconv.Itoa(index+1)
			
		   } else {
	s.Hostname = "k8snode" + strconv.Itoa(index+1)
		   }
			 }
		}
	
	Etcds = hosts.Etcds()
    for index,s := range Etcds {
	if s.Hostname == "" {
		if index < 9 {
			s.Hostname = "k8setcd" + "0"+strconv.Itoa(index+1)
		   } else {
	s.Hostname = "k8setcd" + strconv.Itoa(index+1)
				
		   }
			 }
		}
	
	return nil
	}