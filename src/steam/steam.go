package steam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"web-scraper/common"

	"github.com/bojanz/currency"
)

func GetInfo(title string) ([]common.GameInfo, error) {
	// I don't know why, but sometimes this endpoint returns different results
	// Or even throws "stream error: stream ID 1; INTERNAL_ERROR; received from peer"
	// As far as I tested this problem is not deterministic whatsoever
	// I guess the Steam server gets overloaded from time to time
	// TODO: Find a workaround
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

	if len(matches) == 0 {
		return make([]common.GameInfo, 0), nil
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

	games := make([]common.GameInfo, 0, len(matches))
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
				CurrencyCode string `json:"currency"`
				Price        int64  `json:"final"`
			} `json:"price_overview"`
		} `json:"data"`
	}

	locale := currency.NewLocale("uk") // TODO: Replace this with actual user locale
	formatter := currency.NewFormatter(locale)

	prices := make([]string, len(appIds))
	for i, appId := range appIds {
		var priceResponse PriceResponse
		err = json.Unmarshal(obj[appId], &priceResponse)
		if err != nil || priceResponse.Data.PriceOverview.Price == 0 {
			prices[i] = ""
			continue
		}
		amount, _ := currency.NewAmountFromInt64(priceResponse.Data.PriceOverview.Price, priceResponse.Data.PriceOverview.CurrencyCode)
		prices[i] = formatter.Format(amount)
	}

	return prices, nil
}
