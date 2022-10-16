package service

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	jwtManger       *JWTManager
	accessibleRoles map[string][]string
}

func NewAuthInterceptor(jwtManger *JWTManager, accessibleRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManger:       jwtManger,
		accessibleRoles: accessibleRoles,
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor{
	return func(ctx context.Context,
				req interface{},
				info *grpc.UnaryServerInfo,
				handler grpc.UnaryHandler)(
				interface{},
				error){
					log.Print("-->unary Interceptor: ",info.FullMethod)
					err := interceptor.authorize(ctx,info.FullMethod)
					if err != nil {
						return nil,err
					}
					return handler(ctx,req)
		}
}

func (interceptor *AuthInterceptor)Stream()grpc.StreamServerInterceptor{
	return func (srv interface{}, 
		stream grpc.ServerStream, 
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error{
			log.Print("-->stream Interceptor: ",info.FullMethod)

			err := interceptor.authorize(stream.Context(),info.FullMethod)
			if err != nil {
				return err
			}
		
			return handler(srv,stream)
	}
}

func (interceptor *AuthInterceptor)authorize(ctx context.Context,method string)error {
	accessibleRoles,ok := interceptor.accessibleRoles[method]
	if !ok {
		//every one can access
		return nil
	}

	md,ok :=metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated,"metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0{
		return status.Errorf(codes.Unauthenticated,"authorization token is not provided")
	}
	accesstoken := values[0]
	claims,err := interceptor.jwtManger.Verify(accesstoken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated,"access token is invalid:%v",err)
	}

	for _,role := range accessibleRoles{
		if role == claims.Role {
			return nil
		}
	}
	return status.Errorf(codes.PermissionDenied,"no permission an access this RPC")
}