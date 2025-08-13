package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ports struct {
	ID        uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
	GroupID   uuid.UUID `json:"group_id" gorm:"type:varchar(36);not null;index"`
	Port      int32     `json:"port" gorm:"not null"`
	Status    int8      `json:"status" gorm:"not null;default:-1"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (p *Ports) TableName() string {
	return "ports"
}

func (p *Ports) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
