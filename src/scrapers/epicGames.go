package scrapers

import (
	"context"
	"game-stores-scraper/settings"
	"strings"

	"github.com/machinebox/graphql"
)

type EpicGamesScraper struct{}

func (EpicGamesScraper) GetName() string {
	return "Epic Games Store"
}

func (EpicGamesScraper) GetInfo(ch chan Result, id int, title string) {
	const graphqlEndpoint = "https://graphql.epicgames.com/graphql"
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
	req.Var("country", strings.ToUpper(settings.CountryCode()))
	req.Var("allowCountries", strings.ToUpper(settings.CountryCode()))
	req.Var("locale", settings.Locale())
	req.Var("withPrice", true)
	req.Var("withMapping", true)
	req.Var("freeGame", false)
	req.Var("sortBy", "title")
	req.Var("sortDir", "asc")
	req.Var("count", settings.MaxCount())

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
		ch <- Result{Id: id, Info: nil, Error: err}
		return
	}

	title = alphanumericRegex.ReplaceAllString(strings.ToLower(title), "")

	games := make([]GameInfo, 0, len(response.Catalog.SearchStore.Elements))
	for _, elem := range response.Catalog.SearchStore.Elements {
		formatted := alphanumericRegex.ReplaceAllString(strings.ToLower(elem.Title), "")
		if strings.Contains(formatted, title) {
			games = append(games, GameInfo{Title: elem.Title, FormattedTitle: formatted, Price: elem.Price.TotalPrice.Formatted.DiscountPrice})
		}
	}
	ch <- Result{Id: id, Info: games, Error: nil}
}
