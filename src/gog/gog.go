package gog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"web-scraper/common"

	"github.com/bojanz/currency"
)

func GetInfo(title string) ([]common.GameInfo, error) {
	url := fmt.Sprintf("https://www.gog.com/games/ajax/filtered?mediaType=game&limit%d&search=%s", common.MaxCount, url.QueryEscape(title))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Products []struct {
			Id    int64  `json:"id"`
			Title string `json:"title"`
			Price struct {
				IsFree bool `json:"isFree"`
			} `json:"price"`
		} `json:"products"`
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	games := make([]common.GameInfo, 0)
	ids := make([]int64, 0, len(response.Products))
	for _, game := range response.Products {
		if !game.Price.IsFree {
			ids = append(ids, game.Id)
			formatted := common.Regex.ReplaceAllString(strings.ToLower(game.Title), "")
			games = append(games, common.GameInfo{Title: game.Title, FormattedTitle: formatted})
		}
	}

	if len(games) == 0 {
		return make([]common.GameInfo, 0), nil
	}

	prices, err := fetchPrices(ids)
	if err != nil {
		return nil, err
	}

	for i, price := range prices {
		games[i].Price = price
	}

	return games, nil
}

func fetchPrices(ids []int64) ([]string, error) {
	locale := currency.NewLocale(common.Locale)
	formatter := currency.NewFormatter(locale)

	prices := make([]string, 0, len(ids))
	for _, id := range ids {
		url := fmt.Sprintf("https://api.gog.com/products/%d/prices?countryCode=%s", id, common.CountryCode)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		type Response struct {
			Data struct {
				Prices []struct {
					Price string `json:"finalPrice"`
				} `json:"prices"`
			} `json:"_embedded"`
		}

		var priceResponse Response
		err = json.Unmarshal(body, &priceResponse)
		if err != nil {
			return nil, err
		}

		tokens := strings.Split(priceResponse.Data.Prices[0].Price, " ")
		price, _ := strconv.ParseInt(tokens[0], 10, 64)
		currencyCode := tokens[1]
		amount, _ := currency.NewAmountFromInt64(price, currencyCode)
		prices = append(prices, formatter.Format(amount))
	}

	return prices, nil
}
