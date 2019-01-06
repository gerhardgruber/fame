package models

import "time"

var (
	// SessionT TODO: comment
	SessionT = &Session{}
)

const (
	// SessionTimeout TODO: comment
	SessionTimeout = time.Minute * 180
	// SessionKeyLength TODO: comment
	SessionKeyLength = 48

	// PhoneBrowserInfo is a constant which is saved to BrowserInfo when
	// a session is created by a phone
	PhoneBrowserInfo = "PHONE"
)

// Session TODO: comment
type Session struct {
	FameModel   `gorm:"embedded_prefix:sess_"`
	Key         string
	UserID      *uint64
	User        *User
	BrowserInfo string
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (u *Session) ColumnPrefix() string {
	return "sess_"
}
