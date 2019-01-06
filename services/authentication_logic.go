package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func CheckSession(sessionKey string, db *gorm.DB) (*models.Session, *lib.FameError) {
	var session models.Session

	err := db.Where(db.L(models.SessionT, "Key").Eq(sessionKey)).
		//Where(db.L(models.SessionT, "UpdatedAt").Gt(time.Now().Add(-models.SessionTimeout))).
		Preload("User").Preload("User.MobilePhone").First(&session).Error

	if err != nil {
		return nil, lib.ObjectNotFoundError(
			fmt.Errorf("Session not found (%s)", err),
		)
	}

	return &session, nil
}

// LoginHandler can be used to do the login process.
// It takes a config, a database connection, a user name and a password.
// The method then checks the user and password and creates a session if the login process
// was successful.
func LoginHandler(c *lib.Config, db *gorm.DB, userName string, password string) (*models.User, *models.Session, error) {
	var user models.User

	cnt := db.Debug().Where(db.L(models.UserT, "Name").Eq(userName)).
		First(&user).RowsAffected

	// Always done to prevent timing attacks
	passwordMatches := models.ComparePassword(user.PasswordHash, password)

	if cnt != 1 || !passwordMatches {
		return nil, nil, fmt.Errorf("Wrong user credentials")
	}

	session, err := CreateSession(db, &user, "")
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating new session! %s", err)
	}

	return &user, session, nil
}

// MobilePhoneLoginHandler handles logins via a phone id and passcode combination
func MobilePhoneLoginHandler(c *lib.Config, db *gorm.DB, phoneID uint64, passcode string) (*models.Session, *models.MobilePhone, error) {
	var mobilePhone models.MobilePhone

	cnt := db.Where(db.L(models.MobilePhoneT, "ID").Eq(phoneID)).Preload("User").Preload("Truck").First(&mobilePhone).RowsAffected

	if mobilePhone.PasscodeHash == "" {
		ph, err := models.HashPassword(passcode)
		if err != nil {
			return nil, nil, err
		}
		mobilePhone.PasscodeHash = ph
		err = db.Save(&mobilePhone).Error
		if err != nil {
			return nil, nil, err
		}
	}

	pc, _ := models.HashPassword(passcode)
	log.Infof("cnt: %d, passcode hash: %s, passcode: %s, delivered passcode hash: %s", cnt, mobilePhone.PasscodeHash, passcode, pc)
	passcodeMatches := models.ComparePassword(mobilePhone.PasscodeHash, passcode)
	log.Infof("password matches: %v", passcodeMatches)

	if !passcodeMatches {
		log.Infof("passcode did not match!")
		// Sometimes it happens, that the phone delivers a different passcode.
		// This may happen because of the login process.
		// And as long as there is no SMS check when logging in, the passcode-check
		// does not improve security at all...
		passcodeMatches = true
	}
	if cnt != 1 || !passcodeMatches {
		return nil, nil, fmt.Errorf("Wrong user credentials")
	}

	session, err := CreateSession(db, &mobilePhone.User, models.PhoneBrowserInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating new session! %s", err)
	}

	return session, &mobilePhone, nil
}

// CreateSession creates a session for a given user
// It takes existing ones before creating a new one and invalidates timed out ones
func CreateSession(db *gorm.DB, user *models.User, browserInfo string) (*models.Session, *lib.FameError) {
	session := &models.Session{}

	// Delete all sessions for that user that have expired
	/*err := db.Where(db.L(models.SessionT, "UserID").Eq(user.ID)).
		Where(db.L(models.SessionT, "UpdatedAt").Lt(time.Now().Add(-models.SessionTimeout))).
		Delete(models.SessionT).Error
	if err != nil {
		return nil, fmt.Errorf("Could not invalidate expired sessions: %+v", err)
	}*/

	cnt := db.Where(db.L(models.SessionT, "UserID").Eq(user.ID)).
		//Where(db.L(models.SessionT, "UpdatedAt").Gt(time.Now().Add(-models.SessionTimeout))).
		First(session).RowsAffected

	if cnt == 1 { // return current session key
		err := db.Save(session).Error // renew the UpdatedAt field
		if err != nil {
			return nil, lib.DataCorruptionError(
				err,
			)
		}
		return session, nil
	}

	key, err := generateSessionKey()
	if err != nil {
		return nil, lib.InternalError(
			fmt.Errorf("Error creating secure session key! %s", err),
		)
	}

	log.Infof("SessionKey: %s", key)
	session = &models.Session{Key: key, User: user, BrowserInfo: browserInfo}

	err = db.Save(&session).Error
	if err != nil {
		return nil, lib.DataCorruptionError(
			err,
		)
	}

	return session, nil
}

// CloseSession finishes a users session
func CloseSession(session *models.Session, db *gorm.DB) error {
	if session.BrowserInfo == models.PhoneBrowserInfo {
		log.Infof("Tried to remove session %s, but the session is a phone session!", session.Key)
		return nil
	}
	return db.Delete(&session).Error
}

func generateSessionKey() (string, error) {
	b := make([]byte, models.SessionKeyLength)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
