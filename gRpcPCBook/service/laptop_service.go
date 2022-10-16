package service

import (
	"bytes"
	"context"
	"errors"
	"gRpcPCBook/pb"
	"io"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maximageSize = 1 << 20

type LaptopServer struct {
	laptopStore LaptopStore
	imageStore ImageStore
	ratingStore RatingStore
	pb.UnimplementedLaptopServiceServer
}

// NewLaptopService return a new LaptopService
func NewLaptopServer(laptopstore LaptopStore,imagestore ImageStore,ratingstore RatingStore) *LaptopServer {
	return &LaptopServer{
		laptopStore:laptopstore,
		imageStore: imagestore,
		ratingStore: ratingstore,
	}
}

//CreateLaptop is a unary RPC to create a new laptop
func (service *LaptopServer) CreateLaptop(ctx context.Context, 
	req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
		laptop := req.GetLaptop()
		log.Printf("receive a create-laptop request with id: %s",laptop.Id)

		if len(laptop.Id) > 0 {
			//check if it's a valid uuid
			_, err := uuid.Parse(laptop.Id)
			if err != nil {
				return nil,status.Errorf(codes.InvalidArgument,"latop ID is not a valid UUID: %v",err)
			}
		}else {
			id, err := uuid.NewRandom()
			if err != nil {
				return nil,status.Errorf(codes.Internal,"cannot generate a new laptop ID: %v",err)
			}
			laptop.Id = id.String()
		}
		
		if err := ContextError(ctx);err != nil {
			return nil,err
		}

		//save the laptop to store
		err := service.laptopStore.Save(laptop)
		if err != nil {
			code :=codes.Internal
			if errors.Is(err,ErrAlreadyExists){
				code =codes.AlreadyExists
			}
			return nil,status.Errorf(code,"cannot save laptop to the store: %v",err)
		}
		log.Printf("save laptop with id: %s",laptop.Id)

		response := &pb.CreateLaptopResponse{
			Id: laptop.Id,
		}
		return response,nil
}

func (server *LaptopServer) SerachLaptop(req *pb.SerachLaptopRequest, 
	stream pb.LaptopService_SerachLaptopServer) error{
		filter := req.GetFilter()
		log.Printf("receive a serach-request with filter: %v",filter)

		err := server.laptopStore.Serach(
			stream.Context(),
			filter,
			func (laptop *pb.Laptop)error {
				res := &pb.SerachLaptopResponse{
					Laptop: laptop,
				}

				err := stream.Send(res)
				if err != nil {
					return err
				}

				log.Printf("send laptop with id: %v",laptop.GetId())
				return nil
			},
	)

		if err != nil {
			return status.Errorf(codes.Internal,"unexpected error :%v",err)
		}

		return nil
}

func (server *LaptopServer)UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req ,err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown,"cannot receive image info"))
	}
    
	laptopid := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("receive an upload-image request form laptop id:%v image Type:%v",laptopid,imageType)

	laptop, err2 := server.laptopStore.Find(laptopid)
	if err2 != nil {
		return logError(status.Errorf(codes.Internal,"cannot find laptop:%v",err2))
	}
	if laptop == nil {
		return logError(status.Errorf(codes.InvalidArgument,"laptop %s doesnot exist",laptopid))
	}
	
	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		//check context error
		if err := ContextError(stream.Context());err != nil {
			return err
		}

		log.Print("witing to receive more data")
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown,"cannot receive chunk data:%v",err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("receive a chunk with size:%d",size)

		imageSize += size
		if imageSize > maximageSize {
			return logError(status.Errorf(codes.InvalidArgument,"image is to large %d > %d",imageSize,maximageSize))
		}

		_, err3 := imageData.Write(chunk)
		if err3 != nil {
			return logError(status.Errorf(codes.Internal,"cannot write chunk data: %v",err3))
		}
	}

	imageID, err2 := server.imageStore.Save(laptopid, imageType, imageData)
	if err2 != nil {
		return logError(status.Errorf(codes.Internal,"cannot save image to the store:%v",err2))
	}

	resp := &pb.UploadImageResponse{
		Id: imageID,
		Size: uint32(imageSize),
	}

	err3 := stream.SendAndClose(resp)
	if err3 != nil {
		return logError(status.Errorf(codes.Unknown,"cannot send response: %v",err3))
	}
	log.Printf("saved image with id:%v, size :%d",imageID,imageSize)
	return nil
}


func (server *LaptopServer)RateLaptop(stream pb.LaptopService_RateLaptopServer) error{
	for {
		 err := ContextError(stream.Context())
		 if err!= nil {
			return err
		}

		req,err :=stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown,"cannot receive stream request:%+v",err))
		}

		laptopID := req.GetLaptopId()
		score := req.GetScores()

		log.Printf("receive a-rating-laptop id:%+v request score:%+v",laptopID,score)

		found, err := server.laptopStore.Find(laptopID)
		if err != nil {
			return logError(status.Errorf(codes.Internal,"cannot find laptop:%v",err))
		}
		if found == nil {
			return logError(status.Errorf(codes.NotFound,"laptopId :%v is not found",laptopID))
		}

		r, err := server.ratingStore.Add(laptopID, score)
		if err != nil {
			return logError(status.Errorf(codes.Internal,"cannot add rating to the store:%v",err))
		}

		resp := &pb.RateLaptopResponse{
			LaptopId: laptopID,
			RatedCount: r.Count,
			AverageScore: r.Sum / float64(r.Count),
		}

		err = stream.Send(resp)
		if err != nil {
			return logError(status.Errorf(codes.Unknown,"cannot send stream response:%v",err))
		}

		

	}

	return nil
}

func logError(err error)error {
	if err != nil {
		log.Print(err)
	}
	return err
}

func ContextError(ctx context.Context)error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled,"request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded,"deadline is exceeded"))
	default:
		return nil
	}
}