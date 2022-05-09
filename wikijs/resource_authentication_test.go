package wikijs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAuthentication(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAuthentication,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"wikijs_authentication_resource.foo.strategies.{0}.{key}", "name", regexp.MustCompile("randomkey1")),
				),
			},
		},
	})
}

const testAccResourceAuthentication = `
resource "wikijs_authentication_resource" "foo" {
    strategies {
        key = "randomkey1"
        strategyKey = "ldap"
        config {
        }
        displayName = "ldap"
        order = 0
        isEnabled = true
        selfRegistration = false
        autoEnrollGroups = [2]
        domainWhitelist = []
    }
    strategies {
        key = "randomkey2"
        strategyKey = "local"
        config {
        }
        displayName = "local"
        order = 1
        isEnabled = true
        selfRegistration = true
        autoEnrollGroups = []
        domainWhitelist = []
    }
}
`
