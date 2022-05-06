package schema

import (
	gqlc "github.com/hasura/go-graphql-client"
)

type Group struct {
	Id       gqlc.Int
	Name     gqlc.String
	IsSystem gqlc.Boolean
}

type PageRuleInput struct {
	Id      gqlc.String   `json:"id"`
	Deny    gqlc.Boolean  `json:"deny"`
	Match   gqlc.String   `json:"match"` // this is actually an enum on wikijs graphql
	Roles   []gqlc.String `json:"roles"`
	Path    gqlc.String   `json:"path"`
	Locales []gqlc.String `json:"locales"`
}

type QueryGroupData struct {
	Groups struct {
		Single struct {
			Id       gqlc.Int
			Name     gqlc.String
			IsSystem gqlc.Boolean
		} `graphql:"single(id: $id)"`
	}
}

type CreateGroupData struct {
	Groups struct {
		Create struct {
			Group          Group
			ResponseResult ResponseStatus
		} `graphql:"create(name: $name)"`
	}
}

type UpdateGroupData struct {
	Groups struct {
		Update DefaultResponse `graphql:"update(id: $id, name: $name, redirectOnLogin: $redirectOnLogin, permissions: $permissions, pageRules: $pageRules)"`
	}
}

type DeleteGroupData struct {
	Groups struct {
		Delete DefaultResponse `graphql:"delete(id: $id)"`
	}
}

type QueryGroupListData struct {
	Groups struct {
		List []struct {
			Id       gqlc.String
			Name     gqlc.String
			IsSystem gqlc.Boolean
		} `graphql:"single(id: $id)"`
	}
}
