package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/sayantan-s/easy-send/integrations/cloudinary"
	"github.com/sayantan-s/easy-send/integrations/openai"
	"github.com/sayantan-s/easy-send/utils/res"
)

type SuccessTemplateResponse struct{
	TranscriptPollingUrl string `json:"transcriptPollingUrl"`
	RecordedUrl string `json:"recordedUrl"`
}


func InitiateTranscriptions(c *fiber.Ctx) error{
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
	transcriptionPollingUrl, err := openai.SetUpTranscriptions(uploadPath)
	os.Remove(uploadPath)
	
	if err != nil{
		return res.Failure(c, res.FalureTemplate{
			StatusCode: fiber.StatusInternalServerError,
			Message: "Unable to upload file",
		})
	}

	return res.Success(c, res.SuccessTemplate{
		StatusCode: fiber.StatusCreated,
		Message: "successfully generated LinkedIn messages",
		Data: SuccessTemplateResponse{
			TranscriptPollingUrl: transcriptionPollingUrl,
			RecordedUrl: reposnse.SecureURL,
		},
	})
}

func GetGeneratedTranscripts(c *fiber.Ctx) error{
	payload := c.Body()
	fmt.Println(payload)
	return res.Success(c, res.SuccessTemplate{
		StatusCode: fiber.StatusAccepted,
		Message: "successfully generated LinkedIn messages",
		Data: "Hello world",
	})
}