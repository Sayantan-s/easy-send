package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Transcription struct{
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primary_key" json:"id"`
	Transcription string     `gorm:"type:text"`
	AudioUrl      string     `gorm:"type:text"`
	TranscriptId  string 
	CreatedAt     time.Time  
	UpdatedAt     time.Time  
}

func (Transcription) TableName() string {
	return "transcript_dtl"
}

func AutoMigrateTranscripts(db *gorm.DB) {
	db.AutoMigrate(&Transcription{})
}