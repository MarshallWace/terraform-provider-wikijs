package wikijs

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-wikijs/wikijs/schema"
	gqlc "github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Host       string
	Token      string
	HTTPClient *http.Client
	gqlClient  *gqlc.Client
}

// NewClient creates a new http client with a connection to the provided wikijs endpoint.
// The client is tested for authentication success on creation.
func NewClient(host, token string) (*Client, error) {
	ctx := context.Background()
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}))
	httpClient := http.Client{Transport: oauthClient.Transport, Timeout: 10 * time.Second}
	graphqlClient := gqlc.NewClient(host, &httpClient)
	c := Client{Token: token,
		Host:       host,
		HTTPClient: &httpClient,
		gqlClient:  graphqlClient}

	// Check connection with a simple query
	_, err := c.GetSite()
	if err != nil {
		if strings.Contains(err.Error(), "Message: Forbidden") {
			return nil, fmt.Errorf("failed to login to Wiki.js API. Check that the host and api token are correct")
		}
		return nil, err
	}

	return &c, nil
}

// query POSTS a graphql query through the hasura go-graphql-client
func query[T any](c *Client, variables map[string]interface{}) (*T, error) {
	var data T
	err := c.gqlClient.Query(context.Background(), &data, variables)
	return &data, err
}

// mutate POSTS a graphql mutation through the hasura go-graphql-client
func mutate[T any](c *Client, variables map[string]interface{}) (*T, error) {
	var data T
	err := c.gqlClient.Mutate(context.Background(), &data, variables)
	return &data, err
}

func (c *Client) GetSite() (*schema.SiteData, error) {
	return query[schema.SiteData](c, nil)
}

func (c *Client) GetGroup(id string) (*schema.QueryGroupData, error) {
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		panic(err)
	}
	variables := map[string]interface{}{
		"id": gqlc.Int(idInt),
	}

	return query[schema.QueryGroupData](c, variables)
}

func (c *Client) GetGroupList() (*schema.QueryGroupListData, error) {
	return query[schema.QueryGroupListData](c, nil)
}

func (c *Client) CreateGroup(name string) (*schema.CreateGroupData, error) {
	variables := map[string]interface{}{
		"name": gqlc.String(name),
	}
	return mutate[schema.CreateGroupData](c, variables)
}

func (c *Client) DeleteGroup(id string) (*schema.DeleteGroupData, error) {
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		panic(err)
	}
	variables := map[string]interface{}{
		"id": gqlc.Int(idInt),
	}
	return mutate[schema.DeleteGroupData](c, variables)
}

func (c *Client) UpdateGroup(id string, name string, redirectOnLogin string, permissions []string, pageRules []schema.PageRuleInput) (*schema.UpdateGroupData, error) {
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		panic(err)
	}

	variables := map[string]interface{}{
		"id":              gqlc.Int(idInt),
		"name":            gqlc.String(name),
		"redirectOnLogin": gqlc.String(redirectOnLogin),
		"permissions":     stringArrayToGqlcStringArray(permissions),
		"pageRules":       pageRules,
	}

	return mutate[schema.UpdateGroupData](c, variables)
}
