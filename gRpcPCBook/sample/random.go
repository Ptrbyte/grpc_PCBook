package sample

import (
	"gRpcPCBook/pb"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func init(){
	rand.Seed(time.Now().UnixNano())
}

func randomKeyboardLayout() pb.Keyboard_Layout{
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}

}

func randCPUBrand() string {
	return randStringFormSet("intel","AMD") 
}

func randStringFormSet(a...string) string{
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}

func randCPUName(brand string) string {
	if brand == "intel" {
		return randStringFormSet(
			"Xeon E-2286M",
			"Core i9-9980HK",
			"Core i7-9750H",
			"Core i5-9400F",
			"Core i3-1005G1",
		)
	}
	return randStringFormSet(
		"Ryzen 7 PRO 2700U",
		"Ryzen 5 PRO 3500U",
		"Ryzen 3 PRO 3200GE",
	)
}

func  randomInt(min ,max int)int {
	return rand.Intn( max- min + 1)+min
}

func randomFloat64(min,max float64) float64 {
	return min+rand.Float64()*(max - min)

}

func randGPUBrand()string {
	return randStringFormSet("NVIDIA","AMD")
}

func randGPUName(brand string)string {
	if brand == "NVIDIA" {
		return randStringFormSet(
			"RTX 2060",
			"RTX 2070",
			"RTX 1660-Ti",
			"RTX 1070",
		)
	}
	return randStringFormSet(
		"RX 590",
		"RX 580",
		"RX 5700-XT",
		"RX Vega-56",
	)
}


func randomFloat32(min float32,max float32)float32 {
	return min + rand.Float32() * (max-min)
}

func randomPanel()pb.Screen_Panel {

	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

func randomSceenResolution() *pb.Screen_Resolution {

	height := randomInt(1080,4320)
	width := height * 16 / 9
	resolution := &pb.Screen_Resolution{
		Height: uint32(height),
		Width: uint32(width),
		}
	return resolution
}

func randomID() string {
	return uuid.New().String()
}

func randLaptopBrand()string {
	return randStringFormSet("Apple","Dell","Lenovo")
}

func randLaptopName(brand string) string {

	switch brand {
	case "Apple":
		return randStringFormSet("Macbook Air","Macbook Pro")
	case "Dell":
		return randStringFormSet("Latitude","Vostro","XPS","Alieware")
	default:
		return randStringFormSet("Thinkpad X1","Thinkpad P1","Thinkpad P53")
	}
}

func randomBool() bool {

	return rand.Intn(2) == 1
}

func RandomLaptopScore()float64{
	return float64(randomInt(1,10))
}
	


