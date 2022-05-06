package schema

import (
	"encoding/json"
	"github.com/hasura/go-graphql-client"
)

type ResponseStatus struct {
	Succeeded graphql.Boolean
	ErrorCode graphql.Int
	Slug      graphql.String
	Message   graphql.String
}

func (res ResponseStatus) String() string {
	out, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	return string(out)
}

type DefaultResponse struct {
	ResponseResult ResponseStatus
}

func (res DefaultResponse) String() string {
	out, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	return string(out)
}
