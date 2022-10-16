package sample

import (
	"gRpcPCBook/pb"
	"github.com/golang/protobuf/ptypes"
)

//NewKeyboard return a new Keyboard sample
func NewKeyboard() *pb.Keyboard {
	keyboard := &pb.Keyboard{
		Layout: randomKeyboardLayout(),
		Backlit: randomBool(),
	}
	return keyboard

}

//NewCPU return a new CPU sample
func NewCPU() *pb.CPU {

	brand := randCPUBrand()
	name := randCPUName(brand)
	numberCores := randomInt(2,8)
	numberTreads := randomInt(numberCores,12)
	minghz := randomFloat64(2.0,3.5)
	maxghz := randomFloat64(minghz,5.0)

	cpu := &pb.CPU{
		Brand: brand,
		Name:  name,
		NumberCores: uint32(numberCores),
		NumberThreads: uint32(numberTreads),
		MinGhz: minghz,
		MaxGhz: maxghz,
	}

	return cpu
}

//NewGPU return a new sample GPU
func NewGPU()*pb.GPU {

	brand := randGPUBrand()
	name := randGPUName(brand)
	minGhz := randomFloat64(1.0,1.5)
	maxGhz := randomFloat64(minGhz,2.0)

	memory := &pb.Memory{
		Value: uint64(randomInt(2,6)),
		Unit: pb.Memory_GIGABYTE,
	}

	gpu := &pb.GPU{
		Brand: brand,
		Name:  name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: memory,

	}
	
	return gpu
}

//NewRAM return a new sample Ram 
func NewRAM() *pb.Memory {
	ram := &pb.Memory{
		Value: uint64(randomInt(4,12)),
		Unit: pb.Memory_GIGABYTE,
	}
	return ram
}

//NewSSD return a new sample SSD
func NewSSD()*pb.Storage {
	ssd := &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(128,1024)),
			Unit: pb.Memory_GIGABYTE,
		},
	}
	return ssd
}

//NewHDD return a new smaple HDD 
func NewHDD()*pb.Storage {
	hhd := &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(1,6)),
			Unit: pb.Memory_TERABYE,
		},
	}
	return hhd
}


//NewScreen return a new smaple Screen
func NewScreen()*pb.Screen{
	
	screen :=&pb.Screen{
		SizeInch: randomFloat32(13,17),
		Resolution: randomSceenResolution(),
		Panel: randomPanel(),
		Multitouch: randomBool(),

	}
	return screen
}

//Newlaptop return a new smaple Laptop
func NewLaptop()*pb.Laptop{

	brand := randLaptopBrand()
	name := randLaptopName(brand)

	laptop := &pb.Laptop{
		Id: randomID(),
		Brand: brand,
		Name: name,
		Cpu: NewCPU(),
		Ram: NewRAM(),
		Gpus: []*pb.GPU{NewGPU()},
		Storages: []*pb.Storage{NewSSD(),NewHDD()},
		Screen: NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(1.0,3.0),
		},
		PriceUsd: randomFloat64(1500,3000),
		ReleaseYear: uint32(randomInt(2015,2019)),
		UpdatedAt: ptypes.TimestampNow(),
	}
	return laptop
}

