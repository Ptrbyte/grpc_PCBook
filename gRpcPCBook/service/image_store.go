package service

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type ImageStore interface {
	Save(laptopid string, imageType string, imageData bytes.Buffer)( string, error)
}

type DiskImageStore struct {
	mutex sync.RWMutex
	imageFolder string
	images map[string]*ImageInfo
}

type ImageInfo struct {
	LaptopId string
	Type     string
	Path     string
}

func NewDiskImageStore(imageFolder string)*DiskImageStore{
	return &DiskImageStore{
		imageFolder: imageFolder,
		images: make(map[string]*ImageInfo),
	}
}

func (store *DiskImageStore)Save(laptopid string, imageType string,imageData bytes.Buffer)( string, error){
	imageID,err := uuid.NewRandom()
	if err != nil {
		return "",fmt.Errorf("cannot generate image id: %v",err)
	}

	imagePath := fmt.Sprintf("%s/%s%s",store.imageFolder,imageID,imageType)
	file, err1 := os.Create(imagePath)
	if err1 != nil {
		return "",fmt.Errorf("cannot create image file:%v",err)
	}

	_, err2 := imageData.WriteTo(file)
	if err2 != nil {
		return "",fmt.Errorf("cannot write image to file:%v",err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.images[imageID.String()] = &ImageInfo{
		LaptopId: laptopid,
		Type: imageType,
		Path: imagePath,
	}

	return imageID.String(),nil
}