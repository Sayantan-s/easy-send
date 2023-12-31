package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	database "github.com/sayantan-s/easy-send/db"
	"github.com/sayantan-s/easy-send/integrations/aai"
	"github.com/sayantan-s/easy-send/integrations/cloudinary"
	"github.com/sayantan-s/easy-send/integrations/openai"
	"github.com/sayantan-s/easy-send/models"
	"github.com/sayantan-s/easy-send/utils/res"
)

type SuccessTemplateResponse struct{
	TranscriptId string `json:"transcriptId"`
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
	bucket := cloudinary.GetCient()
	reposnse, err := bucket.Upload.Upload(ctx, fileDoc, uploader.UploadParams{
		Folder: "media",
		ResourceType: "auto",
	})
	if err != nil{
		return res.Failure(c, res.FalureTemplate{
			StatusCode: fiber.StatusInternalServerError,
			Message: "Unable to upload file",
		})
	}
	transcriptId, err := aai.SetUpTranscriptions(uploadPath)
	os.Remove(uploadPath)
	
	if err != nil{
		return res.Failure(c, res.FalureTemplate{
			StatusCode: fiber.StatusInternalServerError,
			Message: "Unable to process file",
		})
	}

	db, _  := database.GetInstance()

	transcriptRecord := models.Transcription{
		TranscriptId: transcriptId,	
		AudioUrl: reposnse.SecureURL,
	}
	errTranscript := db.Create(&transcriptRecord).Error
	
	if errTranscript != nil{
		return res.Failure(c, res.FalureTemplate{
			StatusCode: fiber.StatusInternalServerError,
			Message: "Unable to process file",
		})
	}

	return res.Success(c, res.SuccessTemplate{
		StatusCode: fiber.StatusCreated,
		Message: "successfully generated LinkedIn messages",
		Data: SuccessTemplateResponse{
			TranscriptId: transcriptRecord.ID.String(),
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

	db, _  := database.GetInstance()
	var transcriptionRecord models.Transcription
	db.First(&transcriptionRecord, "transcript_id = ?", transcriptId)
	transcriptionRecord.Transcription = transcriptData
	db.Save(&transcriptionRecord)

	openAiPrompt := fmt.Sprintf(`
		Write 5 cold emails based on this below text:
		%s
	`, transcriptData)

	requestPayload := fmt.Sprintf(`{
		"model": "gpt-3.5-turbo",
		"messages": [
			{"role": "system", "content": "You are a helpful assistant."},
			{"role": "user", "content": "%s"}
		],
		"temperature": 0.7
	}`, openAiPrompt)

	openai.Completions(requestPayload)
	
	return res.Success(c, res.SuccessTemplate{
		StatusCode: fiber.StatusAccepted,
		Message: "successfully generated transcription",
		Data: nil,
	})
}