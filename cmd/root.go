package cmd

import (
	"github.com/urfave/cli/v2"
)

// App is the main urfave/cli.App for k0sctl
var App = &cli.App{
	Name:  "k8sinstall",
	Usage: "k8s install tool",
	Commands: []*cli.Command{
		applyCommand, //cmd/apply.go，apply命令
	},
}   //使用上面提到的cli服务，生成命令。