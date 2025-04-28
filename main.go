package main

import (
	"encoding/json"
	"fmt"
	v "forecast/views"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Coord struct {
	Lat float64
	Lon float64
}

var CITIES = map[string]Coord{
	"Moscow":           Coord{Lat: 55.75222, Lon: 37.61556},
	"Saint-Petersburg": Coord{Lat: 59.9386, Lon: 30.3141},
	"Vladimir":         Coord{Lat: 56.1366, Lon: 40.3966},
	"Sochi":            Coord{Lat: 43.5992, Lon: 39.7257},
	"Anapa":            Coord{Lat: 39.7257, Lon: 37.3239},
	"Ulyanovsk":        Coord{Lat: 54.3282, Lon: 48.3866},
	"Astrakhan":        Coord{Lat: 46.35, Lon: 48.04},
	"Vladivostok":      Coord{Lat: 43.11, Lon: 131.87},
	"Krasnodar":        Coord{Lat: 45.04, Lon: 38.98},
	"Rostov-on-Don":    Coord{Lat: 47.23, Lon: 39.72},
	"Ekaterinburg":     Coord{Lat: 56.85, Lon: 60.61},
	"Volgograd":        Coord{Lat: 48.72, Lon: 44.5},
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())
	e.Static("/static", "static")
	e.GET("/", IndexPage)
	e.Debug = true
	e.Logger.Fatal(e.Start(":8080"))
}

func IndexPage(c echo.Context) error {
	res := make([]v.Forecast, 0)
	for k, v := range CITIES {
		f, err := GetWeatherData(v)
		if err != nil {
			return err
		}
		f.City = k
		res = append(res, *f)

	}
	//return c.JSON(200, res)
	return v.Index(res).Render(c.Request().Context(), c.Response().Writer)
}

func GetWeatherData(c Coord) (*v.Forecast, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&daily=weather_code,temperature_2m_max,temperature_2m_min&current=temperature_2m,weather_code&timeformat=unixtime", c.Lat, c.Lon))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}
	weather := new(v.Forecast)
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, err
	}
	return weather, nil
}
