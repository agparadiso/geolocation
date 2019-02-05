package postgres

import (
	"reflect"
	"strings"
	"testing"

	"github.com/agparadiso/geolocation/pkg/persistence"
)

func TestParseGeoinfo(t *testing.T) {
	in := `200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346`
	expectedGeoinfo := persistence.Geoinfo{
		IPaddres:     "200.106.141.15",
		CountryCode:  "SI",
		Country:      "Nepal",
		City:         "DuBuquemouth",
		Latitude:     "-84.87503094689836",
		Longitude:    "7.206435933364332",
		MisteryValue: "7823011346",
	}
	geoinfo, err := parseGeoinfo(strings.Split(in, ","))
	if err != nil {
		t.Fatal("failed to parse geoinfo: ", err)
	}

	if !reflect.DeepEqual(geoinfo, expectedGeoinfo) {
		t.Fatalf("expected: %v, got: %v", expectedGeoinfo, geoinfo)
	}
}

func TestIncompletedCorrupted(t *testing.T) {
	type testCase struct {
		in       persistence.Geoinfo
		expected error
	}

	cases := []testCase{
		{
			in: persistence.Geoinfo{
				IPaddres:     "200.106.141.15",
				CountryCode:  "SI",
				Country:      "Nepal",
				City:         "DuBuquemouth",
				Latitude:     "-84.87503094689836",
				Longitude:    "7.206435933364332",
				MisteryValue: "7823011346",
			},
			expected: nil,
		},
		{
			in: persistence.Geoinfo{
				IPaddres:     "",
				CountryCode:  "SI",
				Country:      "Nepal",
				City:         "DuBuquemouth",
				Latitude:     "-84.87503094689836",
				Longitude:    "7.206435933364332",
				MisteryValue: "7823011346",
			},
			expected: errIncompletedField,
		},
		{
			in: persistence.Geoinfo{
				IPaddres:     "200.106.141.15",
				CountryCode:  "SI",
				Country:      "Nepal",
				City:         "DuBuquemouth",
				Latitude:     "not a number",
				Longitude:    "7.206435933364332",
				MisteryValue: "7823011346",
			},
			expected: errLatitudeNotFloat32,
		},
	}

	for _, c := range cases {
		actual := incompletedCorrupted(c.in)
		if c.expected != actual {
			t.Fatalf("expected: %v on geoinfo: %v, got: %v", c.expected, c.in, actual)
		}
	}
}
