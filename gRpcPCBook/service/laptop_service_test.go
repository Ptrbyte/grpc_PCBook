package service_test

import (
	"context"
	"gRpcPCBook/pb"
	"gRpcPCBook/sample"
	"gRpcPCBook/service"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServiceCreateLaptop(t *testing.T) {
	t.Parallel()

	NewLaptopNoID := sample.NewLaptop()
	NewLaptopNoID.Id = ""
	
	LaptopInvaildId := sample.NewLaptop()
	LaptopInvaildId.Id = "invaild_uuid"

	LaptopDuplicateId := sample.NewLaptop()
	storeDuplicateId := service.NewInMemoryLaptopStore()
	err := storeDuplicateId.Save(LaptopDuplicateId)
	require.Nil(t,err)



	testCases := []struct {
		name string
		laptop *pb.Laptop
		store service.LaptopStore
		code  codes.Code
	}{
		{
			name:"success_with_id",
			laptop: sample.NewLaptop(),
			store: service.NewInMemoryLaptopStore(),
			code : codes.OK,
		},
		{
			name:"success_no_id",
			laptop: NewLaptopNoID,
			store: service.NewInMemoryLaptopStore(),
			code : codes.OK,
		},
		{
			name:"failure_invaild_id",
			laptop: LaptopInvaildId,
			store: service.NewInMemoryLaptopStore(),
			code : codes.InvalidArgument,
		},
		{
			name:"failure_duplicate_id",
			laptop: LaptopDuplicateId,
			store: storeDuplicateId,
			code : codes.AlreadyExists,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name,func(t *testing.T){
			t.Parallel()

			req := &pb.CreateLaptopRequest{
				Laptop: tc.laptop,
			}

			service := service.NewLaptopServer(tc.store,nil,nil)
			res, err := service.CreateLaptop(context.Background(),req)
			if tc.code == codes.OK {
				require.NoError(t,err)
				require.NotNil(t,res)
				require.NotEmpty(t,res.Id)
				if len(tc.laptop.Id) > 0 {
					require.Equal(t,tc.laptop.Id,res.Id)
				}
			}else {
				require.Error(t,err)
				require.Nil(t,res)
				s, ok := status.FromError(err)
				require.True(t,ok)
				require.Equal(t,tc.code,s.Code())

			}
		})
		
	}

}