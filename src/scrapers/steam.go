package scrapers

import (
	"encoding/json"
	"fmt"
	"game-stores-scraper/settings"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/bojanz/currency"
)

type SteamScraper struct{}

func (SteamScraper) GetName() string {
	return "Steam"
}

func (scraper SteamScraper) GetInfo(ch chan Result, id int, title string) {
	// Sometimes accessing this endpoint throws "stream error: stream ID 1; INTERNAL_ERROR; received from peer"
	// I guess the Steam server gets overloaded from time to time
	// TODO: Find a workaround
	const appListEndpoint = "https://api.steampowered.com/ISteamApps/GetAppList/v2/"
	resp, err := http.Get(appListEndpoint)
	if err != nil {
		ch <- Result{Id: id, Info: nil, Error: err}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- Result{Id: id, Info: nil, Error: err}
		return
	}

	type Response struct {
		List struct {
			Apps []struct {
				AppId uint32 `json:"appid"`
				Name  string `json:"name"`
			} `json:"apps"`
		} `json:"applist"`
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		ch <- Result{Id: id, Info: nil, Error: err}
		return
	}

	title = alphanumericRegex.ReplaceAllString(strings.ToLower(title), "")

	type Match struct {
		Title          string
		FormattedTitle string
		AppId          string
	}

	matches := make([]Match, 0, settings.MaxCount)
	for _, elem := range response.List.Apps {
		if len(matches) == settings.MaxCount {
			break
		}
		formatted := alphanumericRegex.ReplaceAllString(strings.ToLower(elem.Name), "")
		if strings.Contains(formatted, title) {
			matches = append(matches, Match{Title: elem.Name, FormattedTitle: formatted, AppId: fmt.Sprint(elem.AppId)})
		}
	}

	if len(matches) == 0 {
		ch <- Result{Id: id, Info: make([]GameInfo, 0), Error: nil}
		return
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Title < matches[j].Title
	})

	appIds := make([]string, 0, len(matches))
	for _, m := range matches {
		appIds = append(appIds, m.AppId)
	}
	prices, err := scraper.fetchPrices(appIds)
	if err != nil {
		ch <- Result{Id: id, Info: nil, Error: err}
		return
	}

	games := make([]GameInfo, 0, len(matches))
	for i, price := range prices {
		if price != "" {
			games = append(games, GameInfo{Title: matches[i].Title, FormattedTitle: matches[i].FormattedTitle, Price: price})
		}
	}

	ch <- Result{Id: id, Info: games, Error: err}
}

func (SteamScraper) fetchPrices(appIds []string) ([]string, error) {
	ids := strings.Join(appIds, ",")
	const priceQuery = "https://store.steampowered.com/api/appdetails?appids=%s&filters=price_overview&cc=%s"
	url := fmt.Sprintf(priceQuery, ids, settings.UserProfile.CountryCode)
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

	locale := currency.NewLocale(settings.UserProfile.Locale)
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
