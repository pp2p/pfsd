package main

import (
	"flag"

	log "github.com/pp2p/paranoid/logger"
	"github.com/pp2p/pfsd/network"
)

var (
	portFlag = flag.Int("port", 10101, "port to run pfsd on")
)

func main() {
	flag.Parse()

	log.SetLogDirectory("/var/log/paranoid/pfsd")

	n, err := network.New(*portFlag)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(n.Listen())
}
