package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/adrg/xdg"
	"github.com/cschomburg/timeular-linux"
	"github.com/cschomburg/timeular-linux/modules"
	"github.com/cschomburg/timeular-linux/modules/clockify"
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

func readConfig() (config Config, err error) {
	path := os.Getenv("TIMEULAR_CONFIG")
	if path == "" {
		path, err = xdg.ConfigFile("timeular/config.json")
		if err != nil {
			return config, err
		}
	}
	r, err := os.Open(path)
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
