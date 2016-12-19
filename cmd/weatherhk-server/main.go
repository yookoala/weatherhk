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
var forceHTTPS bool
var hostname string

const noticeNonpublicAPI = "This source is not publicly announced by HKO. That means it can break without previous notice."

func init() {
	portStr := os.Getenv("PORT")
	if portStr != "" {
		portInt, _ := strconv.ParseUint(portStr, 10, 32)
		port = int(portInt)
	}
	if port == 0 {
		port = 8080 // fallback value
	}

	forceHTTPSStr := os.Getenv("FORCE_HTTPS")
	forceHTTPS = (forceHTTPSStr == "TRUE")

	hostname = os.Getenv("APP_HOSTNAME")
	if hostname == "" {
		hostname = "localhost"
	}

}

// enforceHTTPS is a simple middleware to enforce HTTPS on demand
func enforceHTTPS(forceHTTPS bool, inner http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if forceHTTPS && r.Header.Get("X-Forwarded-Proto") != "https" {
			redirectURL := *r.URL
			redirectURL.Host = hostname
			redirectURL.Scheme = "https"
			log.Printf("run here: %s", redirectURL)

			w.Header().Set("Location", redirectURL.String())
			w.WriteHeader(http.StatusMovedPermanently)
			return
		}

		// pass through to inner handler
		inner.ServeHTTP(w, r)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `<html><h1>Simple Hong Kong Weather API</h1><ul><li><a href="/api/CurrentWeather.json">Current Weather</a></li><li><a href="/api/region.json">Region Weather</a></li></ul></html>`)
	})

	r.HandleFunc("/api/CurrentWeather.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

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
		w.Header().Set("Expires", data.Expires().Format(time.RFC1123))
		w.WriteHeader(http.StatusOK)
		enc.Encode(struct {
			Status int                    `json:"status"`
			Data   hkodata.CurrentWeather `json:"data"`
			RSS    string                 `json:"rss"`
			// Raw    string                 `json:"raw_data"`
		}{
			Status: http.StatusOK,
			Data:   *data,
			RSS:    "http://rss.weather.gov.hk/rss/CurrentWeather.xml",
			// Raw:    data.Raw,
		})
	})

	r.HandleFunc("/api/region.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		requestID := r.Header.Get("X-Request-ID")
		source := "http://www.hko.gov.hk/wxinfo/json/region_json.xml"

		req, err := http.Get(source)
		if err != nil {
			log.Printf("request-id: %s, err: %s", requestID, err.Error())
			return
		}

		// prepare encoder for output
		enc := json.NewEncoder(w)

		// decode the RSS
		data, err := hkodata.DecodeRegionJSON(req.Body)
		if err != nil {
			log.Printf("request-id: %s, err: %s", requestID, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(struct {
				Status  int    `json:"status"`
				Message string `json:"message"`
				Source  string `json:"source"`
			}{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Source:  source,
			})
			return
		}

		w.Header().Set("Last-Modified", data.PubDate.Format(time.RFC1123))
		w.Header().Set("Expires", data.Expires().Format(time.RFC1123))
		w.WriteHeader(http.StatusOK)

		enc.Encode(struct {
			Status int             `json:"status"`
			Data   hkodata.Regions `json:"data"`
			Source string          `json:"source"`
			Notice string          `json:"notice"`
			// Raw    string                 `json:"raw_data"`
		}{
			Status: http.StatusOK,
			Data:   *data,
			Source: source,
			Notice: noticeNonpublicAPI,
			// Raw:    data.Raw,
		})

	})

	log.Printf("listen at port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), enforceHTTPS(forceHTTPS, r))
}
