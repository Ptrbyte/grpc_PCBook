package client

import (
	"context"
	"gRpcPCBook/pb"
	"time"

	"google.golang.org/grpc"
)

type AuthClient struct {
	server pb.AuthServiceClient
	username string
	password string
}

func NewAuthClient(cnn *grpc.ClientConn,name string, password string)*AuthClient{
	server := pb.NewAuthServiceClient(cnn)
	return &AuthClient{server,name,password}
}

func (client *AuthClient)Login()(string,error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LoginResquest{
		Username: client.username,
		Password: client.password,
	}

	resp, err := client.server.Login(ctx, req)
	if err != nil {
		return "",err
	}

	return resp.GetAccessToken(),nil
}