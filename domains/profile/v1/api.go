package testv1

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

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

	log := echo.GetLog(ec)

	ID, err := strconv.Atoi(ec.Param("id"))
	if err != nil {
		log.Err(err).Msgf("bad request, id %s", ec.Param("id"))
		return echo.BadRequest(ec, err)
	}

	data, err := t.GetProfileByID(int64(ID))
	if err != nil {
		log.Err(err).Msgf("bad request, id %s", ec.Param("id"))
		return echo.BadRequest(ec, err)
	}

	return echo.OK(ec, ProfileDataResult{Body: data})
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
	log := echo.GetLog(ec)

	var request models.Profile

	err = ec.Bind(&request)
	if err != nil {
		log.Err(err).Msg("bad request")
		return echo.BadRequest(ec, err)
	}

	data, err := t.CreateProfile(&request)
	if err != nil {
		log.Err(err).Msgf("create data failed %+v", &request)
		return echo.CreateFailed(ec, err)
	}

	return echo.OK(ec, ProfileDataResult{Body: data})
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

	log := echo.GetLog(ec)

	ID, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		log.Err(err).Msgf("bad request, id %s", ec.Param("id"))
		return echo.BadRequest(ec, err)
	}

	hard := ec.QueryParam("hard")
	if hard == "true" {
		err = t.HardDeleteProfileByID(ID)
	} else {
		err = t.SoftDeleteProfileByID(ID)
	}

	if err != nil {
		log.Err(err).Msgf("bad request, id %s", ec.Param("id"))
		return echo.BadRequest(ec, err)
	}

	return echo.OkResult(ec)
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
	log := echo.GetLog(ec)

	profileID, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		log.Err(err).Msgf("bad request, id %s", ec.Param("id"))
		return echo.BadRequest(ec, err)
	}

	var bodyBytes []byte
	if ec.Request().Body != nil {
		bodyBytes, err = ioutil.ReadAll(ec.Request().Body)

		ec.Request().Body.Close()

		if err != nil {
			log.Err(err).Msgf("data not updated, id %d", profileID)
			return echo.BadRequest(ec, err)
		}
	}

	profileData, err := t.UpdateProfileByID(profileID, &bodyBytes)
	if err != nil {
		log.Err(err).Msgf("bad request, id %d, body %s", profileID, string(bodyBytes))
		return echo.CreateFailed(ec, err)
	}

	return echo.OK(ec, ProfileDataResult{Body: profileData})
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
	log := echo.GetLog(ec)
	var req ArrayOfMapInterface

	err = ec.Bind(&req)
	if err != nil {
		log.Err(err).Msg("bad request")
		return echo.BadRequest(ec, err)
	}

	foundUsersData, err := u.getUserDataByUserData(&req)
	if err != nil {
		log.Err(err).Msgf("not found data, request %+v", req)
		return echo.NotFound(ec, err)
	}

	return echo.OK(ec, ProfileDataResult{Body: foundUsersData})
}
