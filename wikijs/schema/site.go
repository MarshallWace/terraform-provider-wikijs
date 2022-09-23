// SPDX-FileCopyrightText: 2022 2022 Marshall Wace <opensource@mwam.com>
//
// SPDX-License-Identifier: GPL3

package schema

type SiteGraphqlResponse struct {
	Data SiteData `json:"data"`
}
type SiteData struct {
	Site Site `json:"site"`
}

type Site struct {
	Config SiteConfig `json:"config"`
}
type SiteConfig struct {
	Host        string `json:"host"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
