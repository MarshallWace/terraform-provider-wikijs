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

func query[T any](c *Client, variables map[string]interface{}) (*T, error) {
    var data T
    err := c.gqlClient.Query(context.Background(), &data, variables)
    return &data, err
}

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

    // Convert []string to []gqlc.String
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

    return mutate[schema.UpdateGroupData](c, variables)
}
