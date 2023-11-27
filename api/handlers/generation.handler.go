package handlers

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/sayantan-s/easy-send/integrations/cloudinary"
	"github.com/sayantan-s/easy-send/integrations/openai"
	"github.com/sayantan-s/easy-send/utils/res"
)

type SuccessTemplateResponse struct{
	TranscriptionId string
	RecordedUrl string
}


func GenerateTranscriptCVS(c *fiber.Ctx) error{
	uploadDir := "./assets"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return res.Failure(c, res.FalureTemplate{
			StatusCode: fiber.StatusInternalServerError,
			Message: "Unable to process File",
		})
	}
	form, err := c.MultipartForm()
	if err != nil {
		return res.Failure(c, res.FalureTemplate{
			StatusCode: fiber.StatusInternalServerError,
			Message: "Unable to process File",
		})
	}
	file := form.File["file"][0]
	uploadPath := filepath.Join(uploadDir, file.Filename)
	
	if err := c.SaveFile(file, uploadPath); err != nil {
		return res.Failure(c, res.FalureTemplate{
			StatusCode: fiber.StatusInternalServerError,
			Message: "Unable to process file",
		})
	}

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
	transcriptionId, err := openai.Completions(uploadPath)
	os.Remove(uploadPath)

	return res.Success(c, res.SuccessTemplate{
		StatusCode: fiber.StatusCreated,
		Message: "successfully generated LinkedIn messages",
		Data: SuccessTemplateResponse{
			TranscriptionId: transcriptionId,
			RecordedUrl: reposnse.SecureURL,
		},
	})
}