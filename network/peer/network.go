package peer

import (
	"time"

	pb "github.com/pp2p/proto/paranoid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Network contains both the server and the client
type Network struct {
	Self   *Self
	peers  map[string]*Node
	creds  credentials.TransportCredentials
	secure bool
}

// NewNetwork creates a nre peer-to-peer network
func NewNetwork() *Network {
	return &Network{
		peers: make(map[string]*Node),
	}
}

// Peers of the network
func (n Network) Peers() []*Node {
	var p []*Node
	for _, peer := range n.peers {
		p = append(p, peer)
	}
	return p
}

// AddPeer to the network
func (n *Network) AddPeer(address string) (err error) {
	if _, ok := n.peers[address]; ok {
		return nil
	}
	// Attempt to establish connection
	opts := []grpc.DialOption{
		grpc.WithTimeout(5 * time.Second),
	}
	if n.secure {
		opts = append(opts, grpc.WithTransportCredentials(n.creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	nn := &Node{
		address: address,
	}
	nn.conn, err = grpc.Dial(address, opts...)
	if err != nil {
		return err
	}
	nn.client = pb.NewParanoidNetworkClient(nn.conn)

	n.peers[address] = nn
	return nil
}

// Node containing the information about itself and other peers
type Node struct {
	conn    *grpc.ClientConn
	client  pb.ParanoidNetworkClient
	address string // for display purposes
}

// Client connection to individual peer
func (n Node) Client() pb.ParanoidNetworkClient {
	return n.client
}

// Self is an abstraction over the Node, used to send requests out to other
// peers.
type Self struct {
	Node
}
