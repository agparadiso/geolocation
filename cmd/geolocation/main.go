package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"

	geoinfo "github.com/agparadiso/geolocation/pkg/geoinfo/postgres"
	persister "github.com/agparadiso/geolocation/pkg/persistence/postgres"
	"github.com/agparadiso/geolocation/pkg/server"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	go func() {
		logrus.Println("downloading file...")
		err = downloadFile("file.csv", os.Getenv("CSV_URL"))
		if err != nil {
			logrus.Fatalf("failed to download csv: %s", err.Error())
		}

		persister := persister.New(db)
		logrus.Println("parsing and persisting file...")
		err = persister.PersistGeoinfo("file.csv")
		if err != nil {
			logrus.Errorf("failed to persist csv: %s", err.Error())
		}
	}()

	geoinfoSrv := geoinfo.New(db)
	logrus.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), server.New(geoinfoSrv)))
}

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
