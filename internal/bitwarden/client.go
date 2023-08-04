package bitwarden

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Client interface {
	CreateGroup(ctx context.Context, group Group) (*Group, error)
	GetGroup(ctx context.Context, id string) (*Group, error)
	UpdateGroup(ctx context.Context, id string, group Group) (*Group, error)
	DeleteGroup(ctx context.Context, id string) error
}
type client struct {
	apiURL      string
	oauthConfig *clientcredentials.Config
}

// NewClient creates a new BitWarden API client to interact with the BitWarden Public API
//
// See the BitWarden documentation for more information about the API
// https://bitwarden.com/help/public-api/
// https://bitwarden.com/help/api/
func NewClient(_ context.Context, clientID, clientSecret, apiUrl, authUrl string) (Client, error) {
	c := &client{
		apiURL: apiUrl,
		oauthConfig: &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     authUrl,
			Scopes:       []string{"api.organization"},
			AuthStyle:    oauth2.AuthStyleInParams,
		},
	}

	return c, nil
}

func (c *client) doRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")

	res, err := c.oauthConfig.Client(ctx).Do(req)
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
