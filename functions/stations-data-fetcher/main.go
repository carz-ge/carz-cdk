package main

import (
	"context"
	"log"
	"main/src/core"
	"main/src/migrate"
	"main/src/models"
	"main/src/stations"
)

func run() {
	migrate.RunMigrations()

	var stationEntities []models.AutoStationEntity

	chargers, err := stations.GetChargerEntities()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("CHARGERS ->", chargers[0].String())
		stationEntities = append(stationEntities, chargers...)
	}

	socarStations, err := stations.GetSocarStationsEntities()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("SOCAR ->", socarStations[0].String())
		stationEntities = append(stationEntities, socarStations...)
	}

	rompStations, err := stations.GetRompetrolStationsEntities()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("ROMPETROL ->", rompStations[0].String())
		stationEntities = append(stationEntities, rompStations...)
	}

	wissolStations, err := stations.GetWissolStationsEntities()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("WISSOL ->", wissolStations[0].String())
		stationEntities = append(stationEntities, wissolStations...)
	}

	gulfStations, err := stations.GetGulfStationsEntities()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("GULF ->", gulfStations[0].String())
		stationEntities = append(stationEntities, gulfStations...)
	}

	core.SaveStations(stationEntities)
}

func handler(ctx context.Context, event interface{}) error {
	run()
	return nil
}

func main() {
	run()
	//lambda.Start(handler)
}
