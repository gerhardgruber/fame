package models

import "fmt"

var (
	//AddressT TODO: comment
	AddressT = &Address{}
)

//Address TODO: comment
type Address struct {
	FameModel `gorm:"embedded_prefix:addr_"`
	Street    string
	Number    string
	Postcode  string
	City      string
	Country   string
	Longitude *float64
	Latitude  *float64
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (a *Address) ColumnPrefix() string {
	return "addr_"
}

// GetAddress returns the Address in a String in the local Format
// TODO: Local Format Check based on Country (e.g. Postcode in GB)
func (a *Address) GetAddress() string {
	return fmt.Sprintf("%s %s, %s %s, %s", a.Street, a.Number, a.Postcode, a.City, a.Country)
}
