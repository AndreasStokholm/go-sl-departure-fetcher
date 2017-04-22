package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-restit/lzjson"

	"github.com/spf13/viper"

	"github.com/yosssi/gmq/mqtt"
	mqttClient "github.com/yosssi/gmq/mqtt/client"
)

func main() {
	log.Println("Starting...")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	configErr := viper.ReadInConfig()
	if configErr != nil {
		panic(fmt.Errorf("fatal error config file: %s", configErr))
	}

	mqttHost := viper.GetString("mqtt.host")
	mqttClientID := []byte(viper.GetString("mqtt.client_id"))
	mqttUser := []byte(viper.GetString("mqtt.user"))
	mqttPassword := []byte(viper.GetString("mqtt.password"))
	mqttTopic := []byte(viper.GetString("mqtt.topic"))
	mqttStatus := []byte(viper.GetString("mqtt.status"))

	mqttCli := mqttClient.New(&mqttClient.Options{
		ErrorHandler: func(err error) {
			log.Fatal(err)
		},
	})

	defer mqttCli.Terminate()
	mqttErr := mqttCli.Connect(&mqttClient.ConnectOptions{
		Network:  "tcp",
		Address:  mqttHost,
		ClientID: mqttClientID,
		UserName: mqttUser,
		Password: mqttPassword,
	})

	if mqttErr != nil {
		log.Fatal(mqttErr)
	}

	mqttErr = mqttCli.Publish(&mqttClient.PublishOptions{
		QoS:       mqtt.QoS0,
		Retain:    true,
		TopicName: []byte(mqttStatus),
		Message:   []byte("Ready!"),
	})

	if mqttErr != nil {
		log.Fatal(mqttErr)
	}

	for {
		departures := getDepartures()
		if len(departures) > 0 {
			departureString, error := json.Marshal(departures)
			if error != nil {
				fmt.Println(error)
				return
			}

			mqttCli.Publish(&mqttClient.PublishOptions{
				QoS:       mqtt.QoS0,
				Retain:    true,
				TopicName: []byte(mqttTopic),
				Message:   []byte(departureString),
			})
		}

		log.Println("Sleeping for 30 seconds.")
		time.Sleep(30 * time.Second)
	}
}

func getDepartures() []MQTTDeparture {

	var departureSites []DepartureSite
	viper.UnmarshalKey("departures.sites", &departureSites)

	var departureMessage []MQTTDeparture

	for _, site := range departureSites {
		var siteID = strconv.Itoa(site.SiteID)
		var transportMode = site.TransportMode

		response, err := http.Get("http://api.sl.se/api2/realtimedeparturesV4.json?key=" + viper.GetString("departures.api.key") + "&siteid=" + siteID + "&timewindow=" + viper.GetString("departures.timewindow"))
		if err != nil {
			log.Fatal(err)
			return departureMessage
		}

		var departures []Departure
		json := lzjson.Decode(response.Body).Get("ResponseData").Get(transportMode)
		json.Unmarshal(&departures)

		for _, departure := range departures {

			sort.Ints(site.Directions)
			directionCount := sort.Search(len(site.Directions), func(directionCount int) bool {
				return site.Directions[directionCount] >= departure.JourneyDirection
			})

			if directionCount < len(site.Directions) && site.Directions[directionCount] == departure.JourneyDirection {
				formattedDeparture := MQTTDeparture{
					TransportMode:    departure.TransportMode,
					TimeOfDeparture:  departure.DisplayTime,
					JourneyDirection: departure.JourneyDirection,
				}
				departureMessage = append(departureMessage, formattedDeparture)
			}
		}
	}

	fmt.Println(departureMessage)
	return departureMessage
}
