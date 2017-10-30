package peer

import (
	"context"

	log "github.com/pp2p/paranoid/logger"
	pb "github.com/pp2p/proto/paranoid"
)

// Ping implements the Ping RPC.
func (n *Network) Ping(ctx context.Context, _ *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	log.V(1).Infof("Got ping from peer: %+v", ctx)

	// n.AddPeer()
	return &pb.EmptyMessage{}, nil
}

// Ping other peers to update their status
func (s *Self) Ping(ctx context.Context, n *Network) {
	for _, peer := range n.Peers() {
		if _, err := peer.Client().Ping(ctx, &pb.EmptyMessage{}); err != nil {
			// TODO(voy): Remove the node from the list if there is an issue.
			log.Errorf("Can't ping %s: %v", peer.address, err)
		}
	}
}
