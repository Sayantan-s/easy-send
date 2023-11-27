package cloudinary

import (
	"sync"

	storage "github.com/cloudinary/cloudinary-go/v2"
	"github.com/sayantan-s/easy-send/config"
)

var once sync.Once
var instance *storage.Cloudinary

func GetCient() *storage.Cloudinary{
	once.Do(func() {
		CLOUDINARY_CLOUD_NAME := config.GetConfig("CLOUDINARY_CLOUD_NAME")
		CLOUDINARY_API_KEY:= config.GetConfig("CLOUDINARY_API_KEY")
		CLOUDINARY_SECRET_KEY := config.GetConfig("CLOUDINARY_SECRET_KEY")
		cld, _ := storage.NewFromParams(
			CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY, CLOUDINARY_SECRET_KEY,
		)
		instance = cld
	})
	return instance
}