package castai

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/castai/terraform-provider-castai/castai/sdk"
)

type Config struct {
	ApiUrl   string
	ApiToken string
}

type ProviderConfig struct {
	api *sdk.ClientWithResponses
}

func (c *Config) configureProvider() (interface{}, error) {
	httpClientOption := func(client *sdk.Client) error {
		client.Client = &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   1 * time.Minute,
		}
		return nil
	}

	apiTokenOption := sdk.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-API-Key", c.ApiToken)
		return nil
	})

	client, err := sdk.NewClientWithResponses(c.ApiUrl, httpClientOption, apiTokenOption)
	if err != nil {
		return nil, err
	}

	if checkErr := sdk.CheckGetResponse(client.ListAuthTokensWithResponse(context.Background())); checkErr != nil {
		return nil, fmt.Errorf("validating api token (by listing auth tokens): %v", checkErr)
	}

	return &ProviderConfig{
		api: client,
	}, nil
}
