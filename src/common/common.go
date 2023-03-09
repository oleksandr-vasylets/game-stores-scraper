package common

import "regexp"

const MaxCount int = 100

const CountryCode string = "ua" // TODO: Replace these with actual user country code and locale
const Locale string = "uk"

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
