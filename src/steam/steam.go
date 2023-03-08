package steam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"web-scraper/common"
)

func GetInfo(title string) ([]common.GameInfo, error) {
	resp, err := http.Get("https://api.steampowered.com/ISteamApps/GetAppList/v2/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type AppListResponse struct {
		List struct {
			Apps []struct {
				AppId uint32 `json:"appid"`
				Name  string `json:"name"`
			} `json:"apps"`
		} `json:"applist"`
	}

	var appList AppListResponse
	err = json.Unmarshal(body, &appList)
	if err != nil {
		return nil, err
	}

	title = common.Regex.ReplaceAllString(strings.ToLower(title), "")

	type Match struct {
		Title          string
		FormattedTitle string
		AppId          string
	}

	matches := make([]Match, 0)
	for _, elem := range appList.List.Apps {
		if len(matches) == common.MaxCount {
			break
		}
		formatted := common.Regex.ReplaceAllString(strings.ToLower(elem.Name), "")
		if strings.Contains(formatted, title) {
			matches = append(matches, Match{Title: elem.Name, FormattedTitle: formatted, AppId: fmt.Sprint(elem.AppId)})
		}
	}

	games := make([]common.GameInfo, 0, len(matches))
	if len(matches) == 0 {
		return games, nil
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Title < matches[j].Title
	})

	appIds := make([]string, 0, len(matches))
	for _, m := range matches {
		appIds = append(appIds, m.AppId)
	}
	prices, err := fetchPrices(appIds)
	if err != nil {
		return nil, err
	}

	for i, price := range prices {
		if price != "" {
			games = append(games, common.GameInfo{Title: matches[i].Title, FormattedTitle: matches[i].FormattedTitle, Price: price})
		}
	}

	return games, nil
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

	type PriceResponse struct {
		Data struct {
			PriceOverview struct {
				Price string `json:"final_formatted"`
			} `json:"price_overview"`
		} `json:"data"`
	}

	prices := make([]string, len(appIds))
	for i, appId := range appIds {
		var priceResponse PriceResponse
		err = json.Unmarshal(obj[appId], &priceResponse)
		if err != nil {
			prices[i] = ""
			continue
		}
		prices[i] = priceResponse.Data.PriceOverview.Price
	}

	return prices, nil
}
