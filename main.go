package main

import (
	"BIMSupportBot/config"
	"BIMSupportBot/internal/telegram"
	"BIMSupportBot/repository"
	"log"
	"os"
)

func main() {
	log.Println("Starting api server")

	cfg := initConfig()

	psqlDB, err := repository.NewPsqlDB(cfg)
	if err != nil {
		log.Fatalf("Postgresql init: %s", err)
	} else {
		log.Printf("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	pgRepository := repository.NewPgRepository(psqlDB)
	msgHandler := telegram.NewMessageHandler(cfg, pgRepository,echo.)
	telegram.InitBot(cfg)

	defer psqlDB.Close()
}

func initConfig() *config.Config {

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
	return "./config/config-local"
}
