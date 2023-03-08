package main

import (
	"bufio"
	"fmt"
	"os"
	"web-scraper/common"
	"web-scraper/epicGames"
	"web-scraper/steam"
)

func printGameList(list []common.GameInfo) {
	for i, game := range list {
		fmt.Printf("%d. %s\nPrice: %s", i+1, game.Title, game.Price)
		fmt.Println()
	}
}

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

	games, err := epicGames.GetInfo(title)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	fmt.Println("----- Epic Games Store -----")
	printGameList(games)

	games, err = steam.GetInfo(title)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	fmt.Println("----- Steam -----")
	printGameList(games)
}
