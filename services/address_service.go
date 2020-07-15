package services

import (
	"encoding/json"
	"fmt"
	"net/url"
	"runtime/debug"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/jinzhu/gorm"
	geo "github.com/kellydunn/golang-geo"
	log "github.com/sirupsen/logrus"
)

const (
	googleGeocodeOK = "OK"
)

var geoCoder *geo.GoogleGeocoder

// GetAddressByID looks for the address with the given ID
func GetAddressByID(id uint64, db *gorm.DB) (*models.Address, error) {
	address := &models.Address{}

	err := db.Where(db.L(models.AddressT, "ID").Eq(id)).First(address).Error
	if err != nil {
		return nil, err
	}

	return address, err
}

// LookupAddress does a GeoLookup for an address
func LookupAddress(address string, c *lib.Config) (*models.Address, *lib.FameError) {
	geoCode, err := Geocode(address, c)
	if err != nil {
		return nil, &lib.FameError{
			ErrorCode: "GeocodeError",
			Caption:   "ERR_GEOCODE",
			CaptionData: map[string]interface{}{
				"address": address,
			},
			ErrorMessage: fmt.Sprintf("Error geocoding address %s: %s", address, err.Error()),
			StackTrace:   debug.Stack(),
		}
	}

	newA, err := geoCode.getAddress(0)
	if err != nil {
		return nil, &lib.FameError{
			ErrorCode: "GeocodeError",
			Caption:   "ERR_GEOCODE",
			CaptionData: map[string]interface{}{
				"address": address,
			},
			ErrorMessage: fmt.Sprintf("Error fetching address %s from geocode %+v: %s", address, geoCode, err.Error()),
			StackTrace:   debug.Stack(),
		}
	}

	return newA, nil
}

// GetOrCreateAddress checks if an address is already in the database
func GetOrCreateAddress(address string, db *gorm.DB, c *lib.Config) (*models.Address, *lib.FameError) {
	geoCode, err := Geocode(address, c)
	if err != nil {
		log.Errorf("Error looking for address %s! %s", address, err)
	}

	var newA *models.Address
	if geoCode != nil {
		newA, err := geoCode.getAddress(0)
		if err != nil {
			return nil, &lib.FameError{
				ErrorCode: "GeocodeError",
				Caption:   "ERR_GEOCODE",
				CaptionData: map[string]interface{}{
					"address": address,
				},
				ErrorMessage: fmt.Sprintf("Error fetching address %s from geocode %+v: %s", address, geoCode, err.Error()),
				StackTrace:   debug.Stack(),
			}
		}

		oldA := &models.Address{}
		err = db.Where(db.L(models.AddressT, "Street").Eq(newA.Street).
			And(db.L(models.AddressT, "Number").Eq(newA.Number)).
			And(db.L(models.AddressT, "Postcode").Eq(newA.Postcode)).
			And(db.L(models.AddressT, "City").Eq(newA.City)).
			And(db.L(models.AddressT, "Country").Eq(newA.Country))).
			First(oldA).Error
		if err == nil {
			return oldA, nil
		}

		if !gorm.IsRecordNotFoundError(err) {
			return nil, lib.ObjectNotFoundError(
				fmt.Errorf("Could not find matching address: %s", err),
			)
		}
	} else {
		oldA := &models.Address{}
		err = db.Where(db.L(models.AddressT, "Street").Eq(address)).
			First(oldA).Error
		if err == nil {
			return oldA, nil
		}

		if !gorm.IsRecordNotFoundError(err) {
			return nil, lib.ObjectNotFoundError(
				fmt.Errorf("Could not find matching address: %s", err),
			)
		}

		newA = &models.Address{
			Street: address,
		}
	}

	err = db.Create(newA).Error

	if err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not create address %+v: %s", newA, err.Error()),
		)
	}

	return newA, nil
}

// Geocode loads a geocoded object by address from google
func Geocode(address string, c *lib.Config) (*googleGeocodeResponse, error) {
	if geoCoder == nil {
		geo.SetGoogleAPIKey(c.GetGoogleAPIKey())
		geoCoder = &geo.GoogleGeocoder{}
	}

	fmt.Printf("request: %s\n", fmt.Sprintf("address=%s&key=%s", url.QueryEscape(address), c.GoogleAPIKey))
	resData, err := geoCoder.Request(fmt.Sprintf("address=%s&key=%s", url.QueryEscape(address), c.GoogleAPIKey))
	if err != nil {
		log.Errorf("Error while getting Point: %+v", err)
		return nil, err
	}

	res := &googleGeocodeResponse{}
	err = json.Unmarshal(resData, res)
	if err != nil {
		return nil, err
	}

	if res.Status != googleGeocodeOK {
		return nil, fmt.Errorf("Non OK Geocode Status: %s", res.Status)
	}

	return res, nil
}

// GetPositionForAddress returns a Position Object for a given address
func GetPositionForAddress(address string, c *lib.Config) (*models.Position, error) {
	g, err := Geocode(address, c)
	if err != nil {
		return nil, err
	}

	if len(g.Results) < 1 {
		return nil, fmt.Errorf("No result for geocoding '%s'", address)
	}

	position := models.Position{
		Latitude:  g.Results[0].Geometry.Location.Lat,
		Longitude: g.Results[0].Geometry.Location.Lng,
	}
	return &position, nil
}

// googleGeocodeResponse is an internal struct that is used to parse the response of the google
// geocode api
type googleGeocodeResponse struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64
				Lng float64
			}
		} `json:"geometry"`
	} `json:"results"`
	Status string
}

const (
	// Probably not 100% correct handling
	ggStreet  = "route"
	ggNumber  = "street_number"
	ggZip     = "postal_code"
	ggCity    = "locality"
	ggCountry = "country"
)

func (g *googleGeocodeResponse) getAddress(index int) (*models.Address, error) {
	if index >= len(g.Results) {
		return nil, fmt.Errorf("Could not access geocode response index %d of %d", index, len(g.Results)-1)
	}

	result := g.Results[index]

	addr := &models.Address{
		Latitude:  &result.Geometry.Location.Lat,
		Longitude: &result.Geometry.Location.Lng,
	}

	components := result.AddressComponents
	for _, c := range components {
		for _, t := range c.Types {
			if t == ggStreet {
				addr.Street = c.LongName
				break
			} else if t == ggNumber {
				addr.Number = c.LongName
				break
			} else if t == ggZip {
				addr.Postcode = c.LongName
				break
			} else if t == ggCity {
				addr.City = c.LongName
				break
			} else if t == ggCountry {
				addr.Country = c.LongName
				break
			}
		}
	}

	return addr, nil
}
