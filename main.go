package main

import (
	"flag"
	"fmt"
	"github.com/cpssd/paranoid/pfsd/dnetclient"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/icserver"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
	"github.com/cpssd/paranoid/pfsd/pnetserver"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	srv *grpc.Server

	certFile = flag.String("cert", "", "TLS certificate file - if empty connection will be unencrypted")
	keyFile  = flag.String("key", "", "TLS key file - if empty connection will be unencrypted")
)

func startIcAndListen(pfsDir string) {
	defer globals.Wait.Done()

	globals.Wait.Add(1)
	go icserver.RunServer(pfsDir, true)

	for {
		select {
		case message := <-icserver.MessageChan:
			pnetclient.SendRequest(message)
		case _, ok := <-globals.Quit:
			if !ok {
				return
			}
		}
	}
}

func startRPCServer(lis *net.Listener) {
	var opts []grpc.ServerOption
	if *certFile != "" && *keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalln("FATAL: Failed to generate TLS credentials:", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	srv = grpc.NewServer(opts...)
	pb.RegisterParanoidNetworkServer(srv, &pnetserver.ParanoidServer{})
	globals.Wait.Add(1)
	go srv.Serve(*lis)
}

func main() {
	flag.Parse()
	if len(os.Args) < 4 {
		fmt.Print("Usage:\n\tpfsd <paranoid_directory> <Discovery Server> <Discovery Port>\n")
		os.Exit(1)
	}
	discoveryPort, err := strconv.Atoi(os.Args[3])
	if err != nil || discoveryPort < 1 || discoveryPort > 65535 {
		log.Fatalln("FATAL: Discovery port must be a number between 1 and 65535, inclusive.")
	}
	pnetserver.ParanoidDir = os.Args[1]
	globals.Wait.Add(1)
	go startIcAndListen(pnetserver.ParanoidDir)
	globals.Server, err = pnetclient.GetIP()
	if err != nil {
		log.Fatalln("FATAL: Cant get internal IP.")
	}

	if _, err := os.Stat(pnetserver.ParanoidDir); os.IsNotExist(err) {
		log.Fatalln("FATAL: path", pnetserver.ParanoidDir, "does not exist.")
	}
	if _, err := os.Stat(path.Join(pnetserver.ParanoidDir, "meta")); os.IsNotExist(err) {
		log.Fatalln("FATAL: path", pnetserver.ParanoidDir, "is not valid PFS root.")
	}

	//Asking for port 0 requests a random free port from the OS.
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("FATAL: Failed to start listening : %v.\n", err)
	}
	splits := strings.Split(lis.Addr().String(), ":")
	port, err := strconv.Atoi(splits[len(splits)-1])
	if err != nil {
		log.Fatalln("Could not parse port", splits[len(splits)-1], " Error :", err)
	}
	dnetclient.SetDiscovery(os.Args[2], os.Args[3], strconv.Itoa(port))
	globals.Port = port
	dnetclient.JoinDiscovery("_")
	startRPCServer(&lis)
	HandleSignals()
}
