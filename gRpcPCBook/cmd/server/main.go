package main

import (
	//"context"
	"flag"
	"fmt"
	"gRpcPCBook/pb"
	"gRpcPCBook/service"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// func unaryInterceptor(ctx context.Context,req interface{},info *grpc.UnaryServerInfo,
	// handler grpc.UnaryHandler)(interface{},error){
// 
		// log.Print("-->unary Interceptor: ",info.FullMethod)
		// return handler(ctx,req)
	// }
// 
// 
// func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, 
	// handler grpc.StreamHandler) error{
// 
		// log.Print("-->stream Interceptor: ",info.FullMethod)
		// return handler(srv,stream)
// }

const (
	secretKey = "secret"
	tokenDuration = 15 * time.Minute
)


func sendUsers(userstore service.UserStore)error {
	err := createUser(userstore,"admin1","secret","admin")
	if err != nil {
		return err
	}

	return createUser(userstore,"user1","secret","user")
}

func createUser(userStore service.UserStore,name string,password string,role string)error {
	user,err := service.NewUser(name,password,role)
	if err != nil {
		return err
	}

	return userStore.Save(user)
}

func accessibleRoles()map[string][]string {
	const laptopServerPath = "/pb.LaptopService/"

	return map[string][]string {
		laptopServerPath + "CreateLaptop":{"admin"},
		laptopServerPath + "UploadImage":{"admin"},
		laptopServerPath + "RateLaptop":{"admin","user"},
	}
}

func main() {
	port := flag.Int("port",50051,"The server port")
	flag.Parse()
	log.Printf("start server on port %d",*port)

	laptopstore := service.NewInMemoryLaptopStore()
	imagestore := service.NewDiskImageStore("img")
	ratingstore := service.NewInmemoryRatingStore()

	userstore := service.NewInMemoryUserStore()
	err :=sendUsers(userstore)
	if err != nil {
		log.Fatal("cannot send users")
	}

	jwtmanger := service.NewJWTManger(secretKey,tokenDuration)

	authServer := service.NewAuthService(userstore,jwtmanger)

	laptapServer := service.NewLaptopServer(laptopstore,imagestore,ratingstore)

	interceptor := service.NewAuthInterceptor(jwtmanger,accessibleRoles())

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)
	pb.RegisterAuthServiceServer(grpcServer,authServer)
	pb.RegisterLaptopServiceServer(grpcServer,laptapServer)
	reflection.Register(grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d",*port)
	Listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server:",err)
	}

	err2 := grpcServer.Serve(Listener)
	if err2 != nil {
		log.Fatal("cannot start server:",err2)
	}
	

}