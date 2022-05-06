package wikijs

import (
	"context"
	"github.com/hashicorp/terraform-provider-wikijs/wikijs/schema"
	gqlc "github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	Host       string
	Token      string
	HTTPClient *http.Client
	gqlClient  *gqlc.Client
}

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
	return &c, nil
}

func (c *Client) GetSite() *schema.SiteData {
	var data schema.SiteData
	err := c.gqlClient.Query(context.Background(), &data, nil)
	if err != nil {
		return nil
	}
	return &data
}

func (c *Client) GetGroup(id string) *schema.QueryGroupData {
	var data schema.QueryGroupData
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		panic(err)
	}
	variables := map[string]interface{}{
		"id": gqlc.Int(idInt),
	}

	err = c.gqlClient.Query(context.Background(), &data, variables)
	if err != nil {
		panic(err)
	}
	return &data
}

func (c *Client) GetGroupList() *schema.QueryGroupListData {
	var data schema.QueryGroupListData
	err := c.gqlClient.Query(context.Background(), &data, nil)
	if err != nil {
		panic(err)
	}
	return &data
}

func (c *Client) CreateGroup(name string) *schema.CreateGroupData {
	var data schema.CreateGroupData
	variables := map[string]interface{}{
		"name": gqlc.String(name),
	}

	err := c.gqlClient.Mutate(context.Background(), &data, variables)
	if err != nil {
		panic(err)
	}
	return &data
}

func (c *Client) DeleteGroup(id string) (*schema.DeleteGroupData, error) {
	var data schema.DeleteGroupData
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		panic(err)
	}

	variables := map[string]interface{}{
		"id": gqlc.Int(idInt),
	}

	err = c.gqlClient.Mutate(context.Background(), &data, variables)
	if err != nil {
		panic(err)
	}
	return &data, nil
}

//type PageRuleInput []map[string]interface{}

func (c *Client) UpdateGroup(id string, name string, redirectOnLogin string, permissions []string, pageRules []schema.PageRuleInput) *schema.UpdateGroupData {
	idInt, err := strconv.ParseInt(id, 10, 32)

	permissionsGqlc := make([]gqlc.String, len(permissions))
	for i, arg := range permissions {
		permissionsGqlc[i] = gqlc.String(arg)
	}

	variables := map[string]interface{}{
		"id":              gqlc.Int(idInt),
		"name":            gqlc.String(name),
		"redirectOnLogin": gqlc.String(redirectOnLogin),
		"permissions":     permissionsGqlc,
		"pageRules":       pageRules,
	}

	var data schema.UpdateGroupData

	err = c.gqlClient.Mutate(context.Background(), &data, variables)
	if err != nil {
		panic(err)
	}
	return &data
}
