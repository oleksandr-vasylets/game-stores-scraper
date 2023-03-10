package common

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

func MaxCount() int {
	return settings.MaxCount
}

func CountryCode() string {
	return settings.CountryCode
}

func Locale() string {
	return settings.Locale
}

type profileSettings struct {
	MaxCount    int    `json:"MaxCount"`
	CountryCode string `json:"CountryCode"`
	Locale      string `json:"Locale"`
}

var settings profileSettings

const profileSettingsFilename = "settings.json"

func init() {
	settings = profileSettings{MaxCount: 100, CountryCode: "us", Locale: "en-US"}
	if _, err := os.Stat(profileSettingsFilename); os.IsNotExist(err) {
		return
	}
	file, err := os.Open(profileSettingsFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

var AlphanumericRegex *regexp.Regexp = regexp.MustCompile("[^\\p{L}\\p{N}]+")

type GameInfo struct {
	Title          string
	FormattedTitle string
	Price          string
}

type Result struct {
	Id    int
	Info  []GameInfo
	Error error
}

type Scraper interface {
	GetName() string
	GetInfo(ch chan Result, id int, title string)
}
