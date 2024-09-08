package balancerservice

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	model "github.com/yasamprom/cdn-balancer/internal/model"
	pb "github.com/yasamprom/cdn-balancer/internal/pb/api"
)

func (i *Implementation) GetLink(ctx context.Context, req *pb.GetLinkRequest) (*pb.GetLinkResponse, error) {

	if !isValidLink(req.Uri) {
		return nil, status.Error(codes.InvalidArgument, "invalid uri")
	}

	res, err := i.uc.GetLink(ctx, req.Uri)
	if err != nil {
		if err == model.ErrBadRoutePercent {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &pb.GetLinkResponse{
		Uri: res,
	}, nil
}

func isValidLink(link string) bool {
	return model.RegexpLink.Match([]byte(link))
}
