package initializers

import (
	"log"
	"os"
)

func init() {

	region := "eu-west-1"
	config, err := GetAwsConfig(region)
	if err != nil {
		log.Fatal(err)
	}

	stage := os.Getenv("STAGE")

	var envConfig Config

	if stage != "LOCAL" {
		secretName := "stations-fetcher-prod-eu-west-1"
		envConfig, err = LoadSecrets(config, secretName)

	} else {
		envConfig, err = LoadConfig(".")
	}

	if err != nil {
		log.Fatal(err)
	}
	ConnectDB(envConfig)
}
