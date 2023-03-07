package main

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

func main() {
	// fmt.Print("Enter the name of the game: ")
	// var name string
	// scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan()
	// if scanner.Err() != nil {
	// 	fmt.Fprintln(os.Stderr, scanner.Err().Error())
	// 	return
	// }
	// name = scanner.Text()

	// price, err := epicGames.GetPrice(name)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(price)

	client := graphql.NewClient("https://graphql.epicgames.com/graphql")

	// Define the GraphQL query
	query := `
		query searchStoreQuery($allowCountries: String, $category: String, $count: Int, $country: String!, $keywords: String, $locale: String, $namespace: String, $itemNs: String, $sortBy: String, $sortDir: String, $start: Int, $tag: String, $releaseDate: String, $withPrice: Boolean = false, $withPromotions: Boolean = false, $priceRange: String, $freeGame: Boolean, $onSale: Boolean, $effectiveDate: String) {
			Catalog {
			searchStore(
				allowCountries: $allowCountries
				category: $category
				count: $count
				country: $country
				keywords: $keywords
				locale: $locale
				namespace: $namespace
				itemNs: $itemNs
				sortBy: $sortBy
				sortDir: $sortDir
				releaseDate: $releaseDate
				start: $start
				tag: $tag
				priceRange: $priceRange
				freeGame: $freeGame
				onSale: $onSale
				effectiveDate: $effectiveDate
			) {
				elements {
				title
				id
				namespace
				description
				effectiveDate
				keyImages {
					type
					url
				}
				currentPrice
				seller {
					id
					name
				}
				productSlug
				urlSlug
				url
				tags {
					id
				}
				items {
					id
					namespace
				}
				customAttributes {
					key
					value
				}
				categories {
					path
				}
				catalogNs {
					mappings(pageType: "productHome") {
					pageSlug
					pageType
					}
				}
				offerMappings {
					pageSlug
					pageType
				}
				price(country: $country) @include(if: $withPrice) {
					totalPrice {
					discountPrice
					originalPrice
					voucherDiscount
					discount
					currencyCode
					currencyInfo {
						decimals
					}
					fmtPrice(locale: $locale) {
						originalPrice
						discountPrice
						intermediatePrice
					}
					}
					lineOffers {
					appliedRules {
						id
						endDate
						discountSetting {
						discountType
						}
					}
					}
				}
				promotions(category: $category) @include(if: $withPromotions) {
					promotionalOffers {
					promotionalOffers {
						startDate
						endDate
						discountSetting {
						discountType
						discountPercentage
						}
					}
					}
					upcomingPromotionalOffers {
					promotionalOffers {
						startDate
						endDate
						discountSetting {
						discountType
						discountPercentage
						}
					}
					}
				}
				}
				paging {
				count
				total
				}
			}
			}
		}
	`

	// Define the GraphQL variables
	variables := map[string]interface{}{
		"category":       "",
		"keywords":       "assassin",
		"country":        "UA",
		"allowCountries": "UA",
		"locale":         "en-US",
		"sortDir":        "desc",
		"withPrice":      true,
		"withMapping":    true,
	}

	// Create a new request object with the query and variables
	req := graphql.NewRequest(query)
	req.Var("category", variables["category"])
	req.Var("keywords", variables["keywords"])
	req.Var("country", variables["country"])
	req.Var("allowCountries", variables["allowCountries"])
	req.Var("locale", variables["locale"])
	req.Var("sortDir", variables["sortDir"])
	req.Var("withPrice", variables["withPrice"])
	req.Var("withMapping", variables["withMapping"])
	req.Var("freeGame", false)

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
		panic(err)
	}

	for _, elem := range response.Catalog.SearchStore.Elements {
		fmt.Println("Title:", elem.Title)
		fmt.Println("Price:", elem.Price.TotalPrice.Formatted.DiscountPrice)
		fmt.Println()
	}
}
