package authgrpc

import (
	"context"
	ssov1 "github.com/YaraOduvanchik/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(
	context context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(
	context context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(
	context context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	panic("implement me")
}
