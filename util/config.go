package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDrive                string        `mapstructure:"DB_DRIVER"`
	DBSource               string        `mapstructure:"DB_SOURCE"`
	MigrationURL           string        `mapstructure:"MIGRATION_URL"`
	MilvusAddr             string        `mapstructure:"MILVUS_ADDR"`
	RedisAddress           string        `mapstructure:"REDIS_ADDRESS"`
	VectorGrpcAddr         string        `mapstructure:"VECTOR_GRPC_ADDR"`
	ServerAddress          string        `mapstructure:"SERVER_ADDR"`
	TokenSymmetric         string        `mapstructure:"TOKEN_SYMMETRIC"`
	VectorDBCollectionName string        `mapstructure:"VECTOR_DB_COLLECTION_NAME"`
	ZillisEndpoint         string        `mapstructure:"ZILLIS_ENDPOINT"`
	ZillisAPIKey           string        `mapstructure:"ZILLIS_API_KEY"`
	AccessTokenDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	OpenAIAPIKey           string        `mapstructure:"OPENAI_API_KEY"`
	OpenAIURL              string        `mapstructure:"OPENAI_URL"`
	EmailSenderName        string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress     string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword    string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	RefreshTokenDuration   time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
