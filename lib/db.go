package lib

import (
	"github.com/gerhardgruber/fame/models"
	"github.com/jinzhu/gorm"
)

// MigrateDatabase TODO
func MigrateDatabase(db *gorm.DB) {
	db.AutoMigrate(models.AddressT)
	db.AutoMigrate(models.MobilePhoneT)
	db.AutoMigrate(models.MobilePhoneLogT)
	db.AutoMigrate(models.PositionT)
	db.AutoMigrate(models.SessionT)
	db.AutoMigrate(models.UserT)
}
