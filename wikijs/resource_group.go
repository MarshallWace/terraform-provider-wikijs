package wikijs

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wjSchema "github.com/hashicorp/terraform-provider-wikijs/wikijs/schema"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"time"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Updates Wiki.js Groups via it's graphql API. ",

		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "id",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name",
			},
			"is_system": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "isSystem",
			},
			"redirect_on_login": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "redirect on login path",
			},
			"permissions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "permissions",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "createdAt",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "updatedAt",
			},
			"page_rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"deny": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"match": {
							Type:     schema.TypeString,
							Required: true,
						},
						"roles": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"locales": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*Client)
	id := d.Id()
	name := d.Get("name")
	data, err := c.GetGroup(id)
	if data.Groups.Single.Name == "" && data.Groups.Single.Id == 0 {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{Severity: diag.Warning, Summary: fmt.Sprintf("group with id %s "+
			"and name %s no longer exists due to a change outside of terraform. it has been deleted from the state", id, name)})
		return diags
	}

	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", data.Groups.Single.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_system", data.Groups.Single.IsSystem); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("redirect_on_login", data.Groups.Single.RedirectOnLogin); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("permissions", data.Groups.Single.Permissions); err != nil {
		return diag.FromErr(err)
	}
	flattenedPageRules := make([]interface{}, len(data.Groups.Single.PageRules))
	for i, pr := range data.Groups.Single.PageRules {
		oi := make(map[string]interface{})
		oi["id"] = string(pr.Id)
		oi["deny"] = bool(pr.Deny)
		oi["match"] = string(pr.Match)
		oi["roles"] = gqlcStringArrayToStringArray(pr.Roles)
		oi["path"] = string(pr.Path)
		oi["locales"] = gqlcStringArrayToStringArray(pr.Locales)
		flattenedPageRules[i] = oi
	}
	if err := d.Set("page_rules", flattenedPageRules); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_at", data.Groups.Single.CreatedAt); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("updated_at", data.Groups.Single.UpdatedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	data, err := c.CreateGroup(name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(data.Groups.Create.Group.Id)))

	tflog.Trace(ctx, fmt.Sprintf("created a resource with name %s", name))

	// Perform an update because GraphQL Create only creates a group with 'name' input, other values are the default
	updateDiags := resourceGroupUpdate(ctx, d, meta)
	if updateDiags.HasError() {
		resourceGroupDelete(ctx, d, meta)
		return updateDiags
	}

	tflog.Trace(ctx, fmt.Sprintf("Updated resource name %s with its values", name))

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)
	id := d.Id()
	name := d.Get("name").(string)
	redirectOnLogin := d.Get("redirect_on_login").(string)
	if !strings.HasPrefix(redirectOnLogin, "/") {
		return diag.FromErr(fmt.Errorf("redirectOnLogin must start with /"))
	}

	globalPermissions := getGlobalPermissions(d)

	pageRules, diags := getPageRules(d)
	if diags != nil {
		return diags
	}

	validationError := validatePageRules(pageRules, globalPermissions)
	if validationError != nil {
		return validationError
	}

	_, err := c.UpdateGroup(id, name, redirectOnLogin, globalPermissions, pageRules)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
		return diag.FromErr(err)
	}

	tflog.Trace(ctx, fmt.Sprintf("Updated resource name %s with its values", name))

	return resourceGroupRead(ctx, d, meta)
}

func getGlobalPermissions(d *schema.ResourceData) []string {
	_globalPermissions := d.Get("permissions").(*schema.Set).List()
	globalPermissions := make([]string, len(_globalPermissions))
	for i, arg := range _globalPermissions {
		globalPermissions[i] = arg.(string)
	}
	return globalPermissions
}

func getPageRules(d *schema.ResourceData) ([]wjSchema.PageRuleInput, diag.Diagnostics) {
	var diags diag.Diagnostics
	_pageRules := d.Get("page_rules").([]interface{})
	pageRules := make([]wjSchema.PageRuleInput, len(_pageRules))
	for i, arg := range _pageRules {
		var p wjSchema.PageRuleInput
		err := mapstructure.Decode(arg, &p)
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(string(p.Path), "/") {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("page_rules.path \"%s\" must not start with /", p.Path),
				Detail:   fmt.Sprintf("Remove the / from the start of the page_rule path. This is added automatically by wikijs."),
			})
		}
		pageRules[i] = p
	}
	if len(diags) > 0 {
		return nil, diags
	}
	return pageRules, nil
}

func validatePageRules(pageRules []wjSchema.PageRuleInput, globalPermissions []string) diag.Diagnostics {
	var diags diag.Diagnostics
	for _, rule := range pageRules {
		roles := rule.Roles
		for _, role := range roles {
			if !slices.Contains(globalPermissions, string(role)) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Tried to set PageRule role for unallowed global permission '%s' in the pagerule block of id: %s", role, rule.Id),
					Detail: fmt.Sprintf("In order to set a role for a page rule, that role must first be enabled under global permissions." +
						" Add the permission to wikijs_group_resource.permissions."),
				})
			}
		}
	}
	if len(diags) > 0 {
		return diags
	} else {
		return nil
	}
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*Client)
	id := d.Id()
	name := d.Get("name").(string)
	_, err := c.DeleteGroup(id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	tflog.Trace(ctx, fmt.Sprintf("Deleted resource name %s", name))

	return diags
}
