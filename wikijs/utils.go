package wikijs

import gqlc "github.com/hasura/go-graphql-client"

func gqlcStringArrayToStringArray(in []gqlc.String) []string {
	o := make([]string, len(in))
	for i, r := range in {
		o[i] = string(r)
	}
	return o
}
