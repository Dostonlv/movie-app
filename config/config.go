package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"

	"github.com/spf13/cast"
)

const (
	CollectionName = "mac"
)

type Config struct {
	Environment   string
	MongoHost     string
	MongoPort     int
	MongoDatabase string
	MongoPassword string
	MongoUser     string
	LogLevel      string
	Port          string
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue

}
func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	cfg := Config{}

	cfg.Environment = cast.ToString(getOrReturnDefaultValue("ENVIRONMENT", "develop"))
	cfg.MongoHost = cast.ToString(getOrReturnDefaultValue("MONGO_HOST", os.Getenv("MONGO_HOST")))
	cfg.MongoPort = cast.ToInt(getOrReturnDefaultValue("MONGO_PORT", 27017))
	cfg.MongoDatabase = cast.ToString(getOrReturnDefaultValue("MONGO_DATABASE", os.Getenv("MONGO_DATABASE")))
	cfg.MongoUser = cast.ToString(getOrReturnDefaultValue("MONGO_USER", os.Getenv("MONGO_USER")))
	cfg.MongoPassword = cast.ToString(getOrReturnDefaultValue("MONGO_PASSWORD", os.Getenv("MONGO_PASSWORD")))
	cfg.LogLevel = cast.ToString(getOrReturnDefaultValue("LOG_LEVEL", "debug"))
	cfg.Port = cast.ToString(getOrReturnDefaultValue("PORT", os.Getenv("PORT")))

	return cfg

}
