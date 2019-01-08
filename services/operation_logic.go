package services

import (
	"fmt"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/jinzhu/gorm"
)

// GetOperations loads all operations
func GetOperations(db *gorm.DB) (operations *[]models.Operation, serr *lib.FameError) {
	operations = &[]models.Operation{}
	if err := db.Find(operations).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get operations: %s", err),
		)
	}

	return operations, nil
}

func CreateOperation(title, fname, lname string, db *gorm.DB) (*models.Operation, *lib.FameError) {
	operation := &models.Operation{
		Title:     title,
		FirstName: fname,
		LastName:  lname,
	}

	if err := db.Create(operation).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not create operation: %s", err),
		)
	}

	return operation, nil
}

// GetOperationByID loads an operation with the given ID
func GetOperationByID(id uint64, db *gorm.DB) (operation *models.Operation, serr *lib.FameError) {
	operation = &models.Operation{}
	if err := db.First(operation, id).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get operation %d: %s", id, err),
		)
	}

	return operation, nil
}
