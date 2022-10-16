package service

import (
	"context"
	"errors"
	"fmt"
	"gRpcPCBook/pb"
	"log"
	"sync"

	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("record already exists")
//LaptopStore is an interface to store laptop
type LaptopStore interface {
	//Save laptop to the store
	Save(laptop *pb.Laptop) error
	//Find a laptop by ID
	Find(id string)(*pb.Laptop,error)
	//Serach serachs for laptops with filter,returns one by one via the found func
	Serach(ctx context.Context,filter *pb.Filter,found func(laptop *pb.Laptop)error)error
}

//InMemoryLaptopStore store laptop to in memeory
type InMemoryLaptopStore struct{
	mutex sync.RWMutex
	data map[string]*pb.Laptop
}

func NewInMemoryLaptopStore()*InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

func (store *InMemoryLaptopStore)Find(id string)(*pb.Laptop,error){
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil,nil
	}
	
	//deep copy
	return deepCopy(laptop)
}

//Serach serachs for laptops with filter,returns one by one via the found func
func (store *InMemoryLaptopStore)Serach(ctx context.Context,
	filter *pb.Filter,found func(laptop *pb.Laptop)error)error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	log.Print("filter:",filter)
	for _,laptop := range store.data {

		log.Print("Checking laptop id: ", laptop.GetId())

		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Print("context is cancelled")
			return errors.New("context is cancelled")
		}

		if isQualified(filter,laptop){

			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}

			err = found(other)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isQualified(filter *pb.Filter,laptop *pb.Laptop) bool {
	
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}
	
	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}
	
	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}
	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()){
		return false
	}

	return true
}

func toBit(memory *pb.Memory) uint64 {
	val := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return val
	case pb.Memory_BYTE:
		return val << 3 //8=2^3
	case pb.Memory_KILOBYTE:
		return val << 13 //1024 * 8
	case pb.Memory_MEGABYTE:
		return val << 23
	case pb.Memory_GIGABYTE:
		return val << 33
	case pb.Memory_TERABYE:
		return val << 43
	default:
		return 0
	}
}

//deep copy
func deepCopy(laptop *pb.Laptop)(*pb.Laptop,error){
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
	return nil,fmt.Errorf("cannot copy laptop data: %v",err)
	}
	return other,nil
}
//Save laptop to the store
func (store *InMemoryLaptopStore)Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	//deep copy
	other ,err:= deepCopy(laptop)
	if err != nil {
		return err
	}

	store.data[other.Id] = other
	return nil
	
}

