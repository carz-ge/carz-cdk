package initializers

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DBUrl string `mapstructure:"POSTGRES_URL" json:"POSTGRES_URL"`
	//`json:"POSTGRES_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func LoadSecrets(awsConfig aws.Config, secretName string) (config Config, err error) {
	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(awsConfig)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString
	log.Println(secretString)
	err = json.Unmarshal([]byte(secretString), &config)
	// Your code goes here.
	return
}
