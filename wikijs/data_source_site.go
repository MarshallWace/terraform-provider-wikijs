package wikijs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceSite() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource for SiteQuery",

		ReadContext: dataSourceSiteRead,

		Schema: map[string]*schema.Schema{
			"host": {
				// This description is used by the documentation generator and the language server.
				Description: "Wikijs host",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"title": {
				// This description is used by the documentation generator and the language server.
				Description: "Wikijs title",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceSiteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	var diags diag.Diagnostics
	c := meta.(*Client)

	data := c.GetSite()

	if err := d.Set("host", data.Site.Config.Host); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("title", data.Site.Config.Title); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
