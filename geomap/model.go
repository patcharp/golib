package geomap

import (
	"github.com/paulmach/orb/geojson"
)

type Map struct {
	featureCollection *geojson.FeatureCollection
}

type Property struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
	NameEN    string  `json:"name"`
	NameTH    string  `json:"thname"`
}
