package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/jinzhu/gorm"
)

// GetUsers loads all users
func GetUsers(db *gorm.DB) (users *[]models.User, serr *lib.FameError) {
	users = &[]models.User{}
	if err := db.Model(models.UserT).Order(db.L(models.UserT, "LastName").OrderAsc()).Find(users).Error; err != nil {
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

// GetUserByName looks up a User based on their name
func GetUserByName(name string, db *gorm.DB) (user *models.User, serr *lib.FameError) {
	name = strings.TrimSpace(name)

	user = &models.User{}
	err := db.Where(db.L(models.UserT, "Name").Eq(name)).First(user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, lib.ObjectNotFoundError(
				fmt.Errorf("Could not find user with name '%s': %s", name, err),
			)
		}

		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get user by name '%s': %s", name, err),
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

// GetOrCreateUserByName either returns an existing user with the same name
// or creates a new empty one
func GetOrCreateUserByName(name string, db *gorm.DB) (*models.User, *lib.FameError) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, lib.InvalidParamsError(
			fmt.Errorf("Could get or create user without name"),
		)
	}

	user, serr := GetUserByName(name, db)
	if serr != nil {
		if !serr.IsObjectNotFoundError() {
			return nil, lib.DataCorruptionError(
				fmt.Errorf("Could not get user by name '%s' for get or create user: %s", name, serr),
			)
		}
	} else if user != nil {
		return user, nil
	}

	user = &models.User{
		Name: name,
	}

	serr = SaveUser(user, db)
	if serr != nil {
		return nil, serr
	}

	return user, nil
}

// RegisterUser registers a new user if there is no user with the same email address
// if there is then both the previous user as well as an error is returned
func RegisterUser(name, fname, lname, email, lang, pw string, rightType models.UserRightType, db *gorm.DB) (*models.User, *lib.FameError) {
	email = strings.TrimSpace(email)

	if email == "" {
		return nil, lib.InvalidParamsError(
			fmt.Errorf("Could not register user without email"),
		)
	}

	user, serr := GetUserByName(name, db)
	if serr != nil {
		if !serr.IsObjectNotFoundError() {
			return nil, lib.DataCorruptionError(
				fmt.Errorf("Could not get user with name '%s' for registration: %s", name, serr),
			)
		}
	} else if user != nil {
		return user, lib.WorkflowError(
			fmt.Errorf("Can not register over existing name '%s'", name),
		)
	}

	user = &models.User{
		Name:      name,
		EMail:     email,
		FirstName: fname,
		LastName:  lname,
		Lang:      lang,
		PW:        pw,
		RightType: rightType,
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

func GetCurrentPausesByUserID(userID uint64, db *gorm.DB) (*models.PauseAction, *models.PauseAction, *lib.FameError) {
	operationPause := &models.PauseAction{}
	trainingPause := &models.PauseAction{}

	err := db.Model(operationPause).Where(
		db.L(operationPause, "UserID").Eq(userID),
	).Where(
		db.L(operationPause, "Type").Eq(models.OperationPause),
	).Last(operationPause).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			operationPause = nil
		} else {
			return nil, nil, lib.DataCorruptionError(fmt.Errorf("Error selecting pause actions for user %d! %w", userID, err))
		}
	}

	err = db.Model(trainingPause).Where(
		db.L(trainingPause, "UserID").Eq(userID),
	).Where(
		db.L(trainingPause, "Type").Eq(models.TrainingPause),
	).Last(trainingPause).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			trainingPause = nil
		} else {
			return nil, nil, lib.DataCorruptionError(fmt.Errorf("Error selecting pause actions for user %d! %w", userID, err))
		}
	}

	return trainingPause, operationPause, nil
}

func StartPause(userID uint64, pauseType models.PauseType, date time.Time, db *gorm.DB) (*models.PauseAction, *lib.FameError) {
	tp, op, ferr := GetCurrentPausesByUserID(userID, db)
	if ferr != nil {
		return nil, ferr
	}

	if pauseType == models.TrainingPause && tp != nil && tp.EndTime == nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Error starting training pause for user %d! Training pause is already active", userID),
		)
	} else if pauseType == models.OperationPause && op != nil && op.EndTime == nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Error starting operation pause for user %d! Operation pause is already active", userID),
		)
	}

	pause := &models.PauseAction{
		UserID:    userID,
		Type:      pauseType,
		StartTime: &date,
	}
	err := db.Save(pause).Error
	if err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Error saving pause for user %d! %w", userID, err),
		)
	}

	return pause, nil
}

func StopPause(userID uint64, pauseType models.PauseType, date time.Time, db *gorm.DB) (*models.PauseAction, *lib.FameError) {
	tp, op, ferr := GetCurrentPausesByUserID(userID, db)
	if ferr != nil {
		return nil, ferr
	}

	if pauseType == models.TrainingPause && tp == nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Error stopping training pause for user %d! There is no training pause active", userID),
		)
	} else if pauseType == models.OperationPause && op != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Error stopping operation pause for user %d! There is no operation pause active", userID),
		)
	}

	var pause *models.PauseAction
	if pauseType == models.TrainingPause {
		pause = tp
	} else {
		pause = op
	}

	pause.EndTime = &date
	err := db.Save(pause).Error
	if err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Error saving pause for user %d! %w", userID, err),
		)
	}

	return pause, nil
}
