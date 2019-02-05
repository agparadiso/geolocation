package postgres

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/agparadiso/geolocation/pkg/persistence"
	"github.com/pkg/errors"
)

type persister struct {
	db *sql.DB
}

var (
	errIncompletedField    = errors.New("incompleted Field")
	errLatitudeNotFloat32  = errors.New("latitude is not a float32")
	errLongitudeNotFloat32 = errors.New("longitude is not a float32")
)

func New(db *sql.DB) persistence.Persister {
	return &persister{
		db: db,
	}
}

func (p *persister) PersistGeoinfo(csvURL string) error {
	type stats struct {
		duplicated          int
		incompletecorrupted int
	}

	s := &stats{}
	start := time.Now()
	csvFile, err := os.Open(csvURL)
	if err != nil {
		return errors.Wrap(err, "failed to open csv")
	}
	csvReader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrap(err, "failed while reading csv")
		}

		g, err := parseGeoinfo(line)
		if err != nil {
			s.incompletecorrupted++
			continue
		}

		sanitize(&g)
		_, err = p.db.Exec(fmt.Sprintf(insertGeoinfoQuery, g.IPaddres, g.CountryCode, g.Country, g.City, g.Latitude, g.Longitude, g.MisteryValue))
		if err != nil {
			// ignore duplicates errors
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				s.duplicated++
				continue
			}
			return errors.Wrapf(err, "failed to persist geoinfo %v", g)
		}
	}
	logrus.Printf("Duplicated: %d, Corrupted/Incompleted %d", s.duplicated, s.incompletecorrupted)
	logrus.Println("Total time to parse and persist: ", time.Since(start))
	return nil
}

func parseGeoinfo(line []string) (persistence.Geoinfo, error) {
	g := persistence.Geoinfo{
		IPaddres:     line[0],
		CountryCode:  line[1],
		Country:      line[2],
		City:         line[3],
		Latitude:     line[4],
		Longitude:    line[5],
		MisteryValue: line[6],
	}

	return g, incompletedCorrupted(g)
}

var insertGeoinfoQuery = `
INSERT INTO geoinfo 
(ip, country_code, country, city, latitude, longitude, mystery_value)
VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s')
`

func incompletedCorrupted(g persistence.Geoinfo) error {
	if g.IPaddres == "" || g.CountryCode == "" || g.Country == "" || g.City == "" || g.Longitude == "" || g.Latitude == "" || g.MisteryValue == "" {
		return errIncompletedField
	}
	if _, err := strconv.ParseFloat(g.Latitude, 32); err != nil {
		return errLatitudeNotFloat32
	}
	if _, err := strconv.ParseFloat(g.Longitude, 32); err != nil {
		return errLongitudeNotFloat32
	}
	return nil
}

func sanitize(g *persistence.Geoinfo) {
	g.City = strings.Replace(g.City, "'", "''", -1)
	g.MisteryValue = strings.Replace(g.MisteryValue, "'", "''", -1)
	g.Country = strings.Replace(g.Country, "'", "''", -1)
}
