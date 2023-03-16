package main

import (
	"bufio"
	"fmt"
	"game-stores-scraper/scrapers"
	"game-stores-scraper/settings"
	"os"
	"sort"
	"strings"

	"github.com/rodaine/table"
)

func find(title string) {
	scraperList := []scrapers.Scraper{scrapers.EpicGamesScraper{}, scrapers.SteamScraper{}, scrapers.GogScraper{}}
	columnNames := make([]interface{}, 0, len(scraperList)+2)
	columnNames = append(columnNames, "#")
	columnNames = append(columnNames, "Title")

	ch := make(chan scrapers.Result, len(scraperList))

	results := make([]scrapers.Result, len(scraperList))
	for i, scraper := range scraperList {
		columnNames = append(columnNames, scraper.GetName())
		go scraper.GetInfo(ch, i, title)
	}
	for i := 0; i < len(scraperList); i++ {
		result := <-ch
		results[result.Id] = result
	}

	data := make(map[string][]string)
	for i, result := range results {
		if result.Error != nil {
			fmt.Fprintln(os.Stderr, result.Error)
			return
		}
		sort.Slice(result.Info, func(x, y int) bool {
			return result.Info[x].FormattedTitle < result.Info[y].FormattedTitle
		})

		for _, game := range result.Info {
			if entry, ok := data[game.FormattedTitle]; ok {
				entry[i+1] = game.Price
			} else {
				entry := make([]string, len(scraperList)+1)
				entry[0] = game.Title
				for j := 1; j < len(entry); j++ {
					if j == i+1 {
						entry[j] = game.Price
					} else {
						entry[j] = "--"
					}
				}
				data[game.FormattedTitle] = entry
			}
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

	if len(keyValuePairs) == 0 {
		fmt.Println("Game(s) not found!")
		return
	}

	tbl := table.New(columnNames...)

	for i, entry := range keyValuePairs {
		row := make([]interface{}, len(entry.Value)+1)
		row[0] = i + 1
		for i, column := range entry.Value {
			row[i+1] = column
		}
		tbl.AddRow(row...)
	}

	tbl.Print()
}

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("\nSupported commands:\n" +
			"1. find [title]\n" +
			"2. get [--country | --locale]\n" +
			"3. set [--country | --locale] [value]\n" +
			"4. exit\n\n" +
			"Enter a command: ")
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, scanner.Err().Error())
			return
		}
		tokens := strings.Split(scanner.Text(), " ")
		if len(tokens) == 1 {
			if strings.ToLower(tokens[0]) == "exit" {
				return
			}
			fmt.Println("Wrong command! Try again")
			continue
		}
		if strings.ToLower(tokens[0]) == "find" {
			title := strings.Join(tokens[1:], " ")
			find(title)
			continue
		}
		if strings.ToLower(tokens[0]) == "get" {
			property := strings.ToLower(tokens[1])
			if property == "--country" {
				fmt.Println("Country:", settings.UserProfile.CountryCode)
			} else if property == "--locale" {
				fmt.Println("Locale:", settings.UserProfile.Locale)
			} else {
				fmt.Println("Wrong command! Try again")
				continue
			}
		}
		if strings.ToLower(tokens[0]) == "set" {
			if len(tokens) != 3 {
				fmt.Println("Wrong command! Try again")
				continue
			}
			property := strings.ToLower(tokens[1])
			if property == "--country" {
				settings.UserProfile.CountryCode = tokens[2]
				fmt.Println("Country changed!")
			} else if property == "--locale" {
				settings.UserProfile.Locale = tokens[2]
				fmt.Println("Locale changed!")
			} else {
				fmt.Println("Wrong command! Try again")
				continue
			}
			settings.Save()
		}
	}
}
