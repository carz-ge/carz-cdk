package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
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
		stationEntities = append(stationEntities, chargers...)
	}

	log.Println("CHARGERS ->", chargers[0].String())

	socarStations, err := stations.GetSocarStationsEntities()
	if err != nil {
		log.Println(err)
	} else {
		stationEntities = append(stationEntities, socarStations...)
	}
	log.Println("SOCAR ->", socarStations[0].String())

	rompStations, err := stations.GetRompetrolStationsEntities()
	if err != nil {
		log.Println(err)
	} else {
		stationEntities = append(stationEntities, rompStations...)
	}
	log.Println("ROMPETROL ->", rompStations[0].String())

	wissolStations, err := stations.GetWissolStationsEntities()
	if err != nil {
		log.Println(err)
	} else {
		stationEntities = append(stationEntities, wissolStations...)
	}
	log.Println("WISSOL ->", wissolStations[0].String())

	gulfStations, err := stations.GetGulfStationsEntities()
	if err != nil {
		log.Println(err)
	} else {
		stationEntities = append(stationEntities, gulfStations...)
	}
	log.Println("GULF ->", gulfStations[0].String())

	core.SaveStations(stationEntities)
}

func handler(ctx context.Context, event interface{}) error {
	run()
	return nil
}

func main() {
	// 	run()
	lambda.Start(handler)
}
