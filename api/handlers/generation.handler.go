package handlers

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/sayantan-s/easy-send/integrations/cloudinary"
	"github.com/sayantan-s/easy-send/utils/res"
)


func GenerateTranscriptCVS(c *fiber.Ctx) error{
	fileHeader, _ := c.FormFile("file")
	fileDoc, _ := fileHeader.Open()

	ctx := context.Background()
	reposnse, err := cloudinary.GetCient().Upload.Upload(ctx, fileDoc, uploader.UploadParams{
		Folder: "media",
		ResourceType: "auto",
	})
	if err != nil{
		return res.Failure(c, res.FalureTemplate{
			StatusCode: fiber.StatusInternalServerError,
			Message: "Unable to upload file",
		})
	}
	return res.Success(c, res.SuccessTemplate{
		StatusCode: fiber.StatusCreated,
		Message: "successfully generated LinkedIn messages",
		Data: reposnse.SecureURL,
	})
}