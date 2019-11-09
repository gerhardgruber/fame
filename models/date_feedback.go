package models

type DateFeedbackType int

const (
	DateFeedbackTypeYes     DateFeedbackType = 1
	DateFeedbackTypeNo      DateFeedbackType = 2
	DateFeedbackTypeUnknown DateFeedbackType = 3
)

var DateFeedbackTypes = map[string]DateFeedbackType{
	"Yes":     DateFeedbackTypeYes,
	"No":      DateFeedbackTypeNo,
	"Unknown": DateFeedbackTypeUnknown,
}

var DateFeedbackT = &DateFeedback{}

type DateFeedback struct {
	FameModel `gorm:"embedded_prefix:dtf_"`

	DateID uint64
	Date   *Date
	UserID uint64
	User   *User

	Feedback DateFeedbackType
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (dtf *DateFeedback) ColumnPrefix() string {
	return "dtf_"
}
