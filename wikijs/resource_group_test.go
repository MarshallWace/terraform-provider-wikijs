package wikijs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceGroup(t *testing.T) {
	//t.Skip("resource not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"wikijs_group_resource.foo", "name", regexp.MustCompile("test-group")),
				),
			},
			{
				Config: testAccResourceGroupUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"wikijs_group_resource.foo", "name", regexp.MustCompile("test-group-updated")),
				),
			},
		},
	})
}

const testAccResourceGroup = `
resource "wikijs_group_resource" "foo" {
    name = "test-group"
    permissions = ["read:pages", "write:pages"]
    redirect_on_login = ""
    page_rules {
        id = "page_rules_dummy_id"
        deny = false
        match = "START"
        roles = ["read:pages","write:pages"]
        path = "/test"
        locales = []
    }
}
`

const testAccResourceGroupUpdated = `
resource "wikijs_group_resource" "foo" {
    name = "test-group-updated"
    permissions = ["read:pages", "write:pages"]
    redirect_on_login = ""
    page_rules {
        id = "page_rules_dummy_id"
        deny = false
        match = "START"
        roles = ["read:pages","write:pages"]
        path = "/test"
        locales = []
    }
}
`
