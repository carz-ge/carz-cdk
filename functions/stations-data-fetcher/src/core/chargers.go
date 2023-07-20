package core

import (
	"log"
	"main/src/chargers"
	"main/src/initializers"
	"main/src/models"
	"time"
)

func containsCharger(slice []models.ChargerEntity, target models.ChargerEntity) bool {
	for _, c := range slice {
		if c.Name == target.Name &&
			c.Longitude == target.Longitude &&
			c.Latitude == target.Latitude {
			return true
		}
	}
	return false
}

func GetAndUpdateChargers() {
	chargerItems, done := chargers.GetChargers()
	if done {
		return
	}

	//log.Println(chargerItems)
	var chargerEntities []models.ChargerEntity

	for _, charger := range chargerItems {
		chargerEntities = append(chargerEntities, models.ChargerEntity{
			Name:      charger.Name,
			Latitude:  charger.Latitude,
			Longitude: charger.Longitude,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	//log.Println(chargerEntities)
	var oldChargerEntities []models.ChargerEntity
	initializers.DB.Find(&oldChargerEntities)

	if len(oldChargerEntities) > len(chargerEntities) {
		log.Println("new chargers are less that is was before")
	}

	var differentChargerEntities []models.ChargerEntity

	for _, entity := range oldChargerEntities {
		if !containsCharger(chargerEntities, entity) {
			differentChargerEntities = append(differentChargerEntities, entity)
		}
	}

	log.Println(differentChargerEntities)

	if len(differentChargerEntities) > 0 {
		initializers.DB.CreateInBatches(differentChargerEntities, 100)
	}
}
