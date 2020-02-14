package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/wandera/s3syncer/cmd"
)

// main
func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
