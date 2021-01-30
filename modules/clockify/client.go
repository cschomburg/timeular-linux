package clockify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const BaseURL = "https://api.clockify.me/api/v1"

type Client struct {
	config Config
	client *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		config: cfg,
		client: &http.Client{},
	}
}

func (c *Client) GetTags() (TagList, error) {
	path := fmt.Sprintf(
		"/workspaces/%s/tags",
		c.config.WorkspaceId,
	)

	var tags TagList

	err := c.doRequest("GET", path, nil, &tags)

	return tags, err
}

func (c *Client) AddTimeEntry(t TimeEntry) error {
	path := fmt.Sprintf(
		"/workspaces/%s/user/%s/time-entries",
		c.config.WorkspaceId,
		c.config.UserId,
	)

	t.Start = ClockifyTime{time.Now()}

	err := c.doRequest("POST", path, t, nil)

	return err
}

func (c *Client) StopTimer() error {
	path := fmt.Sprintf(
		"/workspaces/%s/user/%s/time-entries",
		c.config.WorkspaceId,
		c.config.UserId,
	)

	t := TimeEntry{
		End: ClockifyTime{time.Now()},
	}

	err := c.doRequest("PATCH", path, t, nil)

	return err
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (c *Client) doRequest(method, path string, payload, response interface{}) error {
	b := &bytes.Buffer{}
	if payload != nil {
		err := json.NewEncoder(b).Encode(payload)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, BaseURL+path, b)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", c.config.ApiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		var msg ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
			return err
		}
		if msg.Message != "" {
			return errors.New("Clockify error response: " + msg.Message)
		}

		return errors.New("Unexpected response status:" + resp.Status)
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return err
		}
	}

	return nil
}
