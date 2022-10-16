package service

import (
	"context"
	"gRpcPCBook/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	userStore  UserStore
	jwtManager *JWTManager
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(userStore UserStore, jwtManager *JWTManager) *AuthService {
	return &AuthService{
		userStore:  userStore,
		jwtManager: jwtManager,
	}
}

func (server *AuthService) Login(ctx context.Context, req *pb.LoginResquest) (*pb.LoginResponse, error) {
	user, err := server.userStore.Find(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal,"cannot find user: %v",err)
	}

	if user == nil || !user.IsCorrectPassWord(req.GetPassword()){
		return nil,status.Errorf(codes.NotFound,"incorrect username/password")

	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil,status.Errorf(codes.Internal,"cannot generate access token")
	}

	resp := &pb.LoginResponse{
		AccessToken: token,
	}
	return resp,nil
}