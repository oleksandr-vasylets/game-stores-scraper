package scrapers

import (
	"regexp"
)

var alphanumericRegex *regexp.Regexp = regexp.MustCompile("[^\\p{L}\\p{N}]+")

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
