package modules

import (
	"log"

	"github.com/cschomburg/timeular-linux"
	"github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
)

func sendNotification(n notify.Notification) (uint32, error) {
	conn, err := dbus.SessionBusPrivate()
	if err != nil {
		return 0, err
	}

	if err = conn.Auth(nil); err != nil {
		return 0, err
	}

	if err = conn.Hello(); err != nil {
		return 0, err
	}

	id, err := notify.SendNotification(conn, n)

	return id, err
}

func StartNotify() chan timeular.Timeular {
	ch := make(chan timeular.Timeular)

	go func() {
		log.Println("Module notify active")
		id := uint32(0)
		var err error
		for state := range ch {
			activityName := "<none>"
			if state.Tracking.Activity != nil {
				activityName = state.Tracking.Activity.Name
			}

			n := notify.Notification{
				AppName:    "Timeular",
				ReplacesID: id,
				Summary:    "Timeular activity changed",
				Body:       "Tracking activity: " + activityName,
			}

			id, err = sendNotification(n)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	return ch
}
