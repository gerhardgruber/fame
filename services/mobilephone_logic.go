package services

import (
	"fmt"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// MayUserSeePhone returns true if the given user may in a normal context see a mobile phone
func MayUserSeePhone(user *models.User, phone *models.MobilePhone, db *gorm.DB) (bool, *lib.FameError) {
	// Check if currently logged in user may create the phone record for the given user
	if phone.UserID == user.ID {
		return true, nil
	}

	return false, nil
}

// GetMobilePhoneOrInsert either returns the phone with the given number if it is in the database
// or creates a new one with the parameters in newPhone
// If the UserID is 0 a new User will be created as well
func GetMobilePhoneOrInsert(newPhone *models.MobilePhone, user *models.User, db *gorm.DB) (created bool, phone *models.MobilePhone, serr *lib.FameError) {
	// TODO: exchange for better validity check
	if newPhone.PhoneNumber == "" {
		return false, nil, &lib.FameError{
			ErrorCode:    "InvalidPhoneNumber",
			ErrorMessage: fmt.Sprintf("Invalid PhoneNumber '%s'", newPhone.PhoneNumber),
			Caption:      "ERR_INVALID_PHONE_NUMBER",
			CaptionData: map[string]interface{}{
				"phoneNumber": newPhone.PhoneNumber,
			},
		}
	}

	phone, err := GetMobilePhoneByNumber(newPhone.PhoneNumber, db)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return false, nil, lib.DataCorruptionError(err)
		}
	} else {
		return true, phone, nil
	}

	phone = &models.MobilePhone{
		PhoneNumber: newPhone.PhoneNumber,
	}

	// Prepare for insert
	phone.ID = 0

	if phone.Passcode != "" {
		var err error
		phone.PasscodeHash, err = models.HashPassword(phone.Passcode)
		if err != nil {
			return false, nil, &lib.FameError{
				ErrorCode:    "HashPasswordError",
				ErrorMessage: fmt.Sprintf("An error occured while hashing a password: %s", err),
				Caption:      "ERR_INTERNAL",
				CaptionData:  nil,
			}
		}
	}

	if user == nil {
		var serr *lib.FameError
		user, serr = GetOrCreateUserByName(phone.PhoneNumber, db)
		if serr != nil {
			return false, nil, serr
		}
	}

	phone.UserID = user.ID

	if err := db.Create(phone).Error; err != nil {
		return false, nil, lib.DataCorruptionError(
			fmt.Errorf("An error occured when creating the phone %+v: %s", phone, err.Error()),
		)
	}

	phone.User.MobilePhone = phone

	return true, phone, nil
}

// GetMobilePhoneByNumber loads the mobilephone with the given number and preloads its user
func GetMobilePhoneByNumber(number string, db *gorm.DB) (phone *models.MobilePhone, err error) {
	phone = &models.MobilePhone{}
	err = db.Where(db.L(models.MobilePhoneT, "PhoneNumber").Eq(number)).
		Preload("User").First(phone).Error
	return
}

// CreateMobilePhone creates the given mobilephone in the database if the PhoneNumber does not yet exist
func CreateMobilePhone(phone *models.MobilePhone, user *models.User, db *gorm.DB) *lib.FameError {
	created, newPhone, serr := GetMobilePhoneOrInsert(phone, user, db)
	if serr != nil {
		return serr
	}
	if created == false {
		return &lib.FameError{
			ErrorCode:    "PhoneExists",
			ErrorMessage: fmt.Sprintf("Phone with number '%s' already exists", phone.PhoneNumber),
			Caption:      "ERR_MOBILE_PHONE_EXISTS",
			CaptionData: map[string]interface{}{
				"phoneNumber": phone.PhoneNumber,
			},
		}
	}

	(*phone) = (*newPhone)

	return nil
}

// GetMobilePhone loads the mobilephone with the given ID
func GetMobilePhone(id uint64, db *gorm.DB) (*models.MobilePhone, *lib.FameError) {
	mobile := &models.MobilePhone{}

	if err := db.Preload("User").First(mobile, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, lib.ObjectNotFoundError(
				fmt.Errorf("Could not find mobile phone with id %d: %s", id, err.Error()),
			)
		} else {
			return nil, lib.DataCorruptionError(
				fmt.Errorf("Error finding mobile phone with id %d: %s", id, err.Error()),
			)
		}
	}

	return mobile, nil
}

// GetUserForMobilephone Loads the user with the mobile phone number
func GetUserForMobilephone(PhoneNumber string, db *gorm.DB) (*models.User, error) {
	log.Infof("GetUserForMobilephone")

	user := &models.User{}
	if err := db.Model(models.UserT).
		Joins(db.InnerJoin(models.MobilePhoneT).
			On(db.L(models.UserT, "ID"), db.L(models.MobilePhoneT, "UserID"))).
		Where(db.L(models.MobilePhoneT, "PhoneNumber").Eq(PhoneNumber)).
		Find(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func UpdatePhone(phone *models.MobilePhone, db *gorm.DB) *lib.FameError {
	err := db.Save(phone).Error
	if err != nil {
		return lib.DataCorruptionError(fmt.Errorf("Error while saving phone %+v: %s", phone, err.Error()))
	}
	return nil
}
