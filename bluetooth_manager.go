package timeular

import (
	"errors"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"log"
	"strings"
	"time"
)

const (
	orientationService        = "c7e70010-c847-11e6-8175-8c89a55d403c"
	orientationCharacteristic = "c7e70012-c847-11e6-8175-8c89a55d403c"
)

type BluetoothManager struct {
	OnOrientationChanged func(side int)
}

func getTimeularDevice(a *adapter.Adapter1) (*device.Device1, error) {
	list, err := a.GetDeviceList()
	if err != nil {
		return nil, err
	}

	for _, path := range list {
		dev, err := device.NewDevice1(path)
		if err != nil {
			return nil, err
		}
		if strings.Contains(dev.Properties.Name, "Timeular") {
			return dev, nil
		}
	}

	return nil, errors.New("Timeular not found")
}

func (bm *BluetoothManager) Run() {
	for {
		err := bm.connectAndRun()
		if err != nil {
			log.Fatalf("connection error: %s", err)
		}
	}
}

func (bm *BluetoothManager) connectAndRun() error {
	log.Println("Trying to connect to the Timeular")

	a, err := api.GetDefaultAdapter()
	if err != nil {
		return err
	}

	dev, err := getTimeularDevice(a)
	if err != nil {
		return err
	}

	if err := dev.Connect(); err != nil {
		return err
	}

	devWatch, err := dev.WatchProperties()
	if err != nil {
		return err
	}

	char, err := dev.GetCharByUUID(orientationCharacteristic)
	if err != nil {
		return err
	}

	val, err := char.ReadValue(nil)
	if err != nil {
		return err
	}

	go bm.OnOrientationChanged(int(val[0]))

	charWatch, err := char.WatchProperties()
	if err != nil {
		return err
	}

	if err := char.StartNotify(); err != nil {
		return err
	}

	log.Println("Subscribed to Timeular side changes")

	tick := time.NewTicker(1 * time.Minute)
	defer tick.Stop()

	for {
		select {
		case prop := <-charWatch:
			if prop == nil {
				return errors.New("No property received")
			}
			if prop.Name != "Value" {
				continue
			}

			val = prop.Value.([]byte)
			go bm.OnOrientationChanged(int(val[0]))

		case prop := <-devWatch:
			if prop == nil {
				return errors.New("No property received")
			}

			if prop.Name == "Connected" {
				connected := prop.Value.(bool)
				if !connected {
					log.Println("Connection lost")
					return nil
				}
			}
		case <-tick.C:
			connected, err := dev.GetConnected()
			if err != nil {
				return err
			}
			if !connected {
				log.Println("tick: Connection lost")
				return nil
			}
		}
	}
}
