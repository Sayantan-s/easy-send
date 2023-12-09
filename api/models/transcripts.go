package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Transcription struct{
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primary_key" json:"id"`
	Transcription string     `gorm:"type:text" json:"transcript"`
	AudioUrl      string     `gorm:"type:text" json:"audioUrl"`
	TranscriptId  string  	 `json:"transcriptId"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

func AutoMigrateTranscripts(db *gorm.DB) {
	db.AutoMigrate(&Transcription{})
}