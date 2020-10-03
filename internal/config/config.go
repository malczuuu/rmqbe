package config

import (
	"github.com/spf13/viper"
)

func ReadConfig() Config {
	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	viper.SetDefault("MONGODB_DATABASE", "rmqbe")
	viper.SetDefault("MONGODB_USERS_COLLECTION", "users")

	viper.BindEnv("MONGODB_URI")
	viper.BindEnv("MONGODB_DATABASE")
	viper.BindEnv("MONGODB_USERS_COLLECTION")

	return Config{
		MongoURI:             viper.GetString("MONGODB_URI"),
		MongoDatabase:        viper.GetString("MONGODB_DATABASE"),
		MongoUsersCollection: viper.GetString("MONGODB_USERS_COLLECTION"),
	}
}

type Config struct {
	MongoURI             string
	MongoDatabase        string
	MongoUsersCollection string
}
