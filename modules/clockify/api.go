package clockify

import (
	"strings"
	"time"
)

const TimeFormat = "2006-01-02T15:04:05Z"

type ClockifyTime struct {
	time.Time
}

func (t ClockifyTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}

	loc, _ := time.LoadLocation("UTC")
	stamp := time.Time(t.Time).In(loc).Format(TimeFormat)
	return []byte(`"` + stamp + `"`), nil
}

func (t *ClockifyTime) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), "`")

	if value == "null" {
		t.Time = time.Time{}
		return nil
	}

	var err error
	t.Time, err = time.Parse(TimeFormat, value)

	return err
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TagList []Tag

func (tags TagList) ByName(name string) Tag {
	for _, t := range tags {
		if t.Name == name {
			return t
		}
	}

	return Tag{}
}

type TimeEntry struct {
	Description string       `json:"description,omitempty"`
	Start       ClockifyTime `json:"start,omitempty"`
	End         ClockifyTime `json:"end,omitempty"`
	TagIDs      []string     `json:"tagIds",omitempty"`
}
