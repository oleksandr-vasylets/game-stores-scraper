package main

import (
	"bufio"
	"fmt"
	"game-stores-scraper/scrapers"
	"os"
	"sort"

	"github.com/rodaine/table"
)

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter the title of the game (more specific titles yield better results). Or just press Enter to exit: ")
		var title string
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Fprintln(os.Stderr, scanner.Err())
			reader.ReadLine()
			return
		}
		title = scanner.Text()
		if title == "" {
			return
		}

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
				reader.ReadLine()
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
			fmt.Println()
			continue
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
		fmt.Println()
	}
}
