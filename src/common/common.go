package common

import "regexp"

const MaxCount int = 100

var Regex *regexp.Regexp = regexp.MustCompile("[^\\p{L}\\p{N}]+")

type GameInfo struct {
	Title          string
	FormattedTitle string
	Price          string
}
