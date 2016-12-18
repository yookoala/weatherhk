package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yookoala/weatherhk/hkodata"
)

var port int

func init() {
	portStr := os.Getenv("PORT")
	if portStr != "" {
		portInt, _ := strconv.ParseUint(portStr, 10, 32)
		port = int(portInt)
	}
	if port == 0 {
		port = 8080 // fallback value
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `<html><h1>Simple Hong Kong Weather API</h1><a href="/api/CurrentWeather.json">Current Weather</a></html>`)
	})

	r.HandleFunc("/api/CurrentWeather.json", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		requestID := r.Header.Get("X-Request-ID")

		req, err := http.Get("http://rss.weather.gov.hk/rss/CurrentWeather.xml")
		if err != nil {
			log.Printf("request-id: %s, err: %s", requestID, err.Error())
			return
		}

		// prepare encoder for output
		enc := json.NewEncoder(w)

		// decode the RSS
		data, err := hkodata.DecodeCurrentWeather(req.Body)
		if err != nil {
			log.Printf("request-id: %s, err: %s", requestID, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(struct {
				Status  int    `json:"status"`
				Message string `json:"message"`
			}{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		// return formatted data

		// RFC2616: Tue, 15 Nov 1994 12:45:26 GMT
		// RFC1123: Mon, 02 Jan 2006 15:04:05 MST
		w.Header().Set("Last-Modified", data.PubDate.Format(time.RFC1123))
		w.WriteHeader(http.StatusOK)
		enc.Encode(struct {
			Status int                    `json:"status"`
			Data   hkodata.CurrentWeather `json:"data"`
			// Raw    string                 `json:"raw_data"`
		}{
			Status: http.StatusOK,
			Data:   *data,
			// Raw:    data.Raw,
		})
	})

	log.Printf("listen at port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
