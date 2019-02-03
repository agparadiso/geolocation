package main

import (
	"net/http"
	"os"

	geoinfo "github.com/agparadiso/geolocalization/pkg/geoinfo/mysql"
	"github.com/agparadiso/geolocalization/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	geoinfoSrv := geoinfo.New()
	logrus.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), server.New(geoinfoSrv)))
}
