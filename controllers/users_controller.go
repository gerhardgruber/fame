package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/gerhardgruber/fame/services"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var (
	userInviteOk               = "ok"
	userInviteRegistered       = "registered"
	userInviteHasCompany       = "company"
	userInviteAlreadyInCompany = "already"
)

func getUsers(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	users, serr := services.GetUsers(db)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"Users": users,
	})
}

func createUser(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	type createUserParams struct {
		Name      string
		FirstName string
		LastName  string
		EMail     string
		PW        string
		RightType models.UserRightType
	}

	decoder := json.NewDecoder(r.Body)
	cup := createUserParams{}
	err := decoder.Decode(&cup)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding User: %s", err),
		))
	}

	u, serr := services.RegisterUser(
		cup.Name,
		cup.FirstName,
		cup.LastName,
		cup.EMail,
		"de",
		cup.PW,
		cup.RightType,
		db,
	)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"User": u,
	})
}

func getUser(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	userID, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("User ID '%s' could not be parsed: %s", params["id"], err),
		))
	}

	user, serr := services.GetUserByID(userID, db)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"User": user,
	})
}

func deleteUser(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("User ID could not be parsed: ", err),
		))
	}

	if id != *sess.UserID && sess.User.RightType != models.AdminUser {
		return Error(*lib.PrivilegeError(
			fmt.Errorf("User ID does not match session or session user is not admin! (%d != %d)", id, *sess.UserID),
		))
	}

	decoder := json.NewDecoder(r.Body)
	data := models.User{}
	err = decoder.Decode(&data)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("InvalidObjectError %+v", err),
		))
	}

	user, serr := services.GetUserByID(id, db)
	if serr != nil {
		return Error(*serr)
	}

	err = db.Delete(user).Error
	if err != nil {
		return Error(*lib.DataCorruptionError(
			fmt.Errorf("DataBaseError %+v", err),
		))
	}

	return Success()
}

func updateUserAPI(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("User ID could not be parsed: ", err),
		))
	}

	if id != *sess.UserID && sess.User.RightType != models.AdminUser {
		return Error(*lib.PrivilegeError(
			fmt.Errorf("User ID does not match session! (%d != %d)", id, *sess.UserID),
		))
	}

	decoder := json.NewDecoder(r.Body)
	data := models.User{}
	err = decoder.Decode(&data)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("InvalidObjectError %+v", err),
		))
	}

	user, serr := services.GetUserByID(id, db)
	if serr != nil {
		return Error(*serr)
	}

	user.FirstName = data.FirstName
	user.LastName = data.LastName
	user.EMail = data.EMail
	user.Lang = data.Lang
	user.RightType = data.RightType

	err = db.Save(user).Error
	if err != nil {
		return Error(*lib.DataCorruptionError(
			fmt.Errorf("DataBaseError %+v", err),
		))
	}

	return Success(map[string]interface{}{
		"user": user,
	})
}

func updateUser(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("User ID '%s' could not be parsed: %s", params["id"], err),
		))
	}

	// TODO: Let future admin-users edit users from the same company
	if id != *sess.UserID {
		return Error(*lib.PrivilegeError(
			fmt.Errorf("User ID %d does not match session User ID %d", id, sess.UserID),
		))
	}

	type updateUserParams struct {
		User models.User
	}

	decoder := json.NewDecoder(r.Body)
	uup := updateUserParams{}
	err = decoder.Decode(&uup)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding User: %s", err),
		))
	}

	u, serr := services.GetUserByID(id, db)
	if serr != nil {
		return Error(*serr)
	}

	u.FirstName = uup.User.FirstName
	u.LastName = uup.User.LastName
	u.EMail = uup.User.EMail

	serr = services.SaveUser(u, db)
	if serr != nil {
		return Error(*serr)
	}

	return Success()
}

func changePassword(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("User ID '%s' could not be parsed: %s", params["id"], err),
		))
	}

	// TODO: Let future admin-users edit users from the same company
	if id != *sess.UserID {
		return Error(*lib.PrivilegeError(
			fmt.Errorf("User ID %d does not match session User ID %d", id, sess.UserID),
		))
	}

	type changePasswordParams struct {
		OldPassword string
		NewPassword string
	}

	decoder := json.NewDecoder(r.Body)
	chpp := changePasswordParams{}
	err = decoder.Decode(&chpp)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding changePasswordParams: %s", err),
		))
	}

	serr := services.ChangePassword(id, chpp.OldPassword, chpp.NewPassword, db)
	if serr != nil {
		return Error(*serr)
	}

	return Success()
}

// RegisterUsersControllerRoutes Registers the functions
func RegisterUsersControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/users", serviceWrapperDBAuthenticated("getUsers", getUsers, config)).Methods("GET")
	router.HandleFunc("/users", serviceWrapperDBAuthenticated("createUser", createUser, config)).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}", serviceWrapperDBAuthenticated("getUser", getUser, config)).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}", serviceWrapperDBAuthenticated("updateUser", updateUserAPI, config)).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}/delete", serviceWrapperDBAuthenticated("deleteUser", deleteUser, config)).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}/password", serviceWrapperDBAuthenticated("changePassword", changePassword, config)).Methods("POST")

	router.HandleFunc("/app/v1/users/{id:[0-9]+}", serviceWrapperDBAuthenticated("updateUser", updateUser, config)).Methods("POST")
}
