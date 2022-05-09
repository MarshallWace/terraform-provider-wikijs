package schema

import (
	gqlc "github.com/hasura/go-graphql-client"
)

type Strategies []AuthenticationStrategy

type AuthenticationStrategy struct {
	Key          gqlc.String    `json:"key"`
	Props        []KeyValuePair `json:"props"`
	Title        gqlc.String    `json:"title"`
	Description  gqlc.String    `json:"description"`
	IsAvailable  gqlc.String    `json:"isAvailable"`
	UseForm      gqlc.Boolean   `json:"useForm"`
	UsernameType gqlc.String    `json:"UsernameType"`
	Logo         gqlc.String    `json:"logo"`
	Color        gqlc.String    `json:"color"`
	Website      gqlc.String    `json:"website"`
	Icon         gqlc.String    `json:"icon"`
}

type AuthenticationActiveStrategy struct {
	Key              gqlc.String            `json:"key"`
	Strategy         AuthenticationStrategy `json:"strategy"`
	DisplayName      gqlc.String            `json:"displayName"`
	Order            gqlc.Int               `json:"order"`
	IsEnabled        gqlc.Boolean           `json:"isEnabled"`
	Config           []KeyValuePair         `json:"config"`
	SelfRegistration gqlc.Boolean           `json:"selfRegistration"`
	DomainWhitelist  gqlc.String            `json:"domainWhitelist"`
	AutoEnrollGroups gqlc.Int               `json:"autoEnrollGroups"`
}

type AuthenticationStrategyInput AuthenticationActiveStrategy

type QueryStrategiesData struct {
	Authentication struct {
		Strategies []AuthenticationStrategy
	}
}

type QueryActiveStrategiesData struct {
	Authentication struct {
		ActiveStrategies []AuthenticationActiveStrategy `graphql:"activeStrategies(enabledOnly: $enabledOnly)"` // enabledOnly: boolean
	}
}

type UpdateActiveStrategiesData struct {
	Authentication struct {
		UpdateStrategies DefaultResponse `graphql:"updateStrategies(strategies: $strategies)"` // strategies: []AuthenticationActiveStrategy
	}
}

type StrategyPropData struct {
	Default   string `json:"default"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Hint      string `json:"hint"`
	Enum      bool   `json:"enum"`
	Multiline bool   `json:"multiline"`
	Sensitive bool   `json:"sensitive"`
	MaxWidth  int    `json:"maxWidth"`
	Order     int    `json:"order"`
	Value     string `json:"value"`
}
