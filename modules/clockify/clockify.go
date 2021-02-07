package clockify

import (
	"log"

	"github.com/cschomburg/timeular-linux"
)

type Config struct {
	WorkspaceId string `json:"workspace_id"`
	UserId      string `json:"user_id"`
	ApiKey      string `json:"api_key"`
}

func Start(config Config) chan timeular.Timeular {
	ch := make(chan timeular.Timeular)

	client := NewClient(config)

	go func() {
		taglist, err := client.GetTags()
		if err != nil {
			log.Println(err)
		}

		log.Println("Module clockify active")
		for state := range ch {
			act := state.Tracking.Activity
			if act != nil && act.Name != "" {
				var tags []string
				if tag := taglist.ByName(act.Tag); tag.ID != "" {
					tags = []string{tag.ID}
				}

				t := TimeEntry{
					Description: act.Name,
					TagIDs:      tags,
				}

				if err := client.AddTimeEntry(t); err != nil {
					log.Println(err)
				}
			} else {
				if err := client.StopTimer(); err != nil {
					log.Println(err)
				}
			}
		}
	}()

	return ch
}
