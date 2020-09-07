package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func main() {
	var sensorID int
	flag.IntVar(&sensorID, "id", 0, "sensor ID")

	var threshold float64
	flag.Float64Var(&threshold, "t", 0, "PM2.5 threshold to alert on")

	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "display this help")

	flag.Parse()

	if sensorID == 0 || threshold == 0 {
		showHelp = true
	}

	if showHelp {
		flag.Usage()
		return
	}

	for {
		result, err := getState(sensorID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			continue
		}
		value, err := strconv.ParseFloat(result.Pm25Atm, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing value as float: %v", err)
			continue
		}
		if value <= threshold {
			fmt.Fprintf(os.Stdout, "Yeay (PM2.5: %f)\n", value)
		} else {
			fmt.Fprintf(os.Stdout, "Ouch (PM2.5: %f)\n", value)
		}
		time.Sleep(5 * time.Second)
	}
}

type purpleAirResponse struct {
	MapVersion       string            `json:"mapVersion"`
	BaseVersion      string            `json:"baseVersion"`
	MapVersionString string            `json:"mapVersionString"`
	Results          []purpleAirResult `json:"results"`
}

type purpleAirResult struct {
	ID              int     `json:"ID"`
	Label           string  `json:"Label"`
	Lat             float64 `json:"Lat"`
	Lon             float64 `json:"Lon"`
	PM25Value       string  `json:"PM2_5Value"`
	LastSeen        int     `json:"LastSeen"`
	Type            string  `json:"Type,omitempty"`
	Hidden          string  `json:"Hidden"`
	Version         string  `json:"Version,omitempty"`
	LastUpdateCheck int     `json:"LastUpdateCheck,omitempty"`
	Created         int     `json:"Created"`
	Uptime          string  `json:"Uptime,omitempty"`
	RSSI            string  `json:"RSSI,omitempty"`
	Adc             string  `json:"Adc"`
	P03Um           string  `json:"p_0_3_um"`
	P05Um           string  `json:"p_0_5_um"`
	P10Um           string  `json:"p_1_0_um"`
	P25Um           string  `json:"p_2_5_um"`
	P50Um           string  `json:"p_5_0_um"`
	P100Um          string  `json:"p_10_0_um"`
	Pm10Cf1         string  `json:"pm1_0_cf_1"`
	Pm25Cf1         string  `json:"pm2_5_cf_1"`
	Pm100Cf1        string  `json:"pm10_0_cf_1"`
	Pm10Atm         string  `json:"pm1_0_atm"`
	Pm25Atm         string  `json:"pm2_5_atm"`
	Pm100Atm        string  `json:"pm10_0_atm"`
	IsOwner         int     `json:"isOwner"`
	Humidity        string  `json:"humidity,omitempty"`
	TempF           string  `json:"temp_f,omitempty"`
	Pressure        string  `json:"pressure,omitempty"`
	AGE             int     `json:"AGE"`
	Stats           string  `json:"Stats"`
	ParentID        int     `json:"ParentID,omitempty"`
}

func getState(sensorID int) (purpleAirResult, error) {
	addr := "https://www.purpleair.com/json"
	query := url.Values{}
	query.Set("show", strconv.Itoa(sensorID))
	fullAddr := addr + "?" + query.Encode()
	resp, err := http.Get(fullAddr)
	if err != nil {
		return purpleAirResult{}, fmt.Errorf("getting state from purple air: %w", err)
	}

	var res purpleAirResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return purpleAirResult{}, fmt.Errorf("decoding state from purple air: %w", err)
	}

	return res.Results[0], nil
}
