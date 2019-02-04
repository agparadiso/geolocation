package postgres

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/agparadiso/geolocation/pkg/persistence"
	"github.com/pkg/errors"
)

type persister struct {
	db *sql.DB
}

func New(db *sql.DB) persistence.Persister {
	return &persister{
		db: db,
	}
}

func (p *persister) PersistGeoinfo(csvURL string) error {
	csvFile, _ := os.Open(csvURL)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	geoinfo, err := p.ParseGeoinfo(reader)
	if err != nil {
		return errors.Wrap(err, "failed to parse Geoinfo")
	}

	for _, g := range geoinfo {
		_, err := p.db.Exec(fmt.Sprintf(insertGeoinfoQuery, g.IPaddres, g.CountryCode, g.Country, g.City, g.Latitude, g.Longitude, g.MisteryValue))
		if err != nil {
			return errors.Wrapf(err, "failed to persist geoinfo %v", g)
		}
	}

	return nil
}

func (p *persister) ParseGeoinfo(csvReader *csv.Reader) ([]persistence.Geoinfo, error) {
	var geoinfo []persistence.Geoinfo
	for {
		line, error := csvReader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			return nil, errors.Wrap(error, "failed while reading csv")
		}
		geoinfo = append(geoinfo, persistence.Geoinfo{
			IPaddres:     line[0],
			CountryCode:  line[1],
			Country:      line[2],
			City:         line[3],
			Latitude:     line[4],
			Longitude:    line[5],
			MisteryValue: line[6],
		})
	}
	return geoinfo, nil
}

var insertGeoinfoQuery = `
INSERT INTO geoinfo 
(ip, country_code, country, city, latitude, longitude, mystery_value)
VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s')
`
