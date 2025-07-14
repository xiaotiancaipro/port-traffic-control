package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Groups struct {
	ID         uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
	Bandwidth  int32     `json:"bandwidth" gorm:"not null"`
	PortMaxNum int32     `json:"port_max_num" gorm:"not null"`
	Status     int8      `json:"status" gorm:"not null;default:1"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (g *Groups) TableName() string {
	return "groups"
}

func (g *Groups) BeforeCreate(_ *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}
