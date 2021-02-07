package modules

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cschomburg/timeular-linux"
)

func StartLogger(path string) chan timeular.Timeular {
	ch := make(chan timeular.Timeular)

	go func() {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()

		log.Println("Module logger active")
		for state := range ch {
			activityName := "<none>"
			if state.Tracking.Activity != nil {
				activityName = state.Tracking.Activity.Name
			}

			_, err := fmt.Fprintf(
				f,
				"%s %d %s\n",
				state.Tracking.StartedAt.Format(time.RFC3339),
				state.CurrentSide,
				activityName,
			)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	return ch
}
