package testv1

import (
	"errors"
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

	data, err := t.GetProfileByID(ID)
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

	ID, err := strconv.Atoi(ec.Param("id"))
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
