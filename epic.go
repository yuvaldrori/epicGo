package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var apiKey string = os.Getenv("NASAAPIKEY")

const nasaUrl = "https://api.nasa.gov/EPIC"
const epicUrl = "http://localhost:8080"

type Date struct {
	Date string
}

type Coordinates struct {
	Lat, Lon float64
}

type Position struct {
	X, Y, Z float64
}

type Quaternions struct {
	Q0, Q1, Q2, Q3 float64
}

type Coords struct {
	Centroid_coordinates  Coordinates
	Dscovr_j2000_position Position
	Lunar_j2000_position  Position
	Sun_j2000_position    Position
	Attitude_quaternions  Quaternions
}

type Image struct {
	Identifier            string
	Caption               string
	Image                 string
	Version               string
	Centroid_coordinates  Coordinates
	Dscovr_j2000_position Position
	Lunar_j2000_position  Position
	Sun_j2000_position    Position
	Attitude_quaternions  Quaternions
	Date                  string
	Coords                Coords
}

func main() {
	nasa, err := AvailableDates(fmt.Sprintf("%s/api/natural/all?api_key=%s", nasaUrl, apiKey))
	if err != nil {
		log.Fatal(err)
	}
	epic, err := AvailableDates(fmt.Sprintf("%s/all.json", epicUrl))
	if err != nil {
		log.Fatal(err)
	}
	missingDates := MissingDates(nasa, epic)
	for index, date := range missingDates {
		if index > 1 {
			return
		}
		images, err := getDate(date)
		if err != nil {
			log.Fatal(err)
		}
		for _, image := range images {
			name := image.Image
			url, err := getImageUrl(image)
			if err != nil {
				log.Fatal(err)
			}
			file, err := DownloadImage(url, name)

		}
	}
}

func AvailableDates(url string) ([]Date, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	dates, err := ParseDates(content)
	if err != nil {
		log.Fatal(err)
	}
	return dates, err
}

func ParseDates(content []byte) (dates []Date, err error) {
	err = json.Unmarshal(content, &dates)
	if err != nil {
		log.Fatal(err)
	}
	return dates, err
}

func MissingDates(nasa, epic []Date) []string {
	m := make(map[string]struct{})
	for _, date := range epic {
		m[date.Date] = struct{}{}
	}
	var missingDates []string
	for _, date := range nasa {
		if _, ok := m[date.Date]; ok {
		} else {
			missingDates = append(missingDates, date.Date)
		}
	}
	return missingDates
}

func Resize(original string, size string, name string) (out string, err error) {
	out = filepath.Join(os.TempDir(), fmt.Sprintf("%s_%sx%s.jpg", name, size, size))
	cmd := exec.Command("convert", original, "-resize", fmt.Sprintf("%sx%s", size, size), out)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out, err
}

func DownloadImage(url string, name string) (out string, err error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	out = filepath.Join(os.TempDir(), fmt.Sprintf("%s.png", name))
	file, err := os.Create(out)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return out, err
}

func getDate(date string) ([]Image, error) {
	url := fmt.Sprintf("%s/api/natural/date/%s?api_key=%s", nasaUrl, date, apiKey)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return ParseDate(body)
}

func getImageUrl(image Image) (string, error) {
	t, err := time.Parse("2006-01-02 15:04:05", image.Date)
	if err != nil {
		log.Fatal(err)
	}
	year := t.Year()
	month := t.Month()
	day := t.Day()
	name := image.Image
	url := fmt.Sprintf("%s/archive/natural/%d/%d/%d/png/%s.png?api_key=%s", nasaUrl, year, month, day, name, apiKey)
	return url, err
}

func ParseDate(body []byte) (images []Image, err error) {
	err = json.Unmarshal(body, &images)
	if err != nil {
		log.Fatal(err)
	}
	return images, err
}

func upload(content []byte, path string) error {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, content)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
