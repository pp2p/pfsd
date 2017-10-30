package network

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	log "github.com/pp2p/paranoid/logger"
	"github.com/pp2p/pfsd/network/client"
	"github.com/pp2p/pfsd/network/peer"

	cpb "github.com/pp2p/proto/client"
	ppb "github.com/pp2p/proto/paranoid"
)

// Network holds all connections of the network
type Network struct {
	lis net.Listener
	srv *grpc.Server
	pn  *peer.Network
}

// New creates a new Network instance
func New(port int) (*Network, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()

	pn := peer.NewNetwork()
	cc := client.NewServer(pn)

	ppb.RegisterParanoidNetworkServer(s, pn)
	cpb.RegisterClientServiceServer(s, cc)

	reflection.Register(s)

	return &Network{
		lis: lis,
		srv: s,
	}, nil
}

// Listen on the network
func (n *Network) Listen() error {
	log.Infof("Starting pfsd on %s", n.lis.Addr().String())
	return n.srv.Serve(n.lis)
}
