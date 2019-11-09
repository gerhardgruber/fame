package models

// MobileDeviceType - defines the devicetype
type MobileDeviceType int

const (
	// DeviceiOS - iOS
	DeviceiOS MobileDeviceType = 1

	// DeviceAndroid - Android
	DeviceAndroid MobileDeviceType = 2

	// DeviceWindows - Windows Phone
	DeviceWindows MobileDeviceType = 3

	// DeviceOther - any other device
	DeviceOther MobileDeviceType = 99
)

var (
	// MobilePhoneT TODO: comment
	MobilePhoneT = &MobilePhone{}
)

// MobilePhone TODO: comment
type MobilePhone struct {
	FameModel   `gorm:"embedded_prefix:mp_"`
	PhoneNumber string
	Device      string
	DeviceType  MobileDeviceType

	// Identifying private phone bound "password" that is used for phone based authentication
	Passcode     string `gorm:"-"`
	PasscodeHash string

	UserID uint64
	User   *User
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (m *MobilePhone) ColumnPrefix() string {
	return "mp_"
}

func (m *MobilePhone) HashPasscode() {
	m.PasscodeHash, _ = HashPassword(m.Passcode)
}
