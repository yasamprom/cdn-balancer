package balancerservice

import (
	"context"

	pb "github.com/yasamprom/cdn-balancer/internal/pb/api"
)

func (i *Implementation) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{}, nil
}
