package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	geoinfo "github.com/agparadiso/geolocation/pkg/geoinfo/postgres"
	"github.com/agparadiso/geolocation/pkg/server"
	"github.com/sirupsen/logrus"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "postgres"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	geoinfoSrv := geoinfo.New(db)
	logrus.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), server.New(geoinfoSrv)))
}
