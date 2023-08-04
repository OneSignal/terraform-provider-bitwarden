package bitwarden

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Group struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	Name       string `json:"name"`
	ExternalId string `json:"externalId"`
	AccessAll  bool   `json:"accessAll"`
}

func (c *client) CreateGroup(ctx context.Context, group Group) (*Group, error) {
	rb, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/groups", c.apiURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	newGroup := Group{}
	err = json.Unmarshal(body, &newGroup)
	if err != nil {
		return nil, err
	}

	return &newGroup, nil
}

func (c *client) GetGroup(ctx context.Context, id string) (*Group, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/groups/%s", c.apiURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	newGroup := Group{}
	err = json.Unmarshal(body, &newGroup)
	if err != nil {
		return nil, err
	}

	return &newGroup, nil
}

func (c *client) UpdateGroup(ctx context.Context, id string, group Group) (*Group, error) {
	rb, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/groups/%s", c.apiURL, id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	newGroup := Group{}
	err = json.Unmarshal(body, &newGroup)
	if err != nil {
		return nil, err
	}

	return &newGroup, nil
}

func (c *client) DeleteGroup(ctx context.Context, id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/groups/%s", c.apiURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(ctx, req)

	return err
}
