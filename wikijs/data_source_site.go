// SPDX-FileCopyrightText: 2022 2022 Marshall Wace <opensource@mwam.com>
//
// SPDX-License-Identifier: GPL3

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
		Description: "Get Site data from the Wiki.js graphql API. This is currently incomplete and does not contain" +
			"all the fields available from the site endpoint.",

		ReadContext: dataSourceSiteRead,

		Schema: map[string]*schema.Schema{
			"host": {
				Description: "Wikijs host",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"title": {
				Description: "Wikijs title",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceSiteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*Client)

	data, err := c.GetSite()
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host", data.Site.Config.Host); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("title", data.Site.Config.Title); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
