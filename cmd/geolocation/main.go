package main

import (
	"net/http"
	"os"

	geoinfo "github.com/agparadiso/geolocation/pkg/geoinfo/mysql"
	"github.com/agparadiso/geolocation/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	geoinfoSrv := geoinfo.New()
	logrus.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), server.New(geoinfoSrv)))
}
