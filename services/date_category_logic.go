package services

import (
	"fmt"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/jinzhu/gorm"
)

// GetDateCategories loads all date categories
func GetDateCategories(db *gorm.DB) (dateCategories *[]models.DateCategory, serr *lib.FameError) {
	dateCategories = &[]models.DateCategory{}
	if err := db.Find(dateCategories).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get date categories: %s", err),
		)
	}

	return dateCategories, nil
}

func CreateDateCategory(c *lib.Config, db *gorm.DB, u *models.User, dc *models.DateCategory) (*models.DateCategory, *lib.FameError) {
	if err := db.Create(dc).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not create date category: %s", err),
		)
	}

	return dc, nil
}

func UpdateDateCategory(c *lib.Config, db *gorm.DB, id uint64, dc *models.DateCategory) (*models.DateCategory, *lib.FameError) {
	dateCategory, ferr := GetDateCategoryByID(db, id)
	if ferr != nil {
		return nil, ferr
	}

	dateCategory.Name = dc.Name

	if err := db.Save(dateCategory).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not update date category: %s", err),
		)
	}

	return dateCategory, nil
}

// GetDateCategoryByID loads a date with the given ID
func GetDateCategoryByID(db *gorm.DB, id uint64) (dateCategory *models.DateCategory, ferr *lib.FameError) {
	dateCategory = &models.DateCategory{}
	if err := db.First(dateCategory, id).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get date category %d: %s", id, err),
		)
	}

	return dateCategory, nil
}
