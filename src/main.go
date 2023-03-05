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

	price, err := epicGames.GetPrice(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(price)
}
