## Game stores web scraper

A simple CLI web scraper for fetching latest prices of games from different online game stores (Epic Games Store, Steam and GOG at the moment)

Example of usage:

![Example of usage](doc/screenshot.jpg)
  
### Installation

The prebuilt executable can be downloaded from *Releases* section here on GitHub.

**Alternatively** you can build it manually:
- You need to have Go installed on your machine (visit official Go installation [guide](https://go.dev/doc/install)).
- Clone this repository and navigate to `src/`.
- Run `go build -o [your folder here]`. This will create executable in the specified folder.
***The application will generate some files, so if you don't want them to be scattered around, make a dedicated folder for the app!***
- Launch the built executable (`web-scraper.exe`) and enjoy!

### User manual
On the first launch of the application you will see the following message: `settings.json not found, loading fallback values`. This means that in the folder where the application is located a *settings.json* file was created. You can open this file with any text editor. It has the following structure:

    {
        "CountryCode": "ua",
        "Locale": "uk"
    }

These values can be changed and will be applied the next time you open the application.

`CountryCode` - **must** be a valid country code (refer to this [list](https://countrycode.org/), where first two letters of the third column are valid country codes).

`Locale` - **must** be a valid locale (refer to this [list](https://www.science.co.il/language/Locale-codes.php), where the third column contains valid locales).

Both `CountryCode` and `Locale` are case-insensitive. If country code and/or locale are invalid, `ua` and `uk` will be chosen as default (Glory to Ukraine!).

The application itself is really simple. Every time you open it, a prompt will ask you for a name of the game. Enter the name and press *Enter*. Fetching will take some time, and after it's done, the prices will be displayed. Press *Enter* to exit the application.

### Development progress
Currently the application has very limited functionality, but I have plans for developing it much further!

Planned features include:
- **Wishlist** - user can select specific games that they are interested in, the application will start tracking these games on stores and notify user about discounts
- **Game releases** - the app periodically gathers information from stores and notifies user about new releases
- **New stores** - the app will support more online game stores
- **GUI client** - GUI will provide better user experience (although CLI edition will still be available)

Check out my [Trello board](https://trello.com/b/0W9bt4xw/game-stores-scraper) to see the development progress.

If you found a bug or have an idea for new feature, then feel free to create an issue here on GitHub.