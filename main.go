package main

import (
	"BIMSupportBot/config"
	"BIMSupportBot/internal/telegram"
	"BIMSupportBot/repository"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	log.Println("Starting server")

	cfg := initConfig()

	// Set client options
	clientOptions := options.Client().ApplyURI(cfg.Mongo.Url)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Mongo init: %s", err)
	}
	fmt.Println("Connected to MongoDB!")

	collection := client.Database(cfg.Mongo.DataBase).Collection(cfg.Mongo.Collection)
	mongoRepository := repository.NewMongoRepository(collection)

	msgHandler := telegram.NewMessageHandler(cfg, mongoRepository, context.TODO())
	telegram.InitBot(cfg, msgHandler)
}

func initConfig() *config.Config {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	configPath := GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}
	return cfg
}

func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local.yml"
}
