package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"web-scraper/epicGames"
	"web-scraper/steam"

	"github.com/rodaine/table"
)

func main() {
	fmt.Print("Enter the title of the game: ")
	var title string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, scanner.Err().Error())
		return
	}
	title = scanner.Text()

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

	type Row struct {
		Title       string
		FirstPrice  string
		SecondPrice string
	}

	combined := make([]Row, 0, len(epicGamesResult)+len(steamResult))
	i, j := 0, 0
	for i < len(epicGamesResult) && j < len(steamResult) {
		l := epicGamesResult[i]
		r := steamResult[j]
		if l.FormattedTitle < r.FormattedTitle {
			combined = append(combined, Row{Title: l.Title, FirstPrice: l.Price, SecondPrice: "--"})
			i++
			continue
		} else if r.FormattedTitle < l.FormattedTitle {
			combined = append(combined, Row{Title: r.Title, FirstPrice: "--", SecondPrice: r.Price})
			j++
			continue
		}
		combined = append(combined, Row{Title: l.Title, FirstPrice: l.Price, SecondPrice: r.Price})
		i++
		j++
	}

	for i < len(epicGamesResult) {
		combined = append(combined, Row{Title: epicGamesResult[i].Title, FirstPrice: epicGamesResult[i].Price, SecondPrice: "--"})
		i++
	}

	for j < len(steamResult) {
		combined = append(combined, Row{Title: steamResult[j].Title, FirstPrice: "--", SecondPrice: steamResult[j].Price})
		j++
	}

	tbl := table.New("#", "Title", "Epic Games Store", "Steam")

	for i, game := range combined {
		tbl.AddRow(i+1, game.Title, game.FirstPrice, game.SecondPrice)
	}

	tbl.Print()
}
