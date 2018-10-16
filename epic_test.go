package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseDates(t *testing.T) {
	c := struct {
		in   []byte
		want []Date
	}{
		[]byte(`[{"date": "2015-06-13"}, {"date": "2015-06-16"}]`),
		[]Date{{Date: "2015-06-13"}, {Date: "2015-06-16"}},
	}
	got, err := ParseDates(c.in)
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(got, c.want) != true {
		t.Errorf("ParseDates(\"%s\")\ngot\t%+v\nwant\t%+v", c.in, got, c.want)
	}
}

func TestAvailableDates(t *testing.T) {
	c := struct {
		in   string
		want []Date
	}{
		"http://localhost:8080/dates.json",
		[]Date{{Date: "2018-09-20"}, {Date: "2018-09-19"}},
	}
	got, err := AvailableDates(c.in)
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(got, c.want) != true {
		t.Errorf("AvailableDates(\"%s\")\ngot\t%+v\nwant\t%+v", c.in, got, c.want)
	}
}

func TestMissingDates(t *testing.T) {
	cases := []struct {
		nasa, epic []Date
		want       []string
	}{
		{
			[]Date{{Date: "2018-09-20"}, {Date: "2018-09-19"}},
			[]Date{{Date: "2018-09-20"}},
			[]string{"2018-09-19"},
		},
		{
			[]Date{{Date: "2018-09-20"}, {Date: "2018-09-19"}},
			[]Date{{Date: "2018-09-20"}, {Date: "2018-09-19"}},
			nil,
		},
	}
	for _, c := range cases {
		got := MissingDates(c.nasa, c.epic)
		if reflect.DeepEqual(got, c.want) != true {
			t.Errorf("MissingDates(%v, %v)\ngot\t%+v\nwant\t%v", c.nasa, c.epic, got, c.want)
		}
	}
}

func TestResize(t *testing.T) {
	name := "epic_1b_20151031003633"
	original := filepath.Join("testdata", "epic_1b_20151031003633.png")
	out, err := Resize(original, "256", name)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(out)
	os.Remove(out)
}

func TestDownloadImage(t *testing.T) {
	var apiKey string = os.Getenv("NASAAPIKEY")
	url := fmt.Sprintf("https://api.nasa.gov/EPIC/archive/natural/2015/10/31/png/epic_1b_20151031003633.png?api_key=%s", apiKey)
	name := "epic_1b_20151031003633"
	out, err := DownloadImage(url, name)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(out)
	os.Remove(out)
}

func TestParseDate(t *testing.T) {
	c := struct {
		in   []byte
		want []Image
	}{
		[]byte(`[{"identifier": "20151031221308","caption": "This image was taken by NASA's EPIC camera onboard the NOAA DSCOVR spacecraft","image": "epic_1b_20151031221308","version": "02","centroid_coordinates": {"lat": -5.36902,"lon": -164.487468},"dscovr_j2000_position": {"x": -1268623.557569,"y": -690889.141203,"z": -136798.451041},"lunar_j2000_position": {"x": -47188.349357,"y": 358516.852648,"z": 118163.91538},"sun_j2000_position": {"x": -117110114.80008,"y": -83785061.649627,"z": -36321513.174041},"attitude_quaternions": {"q0": -0.302,"q1": -0.123,"q2": 0.21848,"q3": 0.91975},"date": "2015-10-31 22:08:19","coords": {"centroid_coordinates": {"lat": -5.36902,"lon": -164.487468},"dscovr_j2000_position": {"x": -1268623.557569,"y": -690889.141203, "z": -136798.451041},"lunar_j2000_position": {"x": -47188.349357,"y": 358516.852648,"z": 118163.91538},"sun_j2000_position": {"x": -117110114.80008,"y": -83785061.649627,"z": -36321513.174041},"attitude_quaternions": {"q0": -0.302,"q1": -0.123,"q2": 0.21848,"q3": 0.91975}}}]`),
		[]Image{{
			Identifier: "20151031221308",
			Caption:    "This image was taken by NASA's EPIC camera onboard the NOAA DSCOVR spacecraft",
			Image:      "epic_1b_20151031221308",
			Version:    "02",
			Centroid_coordinates: Coordinates{
				Lat: -5.36902,
				Lon: -164.487468,
			},
			Dscovr_j2000_position: Position{
				X: -1268623.557569,
				Y: -690889.141203,
				Z: -136798.451041,
			},
			Lunar_j2000_position: Position{
				X: -47188.349357,
				Y: 358516.852648,
				Z: 118163.91538,
			},
			Sun_j2000_position: Position{
				X: -117110114.80008,
				Y: -83785061.649627,
				Z: -36321513.174041,
			},
			Attitude_quaternions: Quaternions{
				Q0: -0.302,
				Q1: -0.123,
				Q2: 0.21848,
				Q3: 0.91975,
			},
			Date: "2015-10-31 22:08:19",
			Coords: Coords{
				Centroid_coordinates: Coordinates{
					Lat: -5.36902,
					Lon: -164.487468,
				},
				Dscovr_j2000_position: Position{
					X: -1268623.557569,
					Y: -690889.141203,
					Z: -136798.451041,
				},
				Lunar_j2000_position: Position{
					X: -47188.349357,
					Y: 358516.852648,
					Z: 118163.91538,
				},
				Sun_j2000_position: Position{
					X: -117110114.80008,
					Y: -83785061.649627,
					Z: -36321513.174041,
				},
				Attitude_quaternions: Quaternions{
					Q0: -0.302,
					Q1: -0.123,
					Q2: 0.21848,
					Q3: 0.91975,
				},
			},
		}},
	}
	got, err := ParseDate(c.in)
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(got, c.want) != true {
		t.Errorf("ParseDate(\"%s\")\ngot\t%+v\nwant\t%+v", c.in, got, c.want)
	}
}
