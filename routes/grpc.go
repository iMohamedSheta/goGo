package routes

import (
	"context"
	sheta "imohamedsheta/gocrud/grpc/v1/sheta/proto"

	"google.golang.org/grpc"
)

func LoadGrpcRoutes(grpcServer *grpc.Server) {
	// Here you can register your grpc routes
	// e.g.
	// grpc.RegisterGreeterServer(s, &server.Server{})
	sheta.RegisterSayShetaServer(grpcServer, &SayShetaRealService{})
}

type SayShetaRealService struct {
	sheta.UnimplementedSayShetaServer
}

func (s *SayShetaRealService) SaySheta(ctx context.Context, req *sheta.SayShetaRequest) (*sheta.SayShetaResponse, error) {
	return &sheta.SayShetaResponse{
		Message: "Hello " + req.Name,
	}, nil
}
