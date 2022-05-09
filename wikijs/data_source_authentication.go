package wikijs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceAuthentication() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource for Authentication Strategies",

		ReadContext: dataSourceAuthenticationRead,

		Schema: map[string]*schema.Schema{
			"strategies": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": { // permanent key of the strategy (ie. LDAP)
							Type:     schema.TypeString,
							Computed: true,
						},
						"props": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isAvailable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"useForm": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"usernameType": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"logo": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"color": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"website": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"icon": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAuthenticationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	var diags diag.Diagnostics
	c := meta.(*Client)

	data, err := c.GetStrategies()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("strategies", data.Authentication.Strategies); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
