package main

import (
	"os"
	"bjsh/installk8s/cmd"
	log "github.com/sirupsen/logrus"
)



func main() {
	err := cmd.App.Run(os.Args) //这里去看cmd/root.go
	if err != nil {
		log.Fatal(err)
	}
	
}
