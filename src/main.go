package main

import (
	"bufio"
	"fmt"
	"os"
	"web-scraper/epicGames"
)

func main() {
	fmt.Print("Enter the name of the game: ")
	var name string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, scanner.Err().Error())
		return
	}
	name = scanner.Text()

	games, err := epicGames.GetInfo(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i, game := range games {
		fmt.Printf("%d. %s\nPrice: %s", i+1, game.Title, game.Price)
		fmt.Println()
	}
}
