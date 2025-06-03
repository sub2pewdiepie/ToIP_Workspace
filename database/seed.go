package database

import (
	"log"
	"space/models"
	"space/repositories"
	"space/utils"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

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

func SeedSubjects(db *gorm.DB, subjectRepo *repositories.SubjectRepository, academicGroupRepo *repositories.AcademicGroupRepository) error {
	// Fetch existing academic groups
	var academicGroups []models.AcademicGroup
	if err := db.Find(&academicGroups).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to fetch academic groups for seeding subjects")
		return err
	}

	if len(academicGroups) == 0 {
		utils.Logger.Warn("No academic groups found for seeding subjects, setting AcademicGroupID to NULL")
	}

	subjects := []models.Subject{
		{
			Name: "Mathematics",

			AcademicGroupID: getAcademicGroupID(academicGroups, 0),
		},
		{
			Name: "Physics",

			AcademicGroupID: getAcademicGroupID(academicGroups, 1),
		},
		{
			Name: "Computer Science",

			AcademicGroupID: getAcademicGroupID(academicGroups, 0),
		},
		{
			Name: "English Literature",

			AcademicGroupID: getAcademicGroupID(academicGroups, 1),
		},
		{
			Name: "History",

			AcademicGroupID: getAcademicGroupID(academicGroups, 0),
		},
	}

	for _, subject := range subjects {
		existing, err := subjectRepo.FindByName(subject.Name)
		if err != nil && err != gorm.ErrRecordNotFound {
			utils.Logger.WithFields(logrus.Fields{
				"error": err,
				"name":  subject.Name,
			}).Error("Failed to check existing subject")
			return err
		}
		if existing == nil {
			if err := subjectRepo.Create(&subject); err != nil {
				utils.Logger.WithFields(logrus.Fields{
					"error": err,
					"name":  subject.Name,
				}).Error("Failed to seed subject")
				return err
			}
			utils.Logger.WithFields(logrus.Fields{
				"name":              subject.Name,
				"academic_group_id": subject.AcademicGroupID,
			}).Info("Seeded subject")
		} else {
			utils.Logger.WithFields(logrus.Fields{
				"name": subject.Name,
			}).Debug("Subject already exists, skipping")
		}
	}
	return nil
}

// Helper function to get AcademicGroupID or 0 (NULL) if index is out of bounds
func getAcademicGroupID(groups []models.AcademicGroup, index int) int32 {
	if index < len(groups) {
		return groups[index].AcademicGroupID
	}
	return 0 // Maps to NULL in database
}
