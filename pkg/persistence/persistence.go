package persistence

type Persister interface {
	PersistGeoinfo(csvURL string) error
}

// Geoinfo represents the geoinfo persisted in db
type Geoinfo struct {
	IPaddres     string `json:"ip_adress"`
	CountryCode  string `json:"country_code"`
	Country      string `json:"country"`
	City         string `json:"city"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	MisteryValue string `json:"mystery_value"`
}
