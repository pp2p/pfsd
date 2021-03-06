package pnetserver

import (
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/pp2p/paranoid/proto/paranoidnetwork"
	"github.com/pp2p/pfsd/globals"
)

// NewGeneration receives requests from nodes asking to create a new KeyPiece
// generation in preparation for joining the cluster.
func (s *ParanoidServer) NewGeneration(ctx context.Context, req *pb.NewGenerationRequest) (*pb.NewGenerationResponse, error) {
	if req.PoolPassword == "" {
		if len(globals.PoolPasswordHash) != 0 {
			return &pb.NewGenerationResponse{}, grpc.Errorf(codes.InvalidArgument,
				"cluster is password protected but no password was given")
		}
	} else {
		err := bcrypt.CompareHashAndPassword(globals.PoolPasswordHash,
			append(globals.PoolPasswordSalt, []byte(req.PoolPassword)...))
		if err != nil {
			return &pb.NewGenerationResponse{}, grpc.Errorf(codes.InvalidArgument,
				"unable to request new generation: password error:", err)
		}
	}

	Log.Info("Requesting new generation")
	generationNumber, peers, err := globals.RaftNetworkServer.RequestNewGeneration(req.GetRequestingNode().Uuid)
	if err != nil {
		return &pb.NewGenerationResponse{}, grpc.Errorf(codes.Unknown, "unable to create new generation")
	}
	return &pb.NewGenerationResponse{
		GenerationNumber: int64(generationNumber),
		Peers:            peers,
	}, nil
}
