package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	figure "github.com/common-nighthawk/go-figure"
	cli "gopkg.in/urfave/cli.v1"
)

// Main weather stats in kelvin
type Main struct {
	Temp     float32
	Pressure int
	Humidity int
	TempMin  float32
	TempMax  float32
}

// Wind stats
type Wind struct {
	Speed float32
	Def   int
}

// Clouds stats
type Clouds struct {
	All int
}

// Sys weather station stats
type Sys struct {
	Type    int `json:"type"`
	ID      int
	Message float32
	Country string
	Sunrise int
	Sunset  int
}

// Weather description
type Weather struct {
	ID          int
	Main        string
	Description string
	Icon        string
}

type Coord struct {
	Lon float32
	Lat float32
}

// CityStats main obj
type CityStats struct {
	Coord      Coord `json:"coord"`
	Weather    []Weather
	Base       string
	Main       Main
	Visibility int
	Wind       Wind
	Clouds     Clouds
	dt         int
	Sys        Sys `json:"sys"`
	ID         int `json:"id"`
	Name       string
	Cod        int `json:"cod"`
}

const (
	absoluteZeroC = 273.15
	absoluteZeroF = 459.67
)

// KelvinToFahrenheit converts kelvin to Farenheight
func KelvinToFahrenheit(k float32) float32 {
	value := math.Round(float64((k * 9 / 5) - absoluteZeroF))
	return float32(value)
}

func main() {

	var cityParam string

	app := cli.NewApp()
	app.Name = "Minimal weather"
	app.Usage = "Get current weather for a city"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "city",
			Value:       "San Francisco, US",
			Destination: &cityParam,
		},
	}

	app.Action = func(ctx *cli.Context) error {
		// On ^C, handle exit
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
			for sig := range c {
				fmt.Printf("Received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

		apiToken := os.getEnv("API_TOKEN")

		weatherAPI := "https://api.openweathermap.org/data/2.5/weather?APPID=" + apiToken + "&q=" + cityParam

		weatherClient := http.Client{
			Timeout: time.Second * 2,
		}

		req, err := http.NewRequest(http.MethodGet, weatherAPI, nil)
		if err != nil {
			log.Fatal(err)
		}

		res, err := weatherClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		jsonBody := CityStats{}
		if err := json.Unmarshal(body, &jsonBody); err != nil {
			panic(err)
		}

		weatherIntoASCII := figure.NewFigure(jsonBody.Name, "doom", true)
		weatherIntoASCII.Print()

		temperatureInF := KelvinToFahrenheit(jsonBody.Main.Temp)
		fmt.Printf("\n\nCurrent: %.2fF \n", temperatureInF)
		fmt.Printf("Wind: %0.1fkph \n", jsonBody.Wind.Speed)
		fmt.Printf("Clouds: %s \n", jsonBody.Weather[0].Description)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
