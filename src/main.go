package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"web-scraper/epicGames"
	"web-scraper/gog"
	"web-scraper/steam"

	"github.com/rodaine/table"
)

func main() {
	fmt.Print("Enter the title of the game: ")
	var title string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, scanner.Err())
		return
	}
	title = scanner.Text()
	if title == "" {
		return
	}

	epicGamesResult, err := epicGames.GetInfo(title)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	sort.Slice(epicGamesResult, func(x, y int) bool {
		return epicGamesResult[x].FormattedTitle < epicGamesResult[y].FormattedTitle
	})

	steamResult, err := steam.GetInfo(title)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	sort.Slice(steamResult, func(x, y int) bool {
		return steamResult[x].FormattedTitle < steamResult[y].FormattedTitle
	})

	gogResult, err := gog.GetInfo(title)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	sort.Slice(gogResult, func(x, y int) bool {
		return gogResult[x].FormattedTitle < gogResult[y].FormattedTitle
	})

	data := make(map[string][]string)

	for _, game := range epicGamesResult {
		if entry, ok := data[game.FormattedTitle]; ok {
			entry[1] = game.Price
		} else {
			data[game.FormattedTitle] = []string{game.Title, game.Price, "--", "--"}
		}
	}

	for _, game := range steamResult {
		if entry, ok := data[game.FormattedTitle]; ok {
			entry[2] = game.Price
		} else {
			data[game.FormattedTitle] = []string{game.Title, "--", game.Price, "--"}
		}
	}

	for _, game := range gogResult {
		if entry, ok := data[game.FormattedTitle]; ok {
			entry[3] = game.Price
		} else {
			data[game.FormattedTitle] = []string{game.Title, "--", "--", game.Price}
		}
	}

	var keyValuePairs []struct {
		Key   string
		Value []string
	}
	for key, value := range data {
		keyValuePairs = append(keyValuePairs, struct {
			Key   string
			Value []string
		}{key, value})
	}

	sort.Slice(keyValuePairs, func(i, j int) bool {
		return keyValuePairs[i].Value[0] < keyValuePairs[j].Value[0]
	})

	tbl := table.New("#", "Title", "Epic Games Store", "Steam", "GOG")

	for i, entry := range keyValuePairs {
		tbl.AddRow(i+1, entry.Value[0], entry.Value[1], entry.Value[2], entry.Value[3])
	}

	tbl.Print()
}
