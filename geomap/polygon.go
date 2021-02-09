package geomap

import (
	"encoding/json"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func NewMap(GeoFile string) (Map, error) {
	m := Map{}
	b, err := ioutil.ReadFile(GeoFile)
	if err != nil {
		return m, err
	}
	m.featureCollection, err = geojson.UnmarshalFeatureCollection(b)
	if err != nil {
		return m, err
	}
	return m, nil
}

func Point(lat float64, long float64) orb.Point {
	return orb.Point{long, lat}
}

// isPointInsidePolygon runs through the MultiPolygon and Polygons within a
// feature collection and checks if a point (long/lat) lies within it.
func (m *Map) IsInside(point orb.Point) (Property, bool) {
	var p Property
	for _, feature := range m.featureCollection.Features {
		// Try on a MultiPolygon to begin
		multiPoly, isMulti := feature.Geometry.(orb.MultiPolygon)
		if isMulti {
			if planar.MultiPolygonContains(multiPoly, point) {
				if err := convertMapToProperty(feature.Properties, &p); err != nil {
					logrus.Errorln("Map struct decode error ->", err)
				}
				return p, true
			}
		} else {
			// Fallback to Polygon
			polygon, isPoly := feature.Geometry.(orb.Polygon)
			if isPoly {
				if planar.PolygonContains(polygon, point) {
					if err := convertMapToProperty(feature.Properties, &p); err != nil {
						logrus.Errorln("Map struct decode error ->", err)
					}
					return p, true
				}
			}
		}
	}
	return p, false
}

func convertMapToProperty(s map[string]interface{}, r *Property) error {
	sBytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(sBytes, r)
}
