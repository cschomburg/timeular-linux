package timeular

import (
	"time"
)

type Activity struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Integration string `json:"integration"`
	DeviceSide  int    `json:"deviceSide"`
	Tag         string `json:"tag"`
}

type CurrentTracking struct {
	Activity  *Activity `json:"activity"`
	StartedAt time.Time `json:"startedAt"`
}

type TimeEntry struct {
	ID       string   `json:"id"`
	Activity Activity `json:"activity"`
	Note     string   `json:"note"`
}

type Timeular struct {
	CurrentSide int
	Tracking    CurrentTracking
	Activities  []Activity
}

func (timeular Timeular) GetActivity(deviceSide int) *Activity {
	for _, a := range timeular.Activities {
		if a.DeviceSide == deviceSide {
			return &a
		}
	}

	return nil
}
