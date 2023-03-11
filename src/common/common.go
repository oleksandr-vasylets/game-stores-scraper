package common

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/biter777/countries"
	"golang.org/x/text/language"
)

func MaxCount() int {
	return 100
}

func CountryCode() string {
	return settings.CountryCode
}

func Locale() string {
	return settings.Locale
}

type profileSettings struct {
	CountryCode string
	Locale      string
}

var settings profileSettings

const profileSettingsFilename = "settings.json"
const defaultCountryCode = "ua"
const defaultLocale = "uk"

func init() {
	settings = profileSettings{CountryCode: defaultCountryCode, Locale: defaultLocale}
	if _, err := os.Stat(profileSettingsFilename); os.IsNotExist(err) {
		fmt.Println("settings.json not found, loading fallback values")
		file, err := os.Create(profileSettingsFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer file.Close()

		data, err := json.MarshalIndent(settings, "", "    ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		_, err = file.Write(data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
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

	cc := countries.ByName(settings.CountryCode)
	if len(settings.CountryCode) != 2 || !cc.IsValid() {
		fmt.Printf("Invalid country code, fallback to \"%s\"\n", defaultCountryCode)
		settings.CountryCode = defaultCountryCode
	}

	_, err = language.Parse(settings.Locale)
	if err != nil {
		fmt.Printf("Invalid locale, fallback to \"%s\"\n", defaultLocale)
		settings.Locale = defaultLocale
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
