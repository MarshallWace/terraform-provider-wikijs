package wikijs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSite(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSite,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.wikijs_site_data_source.test", "host", regexp.MustCompile("https://t-wiki")),
				),
			},
		},
	})
}

const testAccDataSourceSite = `
data "wikijs_site_data_source" "test" {
}
`
