package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"log"
	"time"
)

type ServiceType struct {
	Id      string
	Title   string
	TitleEn string
	Image   string
	Code    string
}

type AutoStationEntity struct {
	gorm.Model
	Id           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdByProvider string    `gorm:"type:varchar(255);not null"`

	StationType  string `gorm:"type:varchar(255);not null"`
	ProviderCode string `gorm:"type:varchar(255);not null"`

	TextHtml   string `gorm:"type:varchar(255)"`
	TextHtmlEn string `gorm:"type:varchar(255)"`

	Name   string `gorm:"type:varchar(255);not null"`
	NameEn string `gorm:"type:varchar(255);not null"`

	Description   string
	DescriptionEn string

	Active bool

	Latitude  string `gorm:"type:varchar(255);not null"`
	Longitude string `gorm:"type:varchar(255);not null"`
	Region    string
	RegionEn  string
	City      string
	CityEn    string
	Address   string
	AddressEn string

	Picture string

	ProductTypes datatypes.JSON // fuels
	ObjectTypes  datatypes.JSON // buildings
	PaymentTypes datatypes.JSON
	ServiceTypes datatypes.JSON // other services

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ch *AutoStationEntity) String() string {
	marshal, err := json.Marshal(ch)
	if err != nil {
		log.Fatal(err)
	}
	return string(marshal)

}

func (ch *AutoStationEntity) TableName() string {
	return "auto_stations"
}
