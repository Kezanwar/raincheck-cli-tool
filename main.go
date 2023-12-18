package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	_ "github.com/joho/godotenv/autoload"
	"github.com/savioxavier/termlink"
)

func json_stringify(not_json interface{}) string {
	data, err := json.Marshal(not_json)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func fail_gracefully(msg string) {
	color.Red(msg)
	os.Exit(1)
}

func res_to_weather(response http.Response) Weather {
	body, ioReadErr := io.ReadAll(response.Body)

	if ioReadErr != nil {
		panic(ioReadErr)
	}

	var weather Weather

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

func print_title(location Location, current Current) {

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

func handle_trailing_numbers(num float64) string {
	text := fmt.Sprint(int(math.Floor(num)))
	return text
}

func get_spacing(len int) string {
	var spacing []string
	for i := 1; i <= len; i++ {
		spacing = append(spacing, " ")
	}
	return strings.Join(spacing, "")
}

func print_hours(hours []Hour) {
	const maxDegree = 2
	const maxChance = 3
	for i, hour := range hours {
		date := time.Unix(int64(hour.TimeEpoch), 0)

		if date.Before(time.Now()) {
			continue
		}

		var temp = handle_trailing_numbers(hour.TempC)

		var chance_of_rain = handle_trailing_numbers(hour.ChanceOfRain)

		spaceDegree := get_spacing(maxDegree - len(temp))
		spaceChance := get_spacing(maxChance - len(chance_of_rain))

		message := fmt.Sprintf("%s             %s°  %s               %s%% %s           %s\n", date.Format("15:04"), temp, spaceDegree, chance_of_rain, spaceChance, hour.Condition.Text)

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

func get_search_term() string {
	var search string
	if len(os.Args) < 2 {

		response, httpErr := http.Get("http://ip-api.com/json/?fields=city")

		if httpErr != nil || response.StatusCode != 200 {
			fail_gracefully("No results found, try searching for a city instead, e.g $ raincheck <city>")
		}

		defer response.Body.Close()

		body, ioReadErr := io.ReadAll(response.Body)

		if ioReadErr != nil {
			panic(ioReadErr)
		}

		var city IPCity

		marshalError := json.Unmarshal(body, &city)

		if marshalError != nil {
			panic(marshalError)
		}

		search = city.City

	} else {
		search = os.Args[1]
	}
	return search
}

var DEV_API_KEY = os.Getenv("API_KEY")
var BUILD_API_KEY string

func main() {

	var search = get_search_term()

	var API_KEY string

	if len(BUILD_API_KEY) > 0 {
		API_KEY = BUILD_API_KEY
	} else {
		API_KEY = DEV_API_KEY
	}

	URL := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&aqi=no", API_KEY, search)

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
