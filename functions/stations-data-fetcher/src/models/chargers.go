package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"
)

type ChargerEntity struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Latitude  string    `gorm:"type:varchar(255);not null"`
	Longitude string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ch *ChargerEntity) String() string {
	marshal, err := json.Marshal(ch)
	if err != nil {
		log.Fatal(err)
	}
	return string(marshal)

}

func (ch *ChargerEntity) TableName() string {
	return "chargers"
}
