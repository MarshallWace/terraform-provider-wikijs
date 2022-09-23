// SPDX-FileCopyrightText: 2022 2022 Marshall Wace <opensource@mwam.com>
//
// SPDX-License-Identifier: GPL3

package wikijs

import gqlc "github.com/hasura/go-graphql-client"

func gqlcStringArrayToStringArray(in []gqlc.String) []string {
	o := make([]string, len(in))
	for i, r := range in {
		o[i] = string(r)
	}
	return o
}

func stringArrayToGqlcStringArray(in []string) []gqlc.String {
	o := make([]gqlc.String, len(in))
	for i, r := range in {
		o[i] = gqlc.String(r)
	}
	return o
}
