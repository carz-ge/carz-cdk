package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"main/src/core"
	"main/src/migrate"
)

func run() {
	migrate.RunMigrations()

	core.GetAndUpdateChargers()
}

func handler(ctx context.Context, event interface{}) error {
	run()
	return nil
}

func main() {
	lambda.Start(handler)
}
