package steam

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

type appListResponse struct {
	List struct {
		Apps []struct {
			AppId uint32 `json:"appid"`
			Name  string `json:"name"`
		} `json:"apps"`
	} `json:"applist"`
}

func GetPrice(name string) (string, error) {
	resp, err := http.Get("https://api.steampowered.com/ISteamApps/GetAppList/v2/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var appList appListResponse
	err = json.Unmarshal(body, &appList)
	if err != nil {
		return "", err
	}

	matches := make([]string, 0)
	appIds := make([]string, 0)
	for _, elem := range appList.List.Apps {
		if len(appIds) >= 100 {
			break
		}
		if strings.Contains(strings.ToLower(elem.Name), strings.ToLower(name)) {
			matches = append(matches, elem.Name)
			appIds = append(appIds, fmt.Sprint(elem.AppId))
		}
	}

	if len(appIds) == 0 {
		return "", fmt.Errorf("No matches found! Try more specific name.")
	}

	prices, err := fetchPrices(appIds)
	if err != nil {
		return "", err
	}

	type app struct {
		name  string
		price string
	}

	apps := make([]app, 0, len(matches))
	for j, price := range prices {
		if price != "" {
			apps = append(apps, app{matches[j], price})
		}
	}

	// Sort the matches by name.
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].name < apps[j].name
	})

	// Display the closest matches to the user's input.
	fmt.Println("Closest matches:")
	for i, m := range apps {
		fmt.Printf("%d. %s\n", i+1, m.name)
	}
	fmt.Println("\nIf the game that you are looking for is not on the list,",
		"then please try again and be more specific next time.")

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

		// Find the price for the selected game.
		for _, app := range apps {
			if app.name == selectedGame {
				return app.price, nil
			}
		}
		fmt.Println("No such game in the list! Please, try again.")
	}
}

type priceResponse struct {
	Data struct {
		PriceOverview struct {
			Price string `json:"final_formatted"`
		} `json:"price_overview"`
	} `json:"data"`
}

func fetchPrices(appIds []string) ([]string, error) {
	param := strings.Join(appIds, ",")
	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s&filters=price_overview", param)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var obj map[string]json.RawMessage
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return nil, err
	}

	prices := make([]string, len(appIds))
	for i, appId := range appIds {
		var priceResp priceResponse
		err = json.Unmarshal(obj[appId], &priceResp)
		if err != nil {
			prices[i] = ""
			continue
		}
		prices[i] = priceResp.Data.PriceOverview.Price
	}

	return prices, nil
}
