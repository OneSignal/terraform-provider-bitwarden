package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Client interface {

	//def get_members(self):
	//
	//def get_groups(self):
	//
	//def get_member_groups(self):
	//
	//def create_member(self, member):
	//
	//def delete_member(self, member):
	//
	//def _get_member_ids(self, members):
	//
	//def create_group(self, group):
	//
	//def delete_group(self, group):
}

type client struct {
	apiURL      string
	oauthConfig *clientcredentials.Config
	ctx         context.Context
}

// NewClient creates a new BitWarden API client to interact with the BitWarden Public API
//
// https://bitwarden.com/help/public-api/
func NewClient(ctx context.Context, clientID, clientSecret, apiUrl string) (Client, error) {
	// Additional Endpoint Params needed to authenticate with BitWarden API
	values := url.Values{}
	values.Set("grant_type", "client_credentials")

	c := &client{
		apiURL: apiUrl,
		oauthConfig: &clientcredentials.Config{
			ClientID:       clientID,
			ClientSecret:   clientSecret,
			TokenURL:       "https://identity.bitwarden.com/connect/token",
			Scopes:         []string{"api.organization"},
			EndpointParams: values,
			AuthStyle:      oauth2.AuthStyleInParams,
		},
		ctx: ctx,
	}

	return c, nil
}

func (c *client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")

	res, err := c.oauthConfig.Client(c.ctx).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
