package network

import (
	"context"
	"fmt"
	"net"

	"github.com/grimthereaper/lockout/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (api *API) startGRPC() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", api.host, api.port))
	if err != nil {
		return err
	}
	api.grpcServer = grpc.NewServer()
	pb.RegisterWhitelistCheckerServer(api.grpcServer, &grpcServer{})

	reflection.Register(api.grpcServer)

	return api.grpcServer.Serve(lis)
}

type grpcServer struct{}

func (server *grpcServer) CheckIP(ctx context.Context, r *pb.IPCheckRequest) (*pb.IPCheckResponse, error) {
	whitelisted, err := checkIP(r.GetIp(), r.GetCountries())

	return &pb.IPCheckResponse{Whitelisted: whitelisted}, err
}
