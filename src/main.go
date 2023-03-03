package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Data struct {
		PriceOverview struct {
			Price string `json:"final_formatted"`
		} `json:"price_overview"`
	} `json:"data"`
}

func main() {
	appid := "105600" // Half-Life 2

	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s&filters=price_overview", appid)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var obj map[string]json.RawMessage
	err = json.Unmarshal(body, &obj)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data Response
	err = json.Unmarshal(obj[appid], &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(data.Data.PriceOverview.Price)
}
