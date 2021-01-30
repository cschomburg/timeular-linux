package main

import (
	"encoding/json"
	"github.com/cschomburg/timeular-linux"
	"github.com/cschomburg/timeular-linux/modules"
	"github.com/cschomburg/timeular-linux/modules/clockify"
	"log"
	"os"
	"time"
)

const (
	orientationServiceNew        = "c7e70010-c847-11e6-8175-8c89a55d403c"
	orientationCharacteristicNew = "c7e70012-c847-11e6-8175-8c89a55d403c"
)

type Config struct {
	Notify     bool                `json:"notify"`
	LogPath    string              `json:"log_path"`
	Clockify   clockify.Config     `json:"clockify"`
	Activities []timeular.Activity `json:"activities"`
}

func readConfig() (Config, error) {
	var config Config
	r, err := os.Open("./config.json")
	if err != nil {
		return config, err
	}

	err = json.NewDecoder(r).Decode(&config)
	return config, err
}

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatalln("Could not read config:", err)
	}

	hub := timeular.NewHub()
	go hub.Run()

	if config.Notify {
		hub.Register(modules.StartNotify())
	}
	if config.LogPath != "" {
		hub.Register(modules.StartLogger(config.LogPath))
	}
	if config.Clockify.ApiKey != "" {
		hub.Register(clockify.Start(config.Clockify))
	}

	state := &timeular.Timeular{
		Activities: config.Activities,
	}

	manager := timeular.BluetoothManager{
		OnOrientationChanged: func(sideID int) {
			activity := state.GetActivity(sideID)
			activityName := "no activity"
			if activity != nil {
				activityName = activity.Name
			}

			log.Printf("Device side: %d - %s", sideID, activityName)

			state.CurrentSide = sideID
			state.Tracking = &timeular.CurrentTracking{
				Activity:  activity,
				StartedAt: time.Now(),
			}

			go hub.Broadcast(state)
		},
	}

	manager.Run()

	<-make(chan struct{})
}
