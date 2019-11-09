package models

import (
	"crypto/sha512"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// UserType - sets the user type (simple user-rights)
type UserRightType int

const (
	StandardUser UserRightType = 0
	AdminUser    UserRightType = 1

	MinPasswordLength = 6
)

var (
	// UserT TODO: comment
	UserT = &User{}
)

// User TODO: comment
type User struct {
	FameModel    `gorm:"embedded_prefix:usr_"`
	Name         string
	FirstName    string
	LastName     string
	PW           string `gorm:"-"`
	PasswordHash string
	Lang         string
	EMail        string
	MobilePhone  *MobilePhone
	RightType    UserRightType
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (u *User) ColumnPrefix() string {
	return "usr_"
}

const (
	passwordPepper = "gGGvM-cdCb4t_SonhqRyFzOOJk4Irz17QsX1CrWaOXM4TvsX8XGr7IGPHFW4I0zI"
)

// ComparePassword compares Password from DB with login PWD
func ComparePassword(hashedPassword string, plaintextPassword string) bool {
	hashedSecret := sha512.Sum512([]byte(passwordPepper + plaintextPassword))
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), hashedSecret[:]) == nil
}

// HashPassword hashes the password before storing in db
func HashPassword(password string) (string, error) {
	hashedSecret := sha512.Sum512([]byte(passwordPepper + password))

	hashedPassword, err := bcrypt.GenerateFromPassword(hashedSecret[:], bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// HashPasswordFromPW hashes the PW Field to Password
func (u *User) HashPasswordFromPW() (err error) {
	u.PasswordHash, err = HashPassword(u.PW)
	return
}

// CheckMandatoryFieldsSet checks if all mandatory Fields are set
func (u *User) CheckMandatoryFieldsSet() bool {
	return u.Name != "" && u.FirstName != "" && u.LastName != "" && u.PW != "" && u.Lang != "" && u.EMail != ""
}

func (u *User) FullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	} else if u.FirstName != "" {
		return u.FirstName
	} else {
		return u.LastName
	}
}

func (u *User) IsFullyRegistered() bool {
	return u.PasswordHash != ""
}
