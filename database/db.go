package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	Port     int    `json:"port"`
	SSLMode  string `json:"sslmode"`
}

func LoadConfig(filepath string) (Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func ConnectDatabase() {
	config, _err := LoadConfig("config/config.json")
	if _err != nil {
		fmt.Printf("Error loading config: %v\n", _err)
		return
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	// err = DB.AutoMigrate(&model.BikePart{})
	// if err != nil {
	// 	log.Fatal("Failed to migrate database:", err)
	// }

	log.Println("Database connected successfully.")
}

// SeedDatabase flushes the database and seeds it with fresh data
func SeedDatabase() {

}
