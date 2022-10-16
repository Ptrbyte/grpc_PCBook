package service_test

import (
	"bufio"
	"context"
	"fmt"
	"gRpcPCBook/pb"
	"gRpcPCBook/sample"
	"gRpcPCBook/serializer"
	"gRpcPCBook/service"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopstore := service.NewInMemoryLaptopStore()
	serverAdr := startTestLaptopService(t, laptopstore, nil,nil)
	laptopCilent := newLaptopCilent(t, serverAdr)

	laptop := sample.NewLaptop()
	expectedID := laptop.Id

	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	res, err := laptopCilent.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	//check that the laptop is saved to the store
	other, err1 := laptopstore.Find(res.Id)
	require.NoError(t, err1)
	require.NotNil(t, other)

	//check that the saved laptop is the same as the one we send
	requireSameLapatop(t, laptop, other)
}

func startTestLaptopService(t *testing.T, laptopstore service.LaptopStore, imagestore service.ImageStore,ratingstore service.RatingStore) string {

	laptopServer := service.NewLaptopServer(laptopstore, imagestore, ratingstore)
	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go grpcServer.Serve(listener)
	return listener.Addr().String()
}

func newLaptopCilent(t *testing.T, serveraddress string) pb.LaptopServiceClient {
	cc, err := grpc.Dial(serveraddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(cc)
}

func requireSameLapatop(t *testing.T, laptop1 *pb.Laptop, laptop2 *pb.Laptop) {
	s1, err := serializer.ProtoBufToJSON(laptop1)
	require.NoError(t, err)

	s2, err := serializer.ProtoBufToJSON(laptop2)
	require.NoError(t, err)
	require.Equal(t, s1, s2)
}

func TestClientSerachLaptop(t *testing.T) {
	t.Parallel()

	filter := &pb.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.2,
		MinRam: &pb.Memory{
			Value: 8,
			Unit:  pb.Memory_GIGABYTE,
		},
	}

	laptopStore := service.NewInMemoryLaptopStore()
	expectedIDs := make(map[string]bool)
	for i := 0; i < 6; i++ {

		laptop := sample.NewLaptop()
		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Ram = &pb.Memory{Value: 4096, Unit: pb.Memory_MEGABYTE}
		case 4:
			laptop.PriceUsd = 1999
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.5
			laptop.Cpu.MaxGhz = 4
			laptop.Ram = &pb.Memory{Value: 16, Unit: pb.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumberCores = 6
			laptop.Cpu.MinGhz = 2.8
			laptop.Cpu.MaxGhz = 5.0
			laptop.Ram = &pb.Memory{Value: 64, Unit: pb.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		}
		err := laptopStore.Save(laptop)
		require.NoError(t, err)
	}

	serverAdr := startTestLaptopService(t, laptopStore, nil,nil)
	laptopClient := newLaptopCilent(t, serverAdr)

	req := &pb.SerachLaptopRequest{Filter: filter}
	stream, err := laptopClient.SerachLaptop(context.Background(), req)
	require.NoError(t, err)

	found := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.Contains(t, expectedIDs, res.GetLaptop().GetId())
		found += 1
	}
	require.Equal(t, len(expectedIDs), found)

}

func TestClientUploadImage(t *testing.T) {
	t.Parallel()

	testImageFolder := "../temp"

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore(testImageFolder)

	laptop := sample.NewLaptop()
	err := laptopStore.Save(laptop)
	require.NoError(t, err)

	serverAdr := startTestLaptopService(t, laptopStore, imageStore,nil)
	laptopClient := newLaptopCilent(t, serverAdr)

	imagePath := fmt.Sprintf("%s/laptop.jpg", testImageFolder)
	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer file.Close()

	stream, err := laptopClient.UploadImage(context.Background())
	require.NoError(t, err)

	imageType := filepath.Ext(imagePath)
	req := &pb.UpLoadImageResquest{
		Data: &pb.UpLoadImageResquest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptop.GetId(),
				ImageType: imageType,
			},
		},
	}

	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		size += n

		req := &pb.UpLoadImageResquest{
			Data: &pb.UpLoadImageResquest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotZero(t, res.GetId())
	require.EqualValues(t, size, res.GetSize())
	savedImagePath := fmt.Sprintf("%s/%s%s", testImageFolder, res.GetId(), imageType)
	require.FileExists(t, savedImagePath)
	//require.NoError(t, os.Remove(savedImagePath))
}

func TestClientRateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := service.NewInMemoryLaptopStore()
	ratingStore := service.NewInmemoryRatingStore()

	laptop := sample.NewLaptop()
	err := laptopStore.Save(laptop)
	require.NoError(t, err)

	serverAdr := startTestLaptopService(t, laptopStore, nil,ratingStore)
	laptopClient := newLaptopCilent(t, serverAdr)

	stream, err := laptopClient.RateLaptop(context.Background())
	require.NoError(t,err)

	scores := []float64{8,7.5,10}
	avages := []float64{8,7.75,8.5}

	n := len(scores)

	for i := 0;i < n ;i++ {
		req := &pb.RateLaptopRequest{
			LaptopId: laptop.GetId(),
			Scores: scores[i],
		}

		err := stream.Send(req)
		require.NoError(t,err)	
	}

	err = stream.CloseSend()
	require.NoError(t,err)

	for idx := 0;;idx ++ {
		res, err := stream.Recv()
		if err == io.EOF {
			require.Equal(t,n,idx)
			return
		}

		require.NoError(t, err)
		require.Equal(t, laptop.GetId(), res.GetLaptopId())
		require.Equal(t, uint32(idx+1), res.GetRatedCount())
		require.Equal(t, avages[idx], res.GetAverageScore())
	}

}