package routes

import (
	"context"

	hello "github.com/iMohamedSheta/xapp/grpc/v1/hello/proto"

	"google.golang.org/grpc"
)

func LoadGrpcRoutes(grpcServer *grpc.Server) {
	// Here you can register your grpc routes
	// e.g.
	// grpc.RegisterGreeterServer(s, &server.Server{})
	hello.RegisterHelloServer(grpcServer, &SayShetaRealService{})
}

type SayShetaRealService struct {
	hello.UnimplementedHelloServer
}

func (s *SayShetaRealService) SaySheta(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{
		Message: "Hello " + req.Name,
	}, nil
}
