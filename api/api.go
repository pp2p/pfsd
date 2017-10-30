package api

import (
  "context"
  pb "github.com/pp2p/proto/client"
  "google.golang.org/grpc"
)

type InitRequest = pb.InitRequest
type InitResponse = pb.InitResponse

// Connection is handing all the aP
type Connection struct {
  client pb.ClientServiceClient
}

// New creates a new connection to the pfsd
func New(pfsdAddress string) (*Connection, error) {
  conn, err := grpc.Dial(pfsdAddress)
  if err != nil {
    return nil, err
  }
	return &Connection{
    client: pb.NewClientServiceClient(conn),
  }, nil
}

func (c *Connection) Init(ctx context.Context, req InitRequest) (*InitResponse, error) {
  return c.client.Init(ctx, &req)
}


// opts := []grpc.DialOption{
//   grpc.WithTimeout(5 * time.Second),
// }
// if n.secure {
//   opts = append(opts, grpc.WithTransportCredentials(n.creds))
// } else {
//   opts = append(opts, grpc.WithInsecure())
// }
//
// nn := &Node{
//   Node: p,
// }
// nn.conn, err = grpc.Dial(fmt.Sprintf("%s:%s", p.GetIp(), p.GetPort()), opts...)
// if err != nil {
//   return err
// }
// nn.client = pb.NewParanoidNetworkClient(nn.conn)
//
// n.peers[p.GetUuid()] = nn
