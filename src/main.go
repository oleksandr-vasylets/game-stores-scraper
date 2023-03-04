package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

type AppListResponse struct {
	List struct {
		Apps []struct {
			AppId uint32 `json:"appid"`
			Name  string `json:"name"`
		} `json:"apps"`
	} `json:"applist"`
}

func getPrice(name string) (string, error) {
	resp, err := http.Get("https://api.steampowered.com/ISteamApps/GetAppList/v2/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var appList AppListResponse
	err = json.Unmarshal(body, &appList)
	if err != nil {
		return "", err
	}

	type app struct {
		dist  int
		price string
	}

	// Find exact match or the closest matches to the user's input using Levenshtein distance.
	matches := make(map[string]app)
	for _, elem := range appList.List.Apps {
		if strings.ToLower(elem.Name) == strings.ToLower(name) {
			price, err := fetchPrice(elem.AppId)
			if err != nil {
				return "", err
			}
			if price != "" {
				return price, nil
			}
		}
		dist := levenshtein.DistanceForStrings([]rune(elem.Name), []rune(name), levenshtein.DefaultOptions)
		if dist <= len(name)/2 || strings.Contains(strings.ToLower(elem.Name), strings.ToLower(name)) {
			price, err := fetchPrice(elem.AppId)
			if err != nil {
				return "", err
			}
			if price != "" {
				matches[elem.Name] = app{dist, price}
			}
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("No matches found! Try more specific name.")
	}

	// Sort the matches by distance.
	type match struct {
		name  string
		price string
		dist  int
	}
	var sortedMatches []match
	for name, elem := range matches {
		sortedMatches = append(sortedMatches, match{name, elem.price, elem.dist})
	}
	sort.Slice(sortedMatches, func(i, j int) bool {
		return sortedMatches[i].dist < sortedMatches[j].dist
	})

	// Display the closest matches to the user's input.
	fmt.Println("Closest matches:")
	for i, m := range sortedMatches {
		fmt.Printf("%d. %s\n", i+1, m.name)
	}

	// Ask the user to select a game from the closest matches.
	for {
		var selectedGame string
		fmt.Print("Enter the name of a game from the list above: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Fprintln(os.Stderr, scanner.Err().Error())
			return "", scanner.Err()
		}
		selectedGame = scanner.Text()

		// Find the appid for the selected game.
		for _, app := range sortedMatches {
			if app.name == selectedGame {
				return app.price, nil
			}
		}
		fmt.Println("No such game in the list! Please, try again.")
	}
}

type PriceResponse struct {
	Data struct {
		PriceOverview struct {
			Price string `json:"final_formatted"`
		} `json:"price_overview"`
	} `json:"data"`
}

func fetchPrice(appId uint32) (string, error) {
	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%d&filters=price_overview", appId)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var obj map[string]json.RawMessage
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return "", err
	}

	var priceResp PriceResponse
	err = json.Unmarshal(obj[fmt.Sprint(appId)], &priceResp)
	if err != nil {
		return "", nil
	}

	return priceResp.Data.PriceOverview.Price, nil
}

func main() {
	var name string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, scanner.Err().Error())
		return
	}
	name = scanner.Text()

	price, err := getPrice(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fmt.Println("Price:", price)
}
