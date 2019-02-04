package postgres

import (
	"database/sql"
	"fmt"

	"github.com/agparadiso/geolocation/pkg/geoinfo"
	"github.com/pkg/errors"
)

type GeoinfoSrv struct {
	db *sql.DB
}

func New(db *sql.DB) *GeoinfoSrv {
	return &GeoinfoSrv{db: db}
}

func (g *GeoinfoSrv) GetGeoinfo(ip string) (*geoinfo.Geoinfo, error) {
	row := g.db.QueryRow(fmt.Sprintf(getInfoquery, ip))
	geoinfo := &geoinfo.Geoinfo{}
	if err := row.Scan(&geoinfo.City, &geoinfo.Country); err != nil {
		if err == sql.ErrNoRows {
			return geoinfo, nil
		}
		return nil, errors.Wrap(err, "failed to query geoinfo")
	}
	return geoinfo, nil
}

const getInfoquery = `SELECT city, country FROM geoinfo WHERE ip = '%s'`
