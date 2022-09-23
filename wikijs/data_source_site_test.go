// SPDX-FileCopyrightText: 2022 2022 Marshall Wace <opensource@mwam.com>
//
// SPDX-License-Identifier: GPL3

package wikijs

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSite(t *testing.T) {
	host := os.Getenv("WIKIJS_HOST")
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSite,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.wikijs_site_data_source.test", "host", regexp.MustCompile(host)),
				),
			},
		},
	})
}

const testAccDataSourceSite = `
data "wikijs_site_data_source" "test" {
}
`
