package epicGames

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func GetPrice(slug string) (string, error) {
	url := fmt.Sprintf("https://store.epicgames.com/p/%s", slug)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// TODO: Load the entire page with all JavaScipt code executed

	text := string(body)
	const sectionName = "priceSpecification\":{\"price\":\""
	startOfPrice := strings.Index(text, sectionName)
	if startOfPrice == -1 {
		return "", fmt.Errorf("Price not found on HTML page!")
	}
	endOfPrice := strings.Index(text[startOfPrice+len(sectionName):], ",")
	if endOfPrice == -1 {
		return "", fmt.Errorf("Price not found on HTML page!")
	}
	return text[startOfPrice+1 : endOfPrice], nil
}

func GetSlugs(name string) ([]string, error) {
	url := "https://store-content.ak.epicgames.com/api/content/productmapping"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	regex := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	name = regex.ReplaceAllString(name, "")
	formattedName := strings.Join(strings.Split(strings.ToLower(name), " "), "-")

	slugs := make([]string, 0)
	for _, slug := range data {
		if strings.Contains(slug, formattedName) {
			slugs = append(slugs, slug)
		}
	}

	return slugs, nil
}
