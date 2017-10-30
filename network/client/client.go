// Package client contains the library that accepts connections for a local
// CLI client.
package client

import (
	"context"

	"github.com/pp2p/pfsd/network/peer"
	pb "github.com/pp2p/proto/client"
)

// Server implements the proto defintion of the server of the client
// communicator.
type Server struct {
	pn *peer.Network
}

// NewServer creates a new server for the client connection
func NewServer(pn *peer.Network) *Server {
	return &Server{
		pn: pn,
	}
}

// Init RPC
func (s *Server) Init(ctx context.Context, req *pb.InitRequest) (*pb.InitResponse, error) {
	return nil, nil
}

// Join RPC
func (s *Server) Join(ctx context.Context, req *pb.JoinRequest) (*pb.EmptyMessage, error) {
	return nil, nil
}

// Status RPC
func (s *Server) Status(ctx context.Context, req *pb.EmptyMessage) (*pb.StatusResponse, error) {
	return nil, nil
}
