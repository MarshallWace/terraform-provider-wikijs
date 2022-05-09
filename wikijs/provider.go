package wikijs

import (
    "context"
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
			Schema: map[string]*schema.Schema{
				"host": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("WIKIJS_HOST", nil),
				},
				"token": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("WIKIJS_TOKEN", nil),
				},
			},
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

func configure() func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		host := d.Get("host").(string)
		token := d.Get("token").(string)
		if host == "" {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error,
				Summary: "Wikijs HOST not declared.",
				Detail:  "Set the value as an env var WIKIJS_HOST or as `host` in the provider block."})
			return nil, diags
		}

		if token == "" {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error,
				Summary: "Wikijs API TOKEN not declared.",
				Detail:  "Set the value as an env var WIKIJS_TOKEN or as `token` in the provider block."})
			return nil, diags
		}

		client, err := NewClient(host, token)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, nil
	}
}
