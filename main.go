package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	env "github.com/joho/godotenv"
)

func print_exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func resp_to_json(response http.Response) string {
	body, ioReadErr := io.ReadAll(response.Body)

	if ioReadErr != nil {
		panic(ioReadErr)
	}

	return string(body)
}

func set_env() map[string]string {
	envMap, envLoadErr := env.Read(".env")

	if envLoadErr != nil {
		panic(envLoadErr)
	}
	return envMap
}

func main() {

	envMap := set_env()

	API_KEY := envMap["API_KEY"]

	URL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=London&aqi=no", API_KEY)

	res, httpErr := http.Get(URL)

	if httpErr != nil {
		print_exit("Weather API not available")
	}

	if res.StatusCode != 200 {
		print_exit("Weather API not available")
	}

	json := resp_to_json(*res)

	println(json)
}
