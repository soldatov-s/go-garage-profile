package testv1

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/soldatov-s/go-garage-profile/models"

	"github.com/soldatov-s/go-garage/providers/httpsrv"
	"github.com/soldatov-s/go-garage/providers/httpsrv/echo"
	echoSwagger "github.com/soldatov-s/go-swagger/echo-swagger"
)

func (t *ProfileV1) profileGetHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		err = errors.New("error")
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("This handler getting data for requested ID").
			SetSummary("Get data by ID").
			AddInPathParameter("id", "ID", reflect.Int).
			AddResponse(http.StatusOK, "Data", ProfileDataResult{Body: models.Profile{}}).
			AddResponse(http.StatusBadRequest, "BAD REQUEST", httpsrv.BadRequest(err)).
			AddResponse(http.StatusNotFound, "NOT FOUND DATA", httpsrv.NotFound(err))

		return nil
	}

	log := ec.GetLog()

	ID, err := ec.GetInt64Param("id")
	if err != nil {
		log.Err(err).Msgf("bad request, id %s", ec.Param("id"))
		return ec.BadRequest(err)
	}

	data, err := t.GetProfileByID(ID)
	if err != nil {
		log.Err(err).Msgf("bad request, id %s", ec.Param("id"))
		return ec.BadRequest(err)
	}

	return ec.OK(ProfileDataResult{Body: data})
}

func (t *ProfileV1) profilePostHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		err = errors.New("error")
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("This handler create new data").
			SetSummary("Create Data Handler").
			AddInBodyParameter("data", "Data", models.Profile{}, true).
			AddResponse(http.StatusOK, "Data", &ProfileDataResult{Body: models.Profile{}}).
			AddResponse(http.StatusBadRequest, "BAD REQUEST", httpsrv.BadRequest(err)).
			AddResponse(http.StatusConflict, "CREATE DATA FAILED", httpsrv.CreateFailed(err))

		return nil
	}

	// Main code of handler
	log := ec.GetLog()

	var request models.Profile

	err = ec.Bind(&request)
	if err != nil {
		log.Err(err).Msg("bad request")
		return ec.BadRequest(err)
	}

	data, err := t.CreateProfile(&request)
	if err != nil {
		log.Err(err).Msgf("create data failed %+v", &request)
		return ec.CreateFailed(err)
	}

	return ec.OK(ProfileDataResult{Body: data})
}

func (t *ProfileV1) profileDeleteHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		err = errors.New("error")
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("This handler deletes data for requested ID").
			SetSummary("Delete data by ID").
			AddInPathParameter("id", "ID", reflect.Int).
			AddInQueryParameter("hard", "Hard delete user, if equal true, delete hard", reflect.Bool, false).
			AddResponse(http.StatusOK, "OK", httpsrv.OkResult()).
			AddResponse(http.StatusBadRequest, "BAD REQUEST", httpsrv.BadRequest(err)).
			AddResponse(http.StatusNotFound, "NOT FOUND DATA", httpsrv.NotFound(err))

		return nil
	}

	log := ec.GetLog()

	ID, err := ec.GetInt64Param("id")
	if err != nil {
		log.Err(err).Msgf("bad request, id %s", ec.Param("id"))
		return ec.BadRequest(err)
	}

	hard := ec.QueryParam("hard")
	if hard == "true" {
		err = t.HardDeleteProfileByID(ID)
	} else {
		err = t.SoftDeleteProfileByID(ID)
	}

	if err != nil {
		log.Err(err).Msgf("bad request, id %s", ec.Param("id"))
		return ec.BadRequest(err)
	}

	return ec.OkResult()
}

func (t *ProfileV1) profilePutHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		err = fmt.Errorf("error")
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("Update User Handler").
			SetSummary("This handler update user data by user_id").
			AddInBodyParameter("user_data", "User data", &models.Profile{}, true).
			AddInPathParameter("id", "User id", reflect.Int64).
			AddResponse(http.StatusOK, "User Data", &ProfileDataResult{Body: models.Profile{}}).
			AddResponse(http.StatusBadRequest, "BAD REQUEST", httpsrv.BadRequest(err)).
			AddResponse(http.StatusConflict, "DATA NOT UPDATED", httpsrv.NotUpdated(err)).
			AddResponse(http.StatusNotFound, "NOT FOUND DATA", httpsrv.NotFound(err))

		return nil
	}

	// Main code of handler
	log := ec.GetLog()

	profileID, err := ec.GetInt64Param("id")
	if err != nil {
		log.Err(err).Msgf("bad request, id %s", ec.Param("id"))
		return ec.BadRequest(err)
	}

	var bodyBytes []byte
	if ec.Request().Body != nil {
		bodyBytes, err = ioutil.ReadAll(ec.Request().Body)

		ec.Request().Body.Close()

		if err != nil {
			log.Err(err).Msgf("data not updated, id %d", profileID)
			return ec.BadRequest(err)
		}
	}

	profileData, err := t.UpdateProfileByID(profileID, &bodyBytes)
	if err != nil {
		log.Err(err).Msgf("bad request, id %d, body %s", profileID, string(bodyBytes))
		return ec.CreateFailed(err)
	}

	return ec.OK(ProfileDataResult{Body: profileData})
}

func (u *ProfileV1) profileSearchPostHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		err = fmt.Errorf("error")
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("Find profile").
			SetSummary("This handler find profile data by any field in Profile data struct. Can be multiple structs in request. Search by meta-fields not work!").
			AddInBodyParameter("users_data", "Users data", &ArrayOfProfileData{}, true).
			AddResponse(http.StatusOK, "Users data", &ProfilesDataResult{Body: ArrayOfProfileData{}}).
			AddResponse(http.StatusBadRequest, "BAD REQUEST", httpsrv.BadRequest(err)).
			AddResponse(http.StatusNotFound, "NOT FOUND DATA", httpsrv.NotFound(err))

		return nil
	}

	// Main code of handler
	log := ec.GetLog()
	var req ArrayOfMapInterface

	err = ec.Bind(&req)
	if err != nil {
		log.Err(err).Msg("bad request")
		return ec.BadRequest(err)
	}

	foundUsersData, err := u.getUserDataByUserData(&req)
	if err != nil {
		log.Err(err).Msgf("not found data, request %+v", req)
		return ec.NotFound(err)
	}

	return ec.OK(ProfileDataResult{Body: foundUsersData})
}
