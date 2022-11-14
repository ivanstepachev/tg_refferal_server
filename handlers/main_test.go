package handlers

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ivanstepachev/tg_refferal/store/db"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/yaml.v3"
)

var testH handlers

func TestMain(m *testing.M) {
	f, err := os.Open("../config.yml")
	if err != nil {
		log.Println(err.Error())
	}
	defer f.Close()

	var cfg db.Config
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
	testH = New(db)
	os.Exit(m.Run())
}
