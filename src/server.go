package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/maciej/bme280"
	"golang.org/x/exp/io/i2c"
)

func getBME280() *bme280.Driver {
	device, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x76)
	if err != nil {
		log.Fatal(err)
	}

	driver := bme280.New(device)
	err = driver.InitWith(bme280.ModeForced, bme280.Settings{
		Filter:                  bme280.FilterOff,
		Standby:                 bme280.StandByTime1000ms,
		PressureOversampling:    bme280.Oversampling16x,
		TemperatureOversampling: bme280.Oversampling16x,
		HumidityOversampling:    bme280.Oversampling16x,
	})

	if err != nil {
		log.Fatal(err)
	}

	return driver
}

func handleJSON(w http.ResponseWriter, r *http.Request) {

	bme280 := getBME280()
	defer bme280.Close()

	response, err := bme280.Read()
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{}{
		"p": response.Pressure,
		"t": response.Temperature,
		"h": response.Humidity,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/api", handleJSON)
	http.ListenAndServe(":80", nil)
}
