package db

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ivanstepachev/tg_refferal/store/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/yaml.v3"
)

type Config struct {
    Database struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		DBName string `yaml:"dbname"`
        Username string `yaml:"user"`
        Password string `yaml:"pass"`
    } `yaml:"database"`
}

// Connection to DB and make automigration by models
func Init() *gorm.DB {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Println(err.Error())
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Println(err.Error())
	}

	dbURL := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.DBName, cfg.Database.Password)
	
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln(err)
	}
	// Using for make migrations to DB before server starting
	// Use flag -migrate=true
	cmdArgs := flag.Bool("migrate", false, "Make migration to DB before server starting")
	flag.Parse()
	if *cmdArgs {
		db.AutoMigrate(&models.User{}, &models.Transaction{})
	}
	return db
}