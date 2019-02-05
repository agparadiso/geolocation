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

func New(db *sql.DB) persistence.Persister {
	return &persister{
		db: db,
	}
}

func (p *persister) PersistGeoinfo(csvURL string) error {
	start := time.Now()
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

	logrus.Println("Total time to parse and persist: ", time.Since(start))

	return nil
}

func (p *persister) ParseGeoinfo(csvReader *csv.Reader) ([]persistence.Geoinfo, error) {
	start := time.Now()
	var geoinfo []persistence.Geoinfo
	type stats struct {
		duplicated          int
		incompletecorrupted int
	}

	s := &stats{}
	alreadyImported := make(map[string]bool, 0)
	for {
		line, error := csvReader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			return nil, errors.Wrap(error, "failed while reading csv")
		}

		g := persistence.Geoinfo{
			IPaddres:     line[0],
			CountryCode:  line[1],
			Country:      line[2],
			City:         line[3],
			Latitude:     line[4],
			Longitude:    line[5],
			MisteryValue: line[6],
		}
		// ignore incompleted or corrupted records
		if !incompletedOrCorrupted(g) {
			// ignore duplicates based on the ipAddress
			if _, ok := alreadyImported[g.IPaddres]; !ok {
				sanitize(&g)
				alreadyImported[g.IPaddres] = true
				geoinfo = append(geoinfo, g)
			} else {
				s.duplicated++
			}
		} else {
			s.incompletecorrupted++
		}
	}
	logrus.Printf("Duplicated: %d, Corrupted/Incompleted %d", s.duplicated, s.incompletecorrupted)
	logrus.Println("Time to parse: ", time.Since(start))
	return geoinfo, nil
}

var insertGeoinfoQuery = `
INSERT INTO geoinfo 
(ip, country_code, country, city, latitude, longitude, mystery_value)
VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s')
`

func incompletedOrCorrupted(g persistence.Geoinfo) bool {
	if g.IPaddres == "" || g.CountryCode == "" || g.Country == "" || g.City == "" || g.Longitude == "" || g.Latitude == "" || g.MisteryValue == "" {
		return false
	}
	if _, err := strconv.ParseFloat(g.Latitude, 32); err != nil {
		return false
	}
	if _, err := strconv.ParseFloat(g.Longitude, 32); err != nil {
		return false
	}
	return true
}

func sanitize(g *persistence.Geoinfo) {
	g.City = strings.Replace(g.City, "'", "''", -1)
	g.MisteryValue = strings.Replace(g.MisteryValue, "'", "''", -1)
	g.Country = strings.Replace(g.Country, "'", "''", -1)
}
