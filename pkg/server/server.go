package server

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/agparadiso/geolocation/pkg/geoinfo"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	geoinfoSrv geoinfo.GeoinfoSrv
}

func New(geoinfoSrv geoinfo.GeoinfoSrv) http.Handler {
	s := &Server{geoinfoSrv: geoinfoSrv}
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1/").Subrouter()
	api.HandleFunc(`/geoinfo`, s.getGeoinfo)
	handler := cors.Default().Handler(r)
	return handler
}

func (s *Server) getGeoinfo(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		logrus.Error("ip parameter not entered")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var geoinfo *geoinfo.Geoinfo
	geoinfo, err := s.geoinfoSrv.GetGeoinfo(ip)
	if err != nil {
		logrus.Error("failed to get geoinfo response: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(geoinfo)
	if err != nil {
		logrus.Error("failed to Marshal geoinfo Response: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
