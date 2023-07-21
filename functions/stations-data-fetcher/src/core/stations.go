package core

import (
	"log"
	"main/src/initializers"
	"main/src/models"
	"main/src/stations"
	"time"
)

func containsCharger(slice []models.AutoStationEntity, target models.AutoStationEntity) bool {
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
	chargerItems, err := stations.GetChargers()
	if err != nil {
		return
	}

	var chargerEntities []models.AutoStationEntity

	for _, charger := range chargerItems {
		chargerEntities = append(chargerEntities, models.AutoStationEntity{
			Name:      charger.Name,
			Latitude:  charger.Latitude,
			Longitude: charger.Longitude,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	SaveStations(chargerEntities)
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
