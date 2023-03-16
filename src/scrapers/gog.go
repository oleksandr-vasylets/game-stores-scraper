package scrapers

import (
	"encoding/json"
	"fmt"
	"game-stores-scraper/settings"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bojanz/currency"
)

type GogScraper struct{}

func (GogScraper) GetName() string {
	return "GoG"
}

func (scraper GogScraper) GetInfo(ch chan Result, id int, title string) {
	const gameListQuery = "https://www.gog.com/games/ajax/filtered?mediaType=game&limit%d&search=%s"
	url := fmt.Sprintf(gameListQuery, settings.MaxCount, url.QueryEscape(title))
	resp, err := http.Get(url)
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
		Products []struct {
			Id    int64  `json:"id"`
			Title string `json:"title"`
			Price struct {
				IsFree bool `json:"isFree"`
			} `json:"price"`
			Buyable bool `json:"buyable"`
		} `json:"products"`
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		ch <- Result{Id: id, Info: nil, Error: err}
		return
	}

	games := make([]GameInfo, 0)
	ids := make([]int64, 0, len(response.Products))
	for _, game := range response.Products {
		if game.Buyable && !game.Price.IsFree {
			ids = append(ids, game.Id)
			formatted := alphanumericRegex.ReplaceAllString(strings.ToLower(game.Title), "")
			games = append(games, GameInfo{Title: game.Title, FormattedTitle: formatted})
		}
	}

	if len(games) == 0 {
		ch <- Result{Id: id, Info: make([]GameInfo, 0), Error: err}
		return
	}

	prices, err := scraper.fetchPrices(ids)
	if err != nil {
		ch <- Result{Id: id, Info: nil, Error: err}
		return
	}

	for i, price := range prices {
		games[i].Price = price
	}

	ch <- Result{Id: id, Info: games, Error: err}
	return
}

func (GogScraper) fetchPrices(ids []int64) ([]string, error) {
	locale := currency.NewLocale(settings.UserProfile.Locale)
	formatter := currency.NewFormatter(locale)
	const priceQuery = "https://api.gog.com/products/%d/prices?countryCode=%s"

	prices := make([]string, 0, len(ids))
	for _, id := range ids {
		url := fmt.Sprintf(priceQuery, id, settings.UserProfile.CountryCode)
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

		if len(priceResponse.Data.Prices) == 0 {
			fmt.Println(id)
		}

		tokens := strings.Split(priceResponse.Data.Prices[0].Price, " ")
		price, _ := strconv.ParseInt(tokens[0], 10, 64)
		currencyCode := tokens[1]
		amount, _ := currency.NewAmountFromInt64(price, currencyCode)
		prices = append(prices, formatter.Format(amount))
	}

	return prices, nil
}
