package epicGames

import (
	"context"
	"strings"

	"web-scraper/common"

	"github.com/machinebox/graphql"
)

type Scraper struct{}

const graphqlEndpoint = "https://graphql.epicgames.com/graphql"

func (Scraper) GetName() string {
	return "Epic Games Store"
}

func (Scraper) GetInfo(ch chan common.Result, id int, title string) {
	client := graphql.NewClient(graphqlEndpoint)

	query := `
	  query searchStoreQuery($allowCountries: String, $count: Int, $country: String!, $keywords: String, $locale: String, $sortBy: String, $sortDir: String, $withPrice: Boolean = false, $freeGame: Boolean) {
	    Catalog {
		  searchStore(
			allowCountries: $allowCountries
			count: $count
			country: $country
			keywords: $keywords
			locale: $locale
			sortBy: $sortBy
			sortDir: $sortDir
			freeGame: $freeGame
		  ) {
		    elements {
			  title
			  price(country: $country) @include(if: $withPrice) {
				totalPrice {
				  fmtPrice(locale: $locale) {
					originalPrice
					discountPrice
				  }
				}
			  }
			}	
		  }
		}
	  }`

	req := graphql.NewRequest(query)
	req.Var("keywords", title)
	req.Var("country", strings.ToUpper(common.CountryCode))
	req.Var("allowCountries", strings.ToUpper(common.CountryCode))
	req.Var("locale", common.Locale)
	req.Var("withPrice", true)
	req.Var("withMapping", true)
	req.Var("freeGame", false)
	req.Var("sortBy", "title")
	req.Var("sortDir", "asc")
	req.Var("count", common.MaxCount)

	type Response struct {
		Catalog struct {
			SearchStore struct {
				Elements []struct {
					Title string `json:"title"`
					Price struct {
						TotalPrice struct {
							Formatted struct {
								DiscountPrice string `json:"discountPrice"`
							} `json:"fmtPrice"`
						} `json:"totalPrice"`
					} `json:"price"`
				} `json:"elements"`
			} `json:"searchStore"`
		} `json:"catalog"`
	}

	var response Response
	err := client.Run(context.Background(), req, &response)
	if err != nil {
		ch <- common.Result{Id: id, Info: nil, Error: err}
	}

	title = common.AlphanumericRegex.ReplaceAllString(strings.ToLower(title), "")

	games := make([]common.GameInfo, 0, len(response.Catalog.SearchStore.Elements))
	for _, elem := range response.Catalog.SearchStore.Elements {
		formatted := common.AlphanumericRegex.ReplaceAllString(strings.ToLower(elem.Title), "")
		if strings.Contains(formatted, title) {
			games = append(games, common.GameInfo{Title: elem.Title, FormattedTitle: formatted, Price: elem.Price.TotalPrice.Formatted.DiscountPrice})
		}
	}
	ch <- common.Result{Id: id, Info: games, Error: nil}
}
