package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	types "raincheck-cli-tool/types"
	"strings"
	"time"

	"github.com/fatih/color"
	env "github.com/joho/godotenv"
	"github.com/savioxavier/termlink"
)

func fail_gracefully(msg string) {
	color.Red(msg)
	os.Exit(1)
}

func get_env() map[string]string {
	envMap, envLoadErr := env.Read(".env")

	if envLoadErr != nil {
		panic(envLoadErr)
	}

	return envMap
}

func res_to_weather(response http.Response) types.Weather {
	body, ioReadErr := io.ReadAll(response.Body)

	if ioReadErr != nil {
		panic(ioReadErr)
	}

	var weather types.Weather

	marshalError := json.Unmarshal(body, &weather)

	if marshalError != nil {
		panic(marshalError)
	}

	return weather
}

func print_blank() {
	fmt.Println("")
}

func print_logo() {
	color.Cyan(`  _______ _(_)__  ____/ /  ___ ____/ /__
 / __/ _ '/ / _ \/ __/ _ \/ -_) __/  '_/
/_/  \_,_/_/_//_/\__/_//_/\__/\__/_/\_\
`)

	print_blank()

	fmt.Print("🔗 ")
	color.Cyan(termlink.Link("github/kezanwar/raincheck-cli-tool", "https://www.github.com/kezanwar/raincheck-cli-tool"))
}

func print_title(location types.Location, current types.Current) {

	title := fmt.Sprintf("%s, %s: %.0f° %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	print_blank()
	print_blank()

	print_logo()

	print_blank()

	color.Magenta(title)

	print_blank()

	color.Magenta("🕛 Time        🌡️  Temp       🌧️  Chance Of Rain        ⛅️ Condition")

	print_blank()
}

func handle_trailing_numbers(num float64, suffix string) string {

	text := fmt.Sprint(int(math.Floor(num)))
	return text + suffix

}

func print_hours(hours []types.Hour) {
	for i, hour := range hours {

		date := time.Unix(int64(hour.TimeEpoch), 0)

		if date.Before(time.Now()) {
			continue
		}

		var temp = handle_trailing_numbers(hour.TempC, "°")
		var chane_of_rain = handle_trailing_numbers(hour.ChanceOfRain, "%")

		var spacing []string
		for range strings.Split(temp+chane_of_rain, "") {
			spacing = append(spacing, " ")
		}

		message := fmt.Sprintf("%s             %s                 %s               %s\n", date.Format("15:04"), temp, chane_of_rain, hour.Condition.Text)

		if hour.ChanceOfRain > 40 {
			color.Red(message)
		} else {
			color.Green(message)
		}

		if i != 23 {
			color.White("----------------------------------------------------------------------")
		} else {
			print_blank()

		}

	}
}

func main() {
	envMap := get_env()

	API_KEY := envMap["API_KEY"]

	URL := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=Manchester&aqi=no", API_KEY)

	response, httpErr := http.Get(URL)

	if httpErr != nil {
		fail_gracefully("Weather API not available")
	}

	if response.StatusCode != 200 {
		fail_gracefully("No results found")
	}

	defer response.Body.Close()

	weather := res_to_weather(*response)

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	print_title(location, current)

	print_hours(hours)

}
