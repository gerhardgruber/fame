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
