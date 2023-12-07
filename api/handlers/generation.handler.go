package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/sayantan-s/easy-send/integrations/aai"
	"github.com/sayantan-s/easy-send/integrations/cloudinary"
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
	transcriptionPollingUrl, err := aai.SetUpTranscriptions(uploadPath)
	os.Remove(uploadPath)
	
	if err != nil{
		return res.Failure(c, res.FalureTemplate{
			StatusCode: fiber.StatusInternalServerError,
			Message: "Unable to upload file",
		})
	}

	// db, _  := database.GetInstance()

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
	var data map[string]interface{}
	err := json.Unmarshal(payload, &data)
	
	if err != nil {
		return res.Failure(c, res.FalureTemplate{
			StatusCode: fiber.StatusInternalServerError,
			Message: "Unable to upload file",
		})
	}

	status:= data["status"].(string)

	if status != "completed"{
		return res.Failure(c, res.FalureTemplate{
			StatusCode: fiber.StatusInternalServerError,
			Message: "Unable to process transcription",
		})
	}

	transcriptId := data["transcript_id"].(string)

	transcriptData, err := aai.FetchTranscriptions(transcriptId)

	if err != nil {
		return res.Failure(c, res.FalureTemplate{
			StatusCode: fiber.StatusInternalServerError,
			Message: "Unable to process transcription",
		})
	}
	
	fmt.Print(transcriptData)

	return res.Success(c, res.SuccessTemplate{
		StatusCode: fiber.StatusAccepted,
		Message: "successfully generated transcription",
		Data: "Hello world",
	})
}