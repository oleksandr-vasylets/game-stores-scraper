package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

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

func getAppId(name string) (uint32, error) {
	resp, err := http.Get("http://api.steampowered.com/ISteamApps/GetAppList/v2/")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var appList AppListResponse
	err = json.Unmarshal(body, &appList)
	if err != nil {
		return 0, err
	}

	type app struct {
		id   uint32
		dist int
	}

	// Find the closest matches to the user's input using Levenshtein distance.
	matches := make(map[string]app)
	for _, elem := range appList.List.Apps {
		dist := levenshtein.DistanceForStrings([]rune(elem.Name), []rune(name), levenshtein.DefaultOptions)
		if dist <= len(name)/2 {
			matches[elem.Name] = app{elem.AppId, dist}
		}
	}

	if len(matches) == 0 {
		return 0, fmt.Errorf("Game not found!")
	}

	// Sort the matches by distance.
	type match struct {
		name string
		id   uint32
		dist int
	}
	var sortedMatches []match
	for name, elem := range matches {
		sortedMatches = append(sortedMatches, match{name, elem.id, elem.dist})
	}
	sort.Slice(sortedMatches, func(i, j int) bool {
		return sortedMatches[i].dist < sortedMatches[j].dist
	})

	// Display the closest matches to the user's input.
	fmt.Println("Closest matches:")
	for i, m := range sortedMatches {
		fmt.Printf("%d. %s (distance %d)\n", i+1, m.name, m.dist)
	}

	// Ask the user to select a game from the closest matches.
	for {
		var selectedGame string
		fmt.Print("Enter the name of a game from the list above: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Fprintln(os.Stderr, scanner.Err().Error())
			return 0, scanner.Err()
		}
		selectedGame = scanner.Text()

		// Find the appid for the selected game.
		for _, app := range sortedMatches {
			if app.name == selectedGame {
				return app.id, nil
			}
		}
		fmt.Println("No such game in the list! Please, try again")
	}
}

type PriceResponse struct {
	Data struct {
		PriceOverview struct {
			Price string `json:"final_formatted"`
		} `json:"price_overview"`
	} `json:"data"`
}

func getPrice(appId string) (string, error) {
	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s&filters=price_overview", appId)
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

	var data PriceResponse
	err = json.Unmarshal(obj[appId], &data)
	if err != nil {
		return "", err
	}
	return data.Data.PriceOverview.Price, nil
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

	appId, err := getAppId(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	fmt.Println(appId)
	price, err := getPrice(fmt.Sprint(appId))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	fmt.Printf(price)
}
