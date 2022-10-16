package serializer

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

//ProtoBufJson converts proto buff message to JSON string
func ProtoBufToJSON(message proto.Message)(string,error) {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}
	return marshaler.MarshalToString(message)

}