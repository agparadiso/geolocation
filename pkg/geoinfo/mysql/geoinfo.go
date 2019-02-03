package mysql

import "github.com/agparadiso/geolocalization/pkg/geoinfo"

type GeoinfoSrv struct{}

func New() *GeoinfoSrv {
	return &GeoinfoSrv{}
}

func (g *GeoinfoSrv) GetGeoinfo(ip string) (*geoinfo.Geoinfo, error) {
	return nil, nil
}
