package epicGames

import (
	"context"
	"regexp"
	"strings"

	"web-scraper/common"

	"github.com/machinebox/graphql"
)

func GetInfo(title string) ([]common.GameInfo, error) {
	client := graphql.NewClient("https://graphql.epicgames.com/graphql")

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
	req.Var("country", "UA")
	req.Var("allowCountries", "UA")
	req.Var("locale", "ua-UA")
	req.Var("withPrice", true)
	req.Var("withMapping", true)
	req.Var("freeGame", false)
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
		return nil, err
	}

	regex := regexp.MustCompile("[^a-z0-9 ]+")
	title = regex.ReplaceAllString(strings.ToLower(title), "")

	games := make([]common.GameInfo, 0, len(response.Catalog.SearchStore.Elements))
	for _, elem := range response.Catalog.SearchStore.Elements {
		formatted := regex.ReplaceAllString(strings.ToLower(elem.Title), "")
		if strings.Contains(formatted, title) {
			games = append(games, common.GameInfo{Title: elem.Title, Price: elem.Price.TotalPrice.Formatted.DiscountPrice})
		}
	}
	return games, nil
}
