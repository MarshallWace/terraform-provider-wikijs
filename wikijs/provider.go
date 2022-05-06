package wikijs

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(_ string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"wikijs_site_data_source": dataSourceSite(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"wikijs_group_resource": resourceGroup(),
			},
		}

		p.ConfigureContextFunc = configure()

		return p
	}
}

func configure() func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
		host := os.Getenv("WIKIJS_HOST")
		token := os.Getenv("WIKIJS_TOKEN")

		client, err := NewClient(host, token)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, nil
	}
}
