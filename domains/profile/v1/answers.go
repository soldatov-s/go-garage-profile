package testv1

import (
	"github.com/soldatov-s/go-garage-profile/models"
	"github.com/soldatov-s/go-garage/providers/httpsrv"
)

type ProfileDataResult httpsrv.ResultAnsw

// Return array of items
type ProfilesDataResult httpsrv.ResultAnsw
type ArrayOfProfileData []models.Profile
type ArrayOfMapInterface []map[string]interface{}
