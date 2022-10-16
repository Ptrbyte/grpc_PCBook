package main

import (
	"fmt"
	"strings"
	"time"

	//"bytes"

	"flag"
	"gRpcPCBook/client"
	"gRpcPCBook/pb"
	"gRpcPCBook/sample"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)



func testCreateLaptop(laptopclient *client.LaptopClient){
	laptopclient.CreateLaptop(sample.NewLaptop())
}

func testSerachLaptop(laptopClient *client.LaptopClient){
	//create 10 laptops
	for i := 0;i < 10 ;i++ {
		laptopClient.CreateLaptop(sample.NewLaptop())
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz: 2.0,
		MinRam: &pb.Memory{
			Value: 8,
			Unit: pb.Memory_GIGABYTE,
		},
	}
	laptopClient.SerachLaptop(filter)
	
}

func testUploadImage(laptopClient *client.LaptopClient){
	laptop := sample.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	laptopClient.UpLoadImage(laptop.GetId(),"temp/laptop.jpg")

}

func testRatelaptop(laptopClient *client.LaptopClient){
	n := 3
	laptopIDs := make([]string,n)

	for i := 0; i < n; i++{
		laptop :=sample.NewLaptop()
		laptopIDs[i] =laptop.GetId()
		laptopClient.CreateLaptop(laptop)
	}

	sores := make([]float64,n)

	for {
		fmt.Print("Rate laptop (y/n)?")
		var answer string
		fmt.Scan(&answer)
		
		if strings.ToLower(answer) != "y"{
			break
		}

		for i := 0 ; i < n ; i ++ {
			sores[i] = sample.RandomLaptopScore()
		}

		err := laptopClient.RateLaptop(laptopIDs, sores)
		if err != nil {
			log.Fatal(err)
		}
	}


}

const (
	username = "admin1"
	//username = "user1"
	password = "secret"
	refreshDuration = 30 *time.Second
)


func authMethods()map[string]bool {
	const laptopServerPath = "/pb.LaptopService/"
	return map[string]bool {
		laptopServerPath + "CreateLaptop":true,
		laptopServerPath + "UploadImage":true,
		laptopServerPath + "RateLaptop":true,
	}
}

func main() {
	serveraddress:= flag.String("addr", "localhost:50051", "the address to connect to")
	flag.Parse()
	log.Printf("Dial server %s",*serveraddress)

	cc1, err := grpc.Dial(*serveraddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial cc1 server:",err)
	}
	defer cc1.Close()

	authCilent := client.NewAuthClient(cc1,username,password)

	interceptor, err := client.NewAuthInterceptor(authCilent, authMethods(), refreshDuration)
	if err != nil {
		log.Fatalf("cannot create auth interceptor: %v",err)
	}

	cc2, err := grpc.Dial(*serveraddress, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),)
	if err != nil {
		log.Fatal("cannot dial cc2 server:",err)
	}
	defer cc2.Close()

	laptopCilent := client.NewLaptopClient(cc2)
	//testSerachLaptop(laptopCilent)

	//testUploadImage(laptopCilent)
	testRatelaptop(laptopCilent)

}
