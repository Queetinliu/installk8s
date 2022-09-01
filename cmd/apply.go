package cmd

import (
	"github.com/urfave/cli/v2"
	"bjsh/installk8s/cluster"
	"bjsh/installk8s/steps"
)

var applyCommand = &cli.Command{    //这里使用github.com/urfave/cli/v2这个工具
	Name:  "apply",
	Usage: "Apply a k8sinstall configuration",
	Flags: []cli.Flag{
		configFlag,
	},
	Before: actions(initLogging, initConfig), 
	//都定义在cmd/flags.go。首先执行这些操作
	Action: func(ctx *cli.Context) error {
	
		config := steps.Manager{ Config: ctx.Context.Value(ctxConfigKey{}).(*cluster.Cluster) }
		//config := ctx.Context.Value(ctxConfigKey{}).(*install.Cluster)
		//stepsconfig.Hosts.ParallelEach(install.Settimezone)
        //fmt.Println(config.Hosts)
		//config.Addstep(&steps.UpdateHostName{Config: config.Config},&steps.SetEtcHosts{Config: config.Config},&steps.DisableFirewall{Config: config.Config},&steps.SwapOff{Config: config.Config},
			//&steps.DisableSeLinux{Config: config.Config},&steps.ModifySysctl{Config: config.Config},&steps.SetTimeZone{Config: config.Config})
			config.Addstep(&steps.PrepareHost{Config: config.Config},&steps.GenerateCerts{Config: config.Config})
		/*
		if err := stepsconfig.Settimezone();err != nil {
			return err
		}
		if err := stepsconfig.UpdateHostname();err != nil {
			return err
		}
		if err := stepsconfig.Disablefirewall();err != nil {
			return err
		}
		if err := stepsconfig.DisableSELinux();err != nil {
			return err
		}
		if err := stepsconfig.Setetchosts();err != nil {
			return err
		}
		if err := stepsconfig.Swapoff();err != nil {
			return err
		}
		if err := stepsconfig.Modifysysctl();err != nil {
			return err
		}
		*/
		if err := config.Run();err != nil {
			return err
		}

		return nil
},
}

