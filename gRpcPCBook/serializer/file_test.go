package serializer_test

import (
	"gRpcPCBook/pb"
	"gRpcPCBook/sample"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"gRpcPCBook/serializer"
)

//testing Write Proto Buf message to Binary File
func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryfile :="../temp/laptop.bin"
	jsonfile :="../temp/laptop.json"

	laptop1 := sample.NewLaptop()
	err := serializer.WriteProtoBufToBinaryFile(laptop1, binaryfile)
	require.NoError(t,err)

	laptop2 := &pb.Laptop{}
	err2 := serializer.ReadProtoBufFromBinaryFile(binaryfile, laptop2)
	require.NoError(t,err2)
	require.True(t,proto.Equal(laptop1,laptop2))

	err3 := serializer.WriteProtoBufToJSONFile(laptop1, jsonfile)
	require.NoError(t,err3)

}

