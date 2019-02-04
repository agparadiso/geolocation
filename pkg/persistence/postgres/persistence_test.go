package postgres

import (
	"encoding/csv"
	"reflect"
	"strings"
	"testing"

	"github.com/agparadiso/geolocation/pkg/persistence"
)

func TestParseGeoinfo(t *testing.T) {
	in := `200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346`
	expectedGeoinfo := []persistence.Geoinfo{persistence.Geoinfo{
		IPaddres:     "200.106.141.15",
		CountryCode:  "SI",
		Country:      "Nepal",
		City:         "DuBuquemouth",
		Latitude:     "-84.87503094689836",
		Longitude:    "7.206435933364332",
		MisteryValue: "7823011346",
	}}
	r := csv.NewReader(strings.NewReader(in))
	persister := NewPersister(nil)
	geoinfo, err := persister.ParseGeoinfo(r)
	if err != nil {
		t.Fatal("failed to parse geoinfo: ", err)
	}

	if !reflect.DeepEqual(geoinfo, expectedGeoinfo) {
		t.Fatalf("expected: %v, got: %v", expectedGeoinfo, geoinfo)
	}

}