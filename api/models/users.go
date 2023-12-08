package models

import (
	"time"

	"github.com/google/uuid"
)

type Transcription struct{
	ID            uuid.UUID  `gorm:"column:id"`
	TranscriptId  string  `gorm:"column:transcriptId"`
	Transcription string     `gorm:"column:transcript;type:text"`
	AudioUrl      string     `gorm:"column:audioUrl;type:text"`
	CreatedAt     time.Time  `gorm:"column:createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt"`
}