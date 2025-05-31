package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"space/models"
	"space/repositories"
	"space/utils"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Printf("Database ensured. %s", err)
		return Config{}, fmt.Errorf("failed to decode config file: %w", err)
	}
	return config, nil
}

func ensureDatabaseExists(config Config) error {
	// Connect to PostgreSQL server without specifying a database
	serverDSN := fmt.Sprintf(
		"host=%s user=%s password=%s port=%d sslmode=%s",
		config.Host, config.User, config.Password, config.Port, config.SSLMode,
	)

	db, err := sql.Open("postgres", serverDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL server: %w", err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)"
	if err := db.QueryRow(query, config.DBName).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check if database %s exists: %w", config.DBName, err)
	}

	// Create database if it doesn't exist
	if !exists {
		query := fmt.Sprintf("CREATE DATABASE %s", quoteIdentifier(config.DBName))
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create database %s: %w", config.DBName, err)
		}
		log.Printf("Database %s created successfully.", config.DBName)
	} else {
		log.Printf("Database %s already exists.", config.DBName)
	}
	return nil
}

// quoteIdentifier escapes a PostgreSQL identifier (e.g., database name) to prevent SQL injection
func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func ConnectDatabase() error {
	config, err := LoadConfig("./config/config.json")
	if err != nil {
		utils.Logger.WithField("error", err).Error("Failed to load .env file")
		return fmt.Errorf("error loading config: %w", err)
	}
	// Ensure the database exists
	// fmt.Printf("Database ensured.")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable SQL logging for debugging
	})
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"database": config.DBName,
			"host":     config.Host, // additional useful context
		}).Error("Failed to connect to database")
		return fmt.Errorf("failed to connect to database %s: %w", config.DBName, err)
	}
	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("failed to connect to database: %w", err)
	// }

	// if err := ensureDatabaseExists(config); err != nil {
	// 	return err
	// }
	// Автоматическая миграция всех таблиц
	err = DB.AutoMigrate(
		&models.User{},          // No dependencies
		&models.AcademicGroup{}, // No dependencies
		&models.Group{},         // Depends on User, AcademicGroup
		&models.GroupModer{},    // Depends on Group, User
		&models.GroupUser{},     // Depends on Group, User
		&models.Subject{},       // Depends on Group
		&models.Task{},          // Depends on Subject, User
		&models.Material{},      // Depends on Subject, User
		&models.TimeSlot{},      // No dependencies
		&models.Schedule{},      // Depends on Group, Subject, TimeSlot
		&models.GroupApplication{},
	)
	if err != nil {
		utils.Logger.
			WithError(err).
			WithField("action", "database_migration").
			Error("Database migration failed")
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	utils.Logger.WithFields(logrus.Fields{
		"event":  "database_startup",
		"status": "success",
	}).Info("Database connected and migrated successfully")
	log.Println("Database connected and migrated successfully.")
	return nil
}

func SeedAcademicGroups(db *gorm.DB, repo *repositories.AcademicGroupRepository) error {
	// Flush academic_groups table
	if err := DB.Exec("DELETE FROM academic_groups").Error; err != nil {
		log.Fatalf("Failed to flush database: %v", err)
	}
	utils.Logger.WithFields(logrus.Fields{
		"action": "database_flush",
		"status": "success",
	}).Info("Database flushed successfully")
	log.Println("Database flushed successfully.")

	groups := []models.AcademicGroup{
		{AcademicGroupID: 1, Name: "ЭФМО-01-24", CreatedAt: time.Now()},
		{AcademicGroupID: 2, Name: "ИКБО-14-20", CreatedAt: time.Now()},
		{AcademicGroupID: 3, Name: "ИКБО-15-20", CreatedAt: time.Now()},
	}

	for _, ac_group := range groups {
		if err := DB.Create(&ac_group).Error; err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"error":      err,
				"group_id":   ac_group.AcademicGroupID,
				"group_name": ac_group.Name,
				"operation":  "seed_academic_group",
			}).Error("Failed to seed academic group")
			log.Printf("Failed to seed academic group: %v", err)
		}
	}
	// Successful seeding
	utils.Logger.WithFields(logrus.Fields{
		"event":  "data_seeding",
		"entity": "academic_groups",
		"count":  len(groups),
		// "duration_ms": time.Since(startTime).Milliseconds(),
	}).Info("Academic groups seeded successfully")
	log.Println("Academic groups seeded.")
	return nil
}
