package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/log"
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

	service, err := selenium.NewChromeDriverService("..\\bin\\chromedriver.exe", 4444)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"--no-sandbox",
		"--disable-gpu",
		"--headless",
		"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
		"--log-level=3",
	}})
	caps.SetLogLevel(log.Driver, log.Off)
	caps.SetLogLevel(log.Browser, log.Off)

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		panic(err)
	}
	defer driver.Quit()

	driver.SetImplicitWaitTimeout(100 * time.Millisecond)
	driver.Get("https://store.epicgames.com/en-US/p/assassins-creed-1")
	div, err := driver.FindElement(selenium.ByID, "_schemaOrgMarkup-Product")
	if err != nil {
		panic(err)
	}
	script, err := div.GetAttribute("innerHTML")
	if err != nil {
		panic(err)
	}

	type Data struct {
		Offers []struct {
			PriceSpecification struct {
				Price         uint
				PriceCurrency string
			}
		}
	}

	var data Data
	err = json.Unmarshal([]byte(script), &data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d%s", data.Offers[0].PriceSpecification.Price, data.Offers[0].PriceSpecification.PriceCurrency)
}
