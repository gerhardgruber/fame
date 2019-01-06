package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/gerhardgruber/fame/services"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// TODO: Clean up, I don't know what the app expects of this endpoint
func createMobilePhonePosition(w http.ResponseWriter, r *http.Request, c *lib.Config) {
	r.ParseForm()
	vars := mux.Vars(r)
	mobilePhoneID := vars["mobile_phone_id"]
	decoder := json.NewDecoder(r.Body)
	position := models.Position{}

	err := decoder.Decode(&position)
	if err != nil {
		log.Errorf("InvalidObjectError %+v", err)
		replyError(w, "InvalidObjectError")
		return
	}
	position.MobilePhoneID, err = strconv.ParseUint(mobilePhoneID, 10, 64)
	if err != nil {
		log.Error("Could not convert mobilePhoneID", err)
		replyError(w, "ConversionError")
		return
	}

	log.Infof("mobilePhoneID: %+v", position.MobilePhoneID)
	log.Infof("Longitude: %+v", position.Longitude)
	log.Infof("Latitude: %+v", position.Latitude)

	db, serr := c.GetDatabaseConnection()
	if serr != nil {
		log.Error("Could not get database for login", serr)
		replyError(w, "DatabaseError")
		return
	}

	mobilePhone := models.MobilePhone{}
	db.First(&mobilePhone, position.MobilePhoneID)

	if db.Create(&position).Error != nil {
		log.Error("Could not create MobilePhonePosition", err)
		replyError(w, "MobilePhonePositionCreateError")
		return
	}

	replyData(w, map[string]interface{}{})
}

func positionTest(w http.ResponseWriter, r *http.Request, c *lib.Config) {
	r.ParseForm()
	ID1 := r.Form.Get("id1")
	ID2 := r.Form.Get("id2")
	Adress := r.Form.Get("Adress")

	log.Infof("Adress: %s", Adress)

	db, err := c.GetDatabaseConnection()
	if err != nil {
		log.Error("Could not get database for login", err)
		replyError(w, "DatabaseError")
		return
	}
	Pos1 := models.Position{}
	Pos2 := models.Position{}
	db.Find(&Pos1, ID1)
	db.Find(&Pos2, ID2)
	distance := services.CalculateDistance(&Pos1, &Pos2)
	log.Infof("Distance between Points: %+v", distance)

	pos, _ := services.GetPositionForAddress(Adress, c)

	log.Infof("Address %+v in Points: %+v %+v", Adress, pos.Longitude, pos.Latitude)

	replyData(w, map[string]interface{}{"DistanceBetweenPoints": distance})
}

// RegisterPositionControllerRoutes Registers the functions
func RegisterPositionControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/mobile_phones/{mobile_phone_id:[0-9]+}/positions", serviceWrapper("CreateMobilePhonePosition", createMobilePhonePosition, config)).Methods("POST")
	router.HandleFunc("/positiontest", serviceWrapper("positionTest", positionTest, config)).Methods("POST")

	router.HandleFunc("/app/v1/mobile_phones/{mobile_phone_id:[0-9]+}/positions", serviceWrapper("CreateMobilePhonePosition", createMobilePhonePosition, config)).Methods("POST")
}
