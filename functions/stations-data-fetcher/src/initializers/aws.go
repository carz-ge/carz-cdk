package initializers

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"os"
)

func GetAwsConfig(region string) (aws.Config, error) {
	var otps []func(*config.LoadOptions) error
	otps = append(otps, config.WithRegion(region))
	if os.Getenv("STAGE") == "LOCAL" {
		otps = append(otps, config.WithSharedConfigProfile("carz"))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), otps...)
	return cfg, err
}
