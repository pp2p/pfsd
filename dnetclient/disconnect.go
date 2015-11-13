package dnetclient

import (
	"errors"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

// Disconnect function used to disconnect from the server
func Disconnect() error {
	conn, err := grpc.Dial(DiscoveryAddr, grpc.WithInsecure())
	if err != nil {
		log.Println("[D] [E] failed to dial discovery server at ", DiscoveryAddr)
		return errors.New("Failed to dial discovery server")
	}
	defer conn.Close()

	dclient := pb.NewDiscoveryNetworkClient(conn)

	_, err = dclient.Disconnect(context.Background(), &pb.DisconnectRequest{Node: &pb.Node{ThisNode.IP, ThisNode.Port}})
	if err != nil {
		log.Println("[D] [E] could not send disconnect message")
		return errors.New("Could not send disconnect message")
	}

	return nil
}
