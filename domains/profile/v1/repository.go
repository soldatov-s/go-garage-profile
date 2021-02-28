package testv1

import (
	"strings"

	"github.com/soldatov-s/go-garage-profile/models"
	"github.com/soldatov-s/go-garage/crypto/sign"
	"github.com/soldatov-s/go-garage/providers/db"
	"github.com/soldatov-s/go-garage/utils"
	"github.com/soldatov-s/go-garage/x/sql"
)

func (t *ProfileV1) GetProfileByID(id int) (*models.Profile, error) {
	data := &models.Profile{}

	if err := sql.SelectByID(t.db.Conn, "production.profile", int64(id), &data); err != nil {
		return nil, err
	}
	t.log.Debug().Msgf("data %+v", data)

	return data, nil
}

func (t *ProfileV1) HardDeleteProfileByID(id int) (err error) {
	return sql.HardDeleteByID(t.db.Conn, "production.profile", int64(id))
}

func (t *ProfileV1) SoftDeleteProfileByID(id int) (err error) {
	data, err := t.GetProfileByID(id)
	if err != nil {
		return err
	}

	if data.DeletedAt.Valid {
		return nil
	}

	data.DeletedAt.Timestamp()

	query := make([]string, 0, len(data.SQLParamsRequest()))
	for _, param := range data.SQLParamsRequest() {
		query = append(query, param+"=:"+param)
	}

	if t.db.Conn == nil {
		return db.ErrDBConnNotEstablished
	}

	_, err = t.db.Conn.NamedExec(
		t.db.Conn.Rebind(utils.JoinStrings(" ", "UPDATE production.profile SET", strings.Join(query, ", "), "WHERE id=:id")),
		data)

	if err != nil {
		return err
	}

	return nil
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

	result, err := sql.InsertInto(t.db.Conn, "production.profile", data)
	if err != nil {
		return nil, err
	}

	return result.(*models.Profile), nil
}
