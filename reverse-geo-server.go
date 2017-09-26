package main

import (
	"github.com/codingsince1985/geo-golang/openstreetmap"
//	"github.com/codingsince1985/geo-golang/mapquest/open"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/magiconair/properties"
	"net/http"
	"strconv"
	"strings"
)

var osmGeocoder = openstreetmap.Geocoder()

func reverseGeo(c echo.Context, lon float64, lat float64, ch chan myHandler) {
	//Reverse Geocoding
	address, err := osmGeocoder.ReverseGeocode(lon, lat)
	s := ""
	if err == nil {
		s = address.FormattedAddress
		c.Logger().Info(s)
	} else {
		s = err.Error()
		c.Logger().Error(s)
	}
	ch <- myHandler{s, err}
}



//Health checking endpoint (similar springboot actuator)
func health(c echo.Context) error {
	r := &Ret{
		Status: "Up",
	}
	//return c.String(http.StatusOK, "{\"Status\":\"Up\"}")
	return c.JSON(http.StatusOK, r)
}

func reverse(c echo.Context) error {

	lon, err := strconv.ParseFloat(c.QueryParam("lon"), 64)
	if err != nil {
		return c.String(http.StatusBadRequest, c.QueryString())
	}
	lat, err := strconv.ParseFloat(c.QueryParam("lat"), 64)
	if err != nil {
		return c.String(http.StatusBadRequest, c.QueryString())
	}

	format := c.QueryParam("format")
	if format == "" {
		format = ""
	}

	status := http.StatusBadRequest
	ch := make(chan myHandler)
	go reverseGeo(c, lon, lat, ch)
	handler := <-ch
	s := handler.s
	if handler.e == nil {
		status = http.StatusOK
	}
	if strings.ToUpper(format) == "JSON" {
		return c.JSON(status, s)
	}
	if strings.ToUpper(format) == "XML" {
		return c.XML(status, s)
	}
	return c.String(status, s)
}

func main() {
	port := ":8080"
	// init from a file
	p, err := properties.LoadFile("./config.properties", properties.UTF8)
	if err == nil {
		port = ":" + strconv.Itoa(p.GetInt("server.port", 8080))
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", health)
	e.GET("/reverse", reverse)
	e.Logger.Fatal(e.Start(port))

}
