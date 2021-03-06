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
	"sync"
	"time"

	"github.com/agparadiso/geolocation/pkg/persistence"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

type stats struct {
	duplicated          int
	incompletecorrupted int
}

const (
	duplicated          = 0
	incompletecorrupted = 1
)

func (p *persister) PersistGeoinfo(csvURL string) error {
	s := &stats{}
	start := time.Now()
	csvFile, err := os.Open(csvURL)
	if err != nil {
		return errors.Wrap(err, "failed to open csv")
	}
	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	linesChan := make(chan []string)
	errChan := make(chan error, 1)
	counterChan := make(chan int, 1)
	var counterWg sync.WaitGroup
	go counters(s, counterChan, &counterWg)
	var wg sync.WaitGroup
	for w := 0; w < 10; w++ {
		go worker(p, linesChan, errChan, counterChan, s, &wg)
	}

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			close(linesChan)
			break
		} else if err != nil {
			close(linesChan)
			close(errChan)
			return errors.Wrap(err, "failed while reading csv")
		}

		select {
		case err := <-errChan:
			return err
		default:
		}
		linesChan <- line
	}
	wg.Wait()
	close(counterChan)
	counterWg.Wait()
	close(errChan)
	logrus.Printf("Duplicated: %d, Corrupted/Incompleted: %d", s.duplicated, s.incompletecorrupted)
	logrus.Println("Total time to parse and persist: ", time.Since(start))
	return nil
}

func worker(p *persister, linesChan <-chan []string, errChan chan<- error, counterChan chan int, s *stats, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for line := range linesChan {
		g, err := parseGeoinfo(line)
		if err != nil {
			counterChan <- incompletecorrupted
			continue
		}
		sanitize(&g)
		_, err = p.db.Exec(fmt.Sprintf(insertGeoinfoQuery, g.IPaddres, g.CountryCode, g.Country, g.City, g.Latitude, g.Longitude, g.MisteryValue))
		if err != nil {
			// ignore duplicates errors
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				counterChan <- duplicated
				continue
			}
			errChan <- errors.Wrapf(err, "failed to persist geoinfo %v", g)
		}
	}
}

func counters(s *stats, operation <-chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for op := range operation {
		switch op {
		case duplicated:
			s.duplicated++
		case incompletecorrupted:
			s.incompletecorrupted++
		}
	}
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
