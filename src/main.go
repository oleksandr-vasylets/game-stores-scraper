package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"web-scraper/common"
	"web-scraper/epicGames"
	"web-scraper/gog"
	"web-scraper/steam"

	"github.com/rodaine/table"
)

func main() {
	fmt.Print("Enter the title of the game (more specific titles yield better results): ")
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

	scrapers := []common.Scraper{epicGames.Scraper{}, steam.Scraper{}, gog.Scraper{}}
	columnNames := make([]interface{}, 0, len(scrapers)+2)
	columnNames = append(columnNames, "#")
	columnNames = append(columnNames, "Title")

	ch := make(chan common.Result, len(scrapers))

	results := make([]common.Result, len(scrapers))
	for i, scraper := range scrapers {
		columnNames = append(columnNames, scraper.GetName())
		go scraper.GetInfo(ch, i, title)
	}
	for i := 0; i < len(scrapers); i++ {
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
				entry := make([]string, len(scrapers)+1)
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

	fmt.Scanln()
}
