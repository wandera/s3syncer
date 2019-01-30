package main

import (
	"github.com/WanderaOrg/s3syncer/cmd"
	log "github.com/sirupsen/logrus"
)

// main
func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
