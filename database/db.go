package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"movie/config"
	"movie/models"
	"movie/pkg/logger"
)

var (
	log logger.Logger
	cfg config.Config
)

func InitDB() models.MongoInstance {
	cfg = config.Load()
	log = logger.New(cfg.LogLevel, "Movie")
	log.Info("main: SQLConfig",
		logger.String("Host", cfg.MongoHost),
		logger.Int("Port", cfg.MongoPort),
		logger.String("Database", cfg.MongoDatabase),
	)

	//credential := options.Credential{
	//	Username: cfg.MongoUser,
	//	Password: cfg.MongoPassword,
	//}
	mongoString := fmt.Sprintf("mongodb://%s:%d", cfg.MongoHost, cfg.MongoPort)

	mongoConn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoString))
	if err != nil {
		log.Error("error to connect to mongo database", logger.Error(err))
	}
	connDB := mongoConn.Database(cfg.MongoDatabase)
	MG := models.MongoInstance{
		Client: mongoConn,
		Db:     mongoConn.Database(cfg.MongoDatabase),
	}

	// email is unique
	index := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err = connDB.Collection("user").Indexes().CreateOne(context.Background(), index)
	log.Info("Connected to MongoDB", logger.Any("database: ", connDB.Name()))
	return MG
}
