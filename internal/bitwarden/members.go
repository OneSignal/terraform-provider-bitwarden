package bitwarden

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type OrganizationUserType int64

const (
	Owner   OrganizationUserType = 0
	Admin   OrganizationUserType = 1
	User    OrganizationUserType = 2
	Manager OrganizationUserType = 3
	Custom  OrganizationUserType = 4
)

type Collection struct {
	ID       string `json:"id"`
	ReadOnly bool   `json:"readOnly"`
}

type Member struct {
	Type                  OrganizationUserType `json:"type"`
	AccessAll             bool                 `json:"accessAll"`
	ExternalId            string               `json:"externalId"`
	Email                 string               `json:"email"`
	ResetPasswordEnrolled bool                 `json:"resetPasswordEnrolled"`
	Collections           []Collection         `json:"collections"`
}

type ResponseMember struct {
	Member
	// Response model
	Object string `json:"object"`
	ID     string `json:"id"`
	Name   string `json:"name"`
	// Status
	//  Invited 	= 0
	//  Accepted 	= 1
	//  Confirmed 	= 2
	//  Revoked 	= -1
	Status int64 `json:"status"`
}

func (c *client) CreateMember(ctx context.Context, member Member) (*ResponseMember, error) {
	// TODO: we don't support a collection at the moment, but an empty Collection array is needed otherwise the member
	// creation fails with the following error:
	//{
	//	"object":"error",
	//	"message":"Errors have occurred.",
	//	"errors": {
	// 		"An error has occurred.":["Value cannot be null. (Parameter 'source')"]
	//	}
	//}
	member.Collections = make([]Collection, 0)

	rb, err := json.Marshal(member)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, string(rb))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/members", c.apiURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, string(body))

	newMember := ResponseMember{}
	err = json.Unmarshal(body, &newMember)
	if err != nil {
		return nil, err
	}

	return &newMember, nil
}

func (c *client) GetMember(ctx context.Context, id string) (*ResponseMember, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/members/%s", c.apiURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	newMember := ResponseMember{}
	err = json.Unmarshal(body, &newMember)
	if err != nil {
		return nil, err
	}

	return &newMember, nil
}

func (c *client) UpdateMember(ctx context.Context, id string, member Member) (*ResponseMember, error) {
	rb, err := json.Marshal(member)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/members/%s", c.apiURL, id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	newMember := ResponseMember{}
	err = json.Unmarshal(body, &newMember)
	if err != nil {
		return nil, err
	}

	return &newMember, nil
}

func (c *client) DeleteMember(ctx context.Context, id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/members/%s", c.apiURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(ctx, req)

	return err
}
