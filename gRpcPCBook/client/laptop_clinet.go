package client

import (
	"bufio"
	"fmt"
	"context"
	"gRpcPCBook/pb"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopClient struct{
	server pb.LaptopServiceClient
}

func NewLaptopClient(cc *grpc.ClientConn)*LaptopClient{
	server := pb.NewLaptopServiceClient(cc)
	return &LaptopClient{server}
}


func (laptopClient *LaptopClient)CreateLaptop(laptop *pb.Laptop){
	//laptop := sample.NewLaptop()
    //laptop.Id = ""
    req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	//set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err2 := laptopClient.server.CreateLaptop(ctx, req)
	if err2 != nil {
		st,ok := status.FromError(err2)
		if ok && st.Code() == codes.AlreadyExists {
			log.Println("laptop already exists")
		}else {
			log.Fatal("cannot create laptop:",err2)
		}
		return
	}
	log.Printf("create laptop with id:%v",res.Id)
}

func (laptopClient *LaptopClient)SerachLaptop(filter *pb.Filter){
	log.Print("serach filter:" ,filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SerachLaptopRequest{Filter: filter}
	stream, err := laptopClient.server.SerachLaptop(ctx, req)
	if err != nil {
		log.Fatal("cannot serach laptop:",err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive response:",err)
		}

		laptop := res.GetLaptop()
		log.Print("- found: ", laptop.GetId())
		log.Print(" +brand: ", laptop.GetBrand())
		log.Print(" +name: ", laptop.GetName())
		log.Print(" +cpu cores: ", laptop.GetCpu().GetNumberCores())
		log.Print(" +cpu min ghz: ", laptop.GetCpu().GetMinGhz())
		log.Print(" +ram: ", laptop.GetRam().GetValue(), laptop.GetRam().GetUnit())
		log.Print(" +price: ", laptop.GetPriceUsd())
	}

}

func (laptopClient *LaptopClient)UpLoadImage(laptopID string,imagePath string){
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file:",err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.server.UploadImage(ctx)
	if err != nil {
		log.Fatal("cannot upLoad image:",err)
	}

	req := &pb.UpLoadImageResquest{
		Data: &pb.UpLoadImageResquest_Info{
			Info: &pb.ImageInfo{
				LaptopId: laptopID,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}
	err2 := stream.Send(req)
	if err2 != nil {
		log.Fatal("cannot send image info:",err2)
	}
	//Todo

	reader := bufio.NewReader(file)
	buffer := make([]byte,1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer:",err)
		}


		req := &pb.UpLoadImageResquest{
			Data: &pb.UpLoadImageResquest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err3 := stream.Send(req)
		if err3 != nil {
			err2 := stream.RecvMsg(nil)
			log.Fatal("cannot send chunk to server:",err3,err2)
		}
	}
	res, err3 := stream.CloseAndRecv()
	if err3 != nil {
		log.Fatal("cannot receive response:",err3)
	}

	log.Printf("image Upload with id:%+v, size:%d",res.GetId(),res.GetSize())

}

func (laptopClient *LaptopClient)RateLaptop(laptopIDs []string,scores []float64)error{
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.server.RateLaptop(ctx)
	if err != nil {
		return fmt.Errorf("connot rate laptop:%v",err)
	}

	waitResponse := make(chan error)
	//receive stream
	go func(){
		for {
			res, err:= stream.Recv()
			if err == io.EOF {
				log.Print("no more response")
				waitResponse <- nil
				return
			}

			if err != nil {
				waitResponse <- fmt.Errorf("cannot receive stream response:%v",err)
				return 
			}
			log.Print("receive response: ",res)
		}
	}()
	
	//send request
	for i,laptopid := range laptopIDs{

		req := &pb.RateLaptopRequest{
			LaptopId: laptopid,
			Scores: scores[i],
		}
		err := stream.Send(req)
		if err != nil {
			return fmt.Errorf("cannot send stream request:%v - %v",err,stream.RecvMsg(nil))
		}

		log.Print("send request: ",req)

	}

	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("cannot close send: %v",err)
	}

	err = <-waitResponse
	return err
}