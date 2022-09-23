// SPDX-FileCopyrightText: 2022 2022 Marshall Wace <opensource@mwam.com>
//
// SPDX-License-Identifier: GPL3

package schema

import (
	gqlc "github.com/hasura/go-graphql-client"
)

type Group struct {
	Id              gqlc.Int
	Name            gqlc.String
	IsSystem        gqlc.Boolean
	RedirectOnLogin gqlc.String
	Permissions     []gqlc.String
	PageRules       []PageRule
	// Users []User not implemented for now as no need in use case
	CreatedAt gqlc.String
	UpdatedAt gqlc.String
}

type PageRule struct {
	Id      gqlc.String   `json:"id"`
	Deny    gqlc.Boolean  `json:"deny"`
	Match   gqlc.String   `json:"match"`
	Roles   []gqlc.String `json:"roles"`
	Path    gqlc.String   `json:"path"`
	Locales []gqlc.String `json:"locales"`
}

type PageRuleInput PageRule

type QueryGroupData struct {
	Groups struct {
		Single Group `graphql:"single(id: $id)"`
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
