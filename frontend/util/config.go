package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/viper"
)

type Config struct {
	Environment          string `mapstructure:"ENVIRONMENT" json:"ENVIRONMENT"`
	DBSource             string `mapstructure:"DB_SOURCE" json:"DB_SOURCE"`
	FrontendPort         string `mapstructure:"FRONTEND_PORT" json:"FRONTEND_PORT"`
	TokenSymmetricKey    string `mapstructure:"TOKEN_SYMMETRIC_KEY" json:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  string `mapstructure:"ACCESS_TOKEN_DURATION" json:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration string `mapstructure:"REFRESH_TOKEN_DURATION" json:"REFRESH_TOKEN_DURATION"`
	MainServerAddr       string `mapstructure:"MAIN_SERVER_ADDR" json:"MAIN_SERVER_ADDR"`
	InvoiceGrpcAddr      string `mapstructure:"INVOICE_GRPC_ADDR" json:"INVOICE_GRPC_ADDR"`
	FrontendAddr         string `mapstructure:"FRONTEND_ADDR" json:"FRONTEND_ADDR"`
	StripeKey            string `mapstructure:"STRIPE_KEY" json:"STRIPE_KEY"`
	StripeSecret         string `mapstructure:"STRIPE_SECRET" json:"STRIPE_SECRET"`
}

type RuntimeConfig struct {
	Config
	AccessTokenDurationParsed  time.Duration
	RefreshTokenDurationParsed time.Duration
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(ctx context.Context, path string) (Config, error) {
	environment := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if environment == "" {
		environment = "develop"
	}
	fmt.Printf("Loading config for environment: %s\n", environment)

	var config Config

	switch environment {
	case "develop":
		viper.AddConfigPath(path)
		viper.SetConfigFile(".env") // Specify exact file name
		viper.SetConfigType("env")
		_ = viper.ReadInConfig()
		viper.AutomaticEnv()
		err := viper.Unmarshal(&config)
		if err != nil {
			return config, fmt.Errorf("viper unmarshal error: %w", err)
		}
		return config, nil

	case "staging", "production":
		// secretName := fmt.Sprintf("%s/banking-system", environment)
		secretName := "banking-system"
		awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion("ap-southeast-1"))
		if err != nil {
			return config, fmt.Errorf("unable to load AWS config: %w", err)
		}
		svc := secretsmanager.NewFromConfig(awsCfg)
		secretValue, err := svc.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretName),
		})
		if err != nil {
			return config, fmt.Errorf("failed to get secret: %w", err)
		}
		err = json.Unmarshal([]byte(*secretValue.SecretString), &config)
		if err != nil {
			return config, fmt.Errorf("unmarshal secret: %w", err)
		}
		return config, nil

	default:
		return config, errors.New("invalid ENVIRONMENT: must be one of develop/staging/production")
	}
}

func NewRuntimeConfig(cfg Config) (RuntimeConfig, error) {
	atd, err := time.ParseDuration(cfg.AccessTokenDuration)
	if err != nil {
		return RuntimeConfig{}, fmt.Errorf("invalid ACCESS_TOKEN_DURATION: %w", err)
	}
	rtd, err := time.ParseDuration(cfg.RefreshTokenDuration)
	if err != nil {
		return RuntimeConfig{}, fmt.Errorf("invalid REFRESH_TOKEN_DURATION: %w", err)
	}
	return RuntimeConfig{
		Config:                     cfg,
		AccessTokenDurationParsed:  atd,
		RefreshTokenDurationParsed: rtd,
	}, nil
}
