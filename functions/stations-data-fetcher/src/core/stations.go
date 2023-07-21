package core

import (
	"log"
	"main/src/initializers"
	"main/src/models"
)

func containsCharger(slice []models.AutoStationEntity, target models.AutoStationEntity) bool {
	for _, c := range slice {
		if c.Name == target.Name &&
			c.Longitude == target.Longitude &&
			c.Latitude == target.Latitude &&
			c.ProviderCode == target.ProviderCode {
			return true
		}
	}
	return false
}

func SaveStations(stationEntities []models.AutoStationEntity) {
	var oldChargerEntities []models.AutoStationEntity
	initializers.DB.Find(&oldChargerEntities)
	log.Printf("stationEntities length: %d, old stations: %d", len(stationEntities), len(oldChargerEntities))

	if len(oldChargerEntities) > len(stationEntities) {
		log.Println("new stations are less that is was before")
	}

	var differentChargerEntities []models.AutoStationEntity

	for _, entity := range stationEntities {
		if !containsCharger(oldChargerEntities, entity) {
			differentChargerEntities = append(differentChargerEntities, entity)
		}
	}

	log.Println("diff ", len(differentChargerEntities))

	if len(differentChargerEntities) > 0 {
		initializers.DB.CreateInBatches(differentChargerEntities, 100)
	}
}
