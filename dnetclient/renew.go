package dnetclient

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func Renew() {
	conn, err := grpc.Dial(DiscoveryAddr, grpc.WithInsecure())
	if err != nil {
		log.Println("[D] [E] failed to dial discovery server at ", DiscoveryAddr)
		return
	}
	defer conn.Close()

	dclient := pb.NewDiscoveryNetworkClient(conn)

	_, err = dclient.Renew(context.Background(), &pb.JoinRequest{Pool: "_", Node: &thisNode})
	if err != nil {
		log.Println("[D] [E] could not join")
		return
	}
}