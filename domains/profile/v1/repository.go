package profilev1

import (
	"encoding/json"
	"strconv"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/jmoiron/sqlx"
	"github.com/soldatov-s/go-garage-profile/models"
	"github.com/soldatov-s/go-garage/crypto/sign"
	"github.com/soldatov-s/go-garage/providers/db"
	"github.com/soldatov-s/go-garage/types"
	"github.com/soldatov-s/go-garage/x/sql"
)

const (
	profilesTable = "production.profiles"
)

func (t *ProfileV1) GetProfileByID(id int64) (*models.Profile, error) {
	data := &models.Profile{}

	if err := sql.SelectByID(t.db.Conn, profilesTable, id, &data); err != nil {
		return nil, err
	}
	t.log.Debug().Msgf("data %+v", data)

	return data, nil
}

func (t *ProfileV1) HardDeleteProfileByID(id int64) (err error) {
	return sql.HardDeleteByID(t.db.Conn, profilesTable, id)
}

func (t *ProfileV1) SoftDeleteProfileByID(id int64) (err error) {
	return sql.SoftDeleteByID(t.db.Conn, profilesTable, id)
}

func (t *ProfileV1) CreateProfile(data *models.Profile) (*models.Profile, error) {
	if t.db.Conn == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	data.CreateTimestamp()

	signature, err := sign.Generate()
	if err != nil {
		return nil, err
	}

	data.PrivateKey, data.PublicKey = sign.Marshal(signature)

	result, err := sql.InsertInto(t.db.Conn, profilesTable, data)
	if err != nil {
		return nil, err
	}

	return result.(*models.Profile), nil
}

func mergeProfileData(oldData *models.Profile, patch *[]byte) (newData *models.Profile, err error) {
	id := oldData.ID

	original, err := json.Marshal(oldData)
	if err != nil {
		return
	}

	merged, err := jsonpatch.MergePatch(original, *patch)
	if err != nil {
		return
	}

	err = json.Unmarshal(merged, &newData)
	if err != nil {
		return
	}

	// Protect ID from changes
	newData.ID = id

	newData.UpdateTimestamp()

	return newData, nil
}

func (t *ProfileV1) UpdateProfileByID(id int64, patch *[]byte) (writeData *models.Profile, err error) {
	data, err := t.GetProfileByID(id)
	if err != nil {
		return
	}

	writeData, err = mergeProfileData(data, patch)
	if err != nil {
		return
	}

	_, err = sql.Update(t.db.Conn, profilesTable, writeData)
	if err != nil {
		return nil, err
	}

	return writeData, err
}

func (u *ProfileV1) checkSearchParameter(key string) bool {
	if key == "user_id" {
		return true
	}
	for _, v := range (&models.Profile{}).SQLParamsRequest() {
		if key == v {
			return true
		}
	}

	return false
}

func (u *ProfileV1) getUserDataByUserData(req *ArrayOfMapInterface) (data *ArrayOfProfileData, err error) {
	fullQuery := "("
	queryMap := make(map[string]interface{})

	for i, item := range *req {
		var query []string

		for key, field := range item {
			if !u.checkSearchParameter(key) {
				return nil, ErrKeyDoNotMatch
			}

			queryMap[key+strconv.Itoa(i)] = field

			if field == nil {
				query = append(query, key+" is null")
				continue
			}

			if key == "user_meta" {
				queryMap[key+strconv.Itoa(i)] = field.(types.NullMeta)
			}

			query = append(query, key+"=:"+key+strconv.Itoa(i))
		}

		fullQuery = fullQuery + strings.Join(query, " and ") + ") or ("
	}

	fullQuery = strings.TrimSuffix(fullQuery, " or (")

	if u.db.Conn == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	var rows *sqlx.Rows
	if len(*req) > 0 {
		rows, err = u.db.Conn.NamedQuery("select * from production.profiles where ("+fullQuery+")", queryMap)
	} else {
		rows, err = u.db.Conn.Queryx("select * from production.profiles")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data = &ArrayOfProfileData{}

	for rows.Next() {
		var item models.Profile

		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		*data = append(*data, item)
	}

	return data, err
}
