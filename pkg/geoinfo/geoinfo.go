package geoinfo

// GeoinfoSrv represents the interface of the Geoinfo Retriever
type GeoinfoSrv interface {
	GetGeoinfo(ip string) (*Geoinfo, error)
}

// Geoinfo represents the returned information
type Geoinfo struct {
	Country string `json:"country"`
	City    string `json:"city"`
}
