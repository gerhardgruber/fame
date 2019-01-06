package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/jinzhu/gorm"
)

// GetUsers loads all users
func GetUsers(db *gorm.DB) (users *[]models.User, serr *lib.FameError) {
	users = &[]models.User{}
	if err := db.Find(users).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get users: %s", err),
		)
	}

	return users, nil
}

// GetUserByID loads a user with the given ID
func GetUserByID(id uint64, db *gorm.DB) (user *models.User, serr *lib.FameError) {
	user = &models.User{}
	if err := db.First(user, id).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get user %d: %s", id, err),
		)
	}

	return user, nil
}

// GetUserByEMail looks up a User based on their email address
func GetUserByEMail(email string, db *gorm.DB) (user *models.User, serr *lib.FameError) {
	email = strings.TrimSpace(email)

	user = &models.User{}
	err := db.Where(db.L(models.UserT, "EMail").Eq(email)).First(user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, lib.ObjectNotFoundError(
				fmt.Errorf("Could not find user with email '%s': %s", email, err),
			)
		}

		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get user by email '%s': %s", email, err),
		)
	}
	return user, nil
}

// ChangePassword will fetch the user with the given userID, checks the
// given old (=current password) sets its PW to the new password, hashes
// the password and stores the user record to the database.
func ChangePassword(userID uint64, oldPassword string, newPassword string, db *gorm.DB) *lib.FameError {
	u, err := GetUserByID(userID, db)
	if err != nil {
		return err
	}

	passwordMatches := models.ComparePassword(u.PasswordHash, oldPassword)
	if !passwordMatches {
		return &lib.FameError{
			ErrorCode:    "WrongPassword",
			ErrorMessage: "The entered password is wrong",
			Caption:      "ERR_WRONG_PASSWORD",
			CaptionData:  nil,
		}
	}

	u.PW = newPassword
	err = CheckPassword(u)
	if err != nil {
		return err
	}

	u.HashPasswordFromPW()

	return SaveUser(u, db)
}

// SaveUser writes a user to the database
func SaveUser(user *models.User, db *gorm.DB) *lib.FameError {
	if err := db.Save(user).Error; err != nil {
		return lib.DataCorruptionError(
			fmt.Errorf("Could not save user %d: %s", user.ID, err),
		)
	}
	return nil
}

// GetOrCreateUserByEMail either returns an existing user with the same email address
// or creates a new empty one
func GetOrCreateUserByEMail(email string, db *gorm.DB) (*models.User, *lib.FameError) {
	email = strings.TrimSpace(email)

	if email == "" {
		return nil, lib.InvalidParamsError(
			fmt.Errorf("Could get or create user without email"),
		)
	}

	user, serr := GetUserByEMail(email, db)
	if serr != nil {
		if !serr.IsObjectNotFoundError() {
			return nil, lib.DataCorruptionError(
				fmt.Errorf("Could not get user by email '%s' for get or create user: %s", email, serr),
			)
		}
	} else if user != nil {
		return user, nil
	}

	user = &models.User{
		EMail: email,
		Name:  email,
	}

	serr = SaveUser(user, db)
	if serr != nil {
		return nil, serr
	}

	return user, nil
}

// RegisterUser registers a new user if there is no user with the same email address
// if there is then both the previous user as well as an error is returned
func RegisterUser(fname, lname, email, lang, pw string, db *gorm.DB) (*models.User, *lib.FameError) {
	email = strings.TrimSpace(email)

	if email == "" {
		return nil, lib.InvalidParamsError(
			fmt.Errorf("Could not register user without email"),
		)
	}

	user, serr := GetUserByEMail(email, db)
	if serr != nil {
		if !serr.IsObjectNotFoundError() {
			return nil, lib.DataCorruptionError(
				fmt.Errorf("Could not get user with email '%s' for registration: %s", serr),
			)
		}
	} else if user != nil {
		return user, lib.WorkflowError(
			fmt.Errorf("Can not register over existing EMail '%s'", email),
		)
	}

	user = &models.User{
		EMail:     email,
		FirstName: fname,
		LastName:  lname,
		Lang:      lang,
		PW:        pw,
	}

	serr = CheckPassword(user)
	if serr != nil {
		return nil, serr
	}
	if err := user.HashPasswordFromPW(); err != nil {
		return nil, lib.InternalError(
			fmt.Errorf("Could not hash password for registration: %s", err),
		)
	}

	if err := db.Create(user).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not create user in registration: %s", err),
		)
	}

	return user, nil
}

func CheckPassword(user *models.User) *lib.FameError {
	if len(user.PW) < models.MinPasswordLength {
		return &lib.FameError{
			Caption: "PASSWORD_TOO_SHORT",
			CaptionData: map[string]interface{}{
				"MinLength": strconv.Itoa(models.MinPasswordLength),
			},
			ErrorCode:    "PasswordTooShort",
			ErrorMessage: "The entered password is too short",
		}
	}

	return nil
}
