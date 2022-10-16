package serializer

import (
	"fmt"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
)

//WriteProtoBufToBinaryFile write Proto Buffers to JSON File
func WriteProtoBufToJSONFile(message proto.Message,filename string) error {
	data,err := ProtoBufToJSON(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto buf message to Json string:%w",err)
	}

	err = ioutil.WriteFile(filename, []byte(data),0644)
	if err != nil {
		return fmt.Errorf("cannot write JSON data to file: %w",err)
	}

	return nil
}

//WriteProtoBufToBinaryFile write Proto Buffers to Binary File
func WriteProtoBufToBinaryFile(message proto.Message,filename string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary :%w",err)
	}
	err = ioutil.WriteFile(filename,data,0644)
	if err != nil {
		return fmt.Errorf("cannot write binary data to file:%w",err)
	}

	return nil
}

//ReadProtoBufFromBinaryFile reads proto buf meaasge form binary file
func ReadProtoBufFromBinaryFile(filename string,meaasge proto.Message) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read data from file: %w",err)
	}

	err = proto.Unmarshal(data, meaasge)
	if err !=nil {
		return fmt.Errorf("cannot Unmarshal binary to proto message: %w",err)
	}
	return nil
}