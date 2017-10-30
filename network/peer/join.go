package peer

import (
	"context"

	pb "github.com/pp2p/proto/paranoid"
)

// Join rpc
func (n *Network) Join(ctx context.Context, req *pb.JoinRequest) (*pb.EmptyMessage, error) {
	return nil, nil
}
