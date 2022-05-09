package wikijs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wjSchema "github.com/hashicorp/terraform-provider-wikijs/wikijs/schema"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slices"
	"log"
	"strconv"
	"time"
)

func resourceAuthentication() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Authentication",

		CreateContext: resourceAuthenticationCreate,
		ReadContext:   resourceAuthenticationRead,
		UpdateContext: resourceAuthenticationUpdate,
		DeleteContext: resourceAuthenticationDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "id",
				Computed:    true,
			},
			"strategies": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true, // computed from key
						},
						"key": { // key of the active strategy (user generated)
							Type:     schema.TypeString,
							Required: true,
						},
						"strategyKey": { // permanent key of the strategy (ie. LDAP)
							Type:     schema.TypeString,
							Required: true,
						},
						"config": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"displayName": {
							Type:     schema.TypeString,
							Required: true,
						},
						"order": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"isEnabled": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"selfRegistration": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"domainWhitelist": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"autoEnrollGroups": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
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

func resourceAuthenticationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*Client)

	existingStrategies, err := c.GetActiveStrategies(false)
	if err != nil {
		panic(err)
	}

	if err = d.Set("strategies", existingStrategies.Authentication.ActiveStrategies); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAuthenticationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)
	var diags diag.Diagnostics

	// todo
	// get strategies for the props
	// get desired strategies from tf
	// get existing strategies
	// merge above two, since there's only an updateStategies endpoint which takes a list
	// create the json needed based on the props
	// validate the strategies that can exist
	// post

	allStrategies, err := c.GetStrategies()
	if err != nil {
		panic(err)
	}

	allStrategiesPropsMap := make(map[string]interface{}, len(allStrategies.Authentication.Strategies))
	for _, arg := range allStrategies.Authentication.Strategies {
		strategyKey := string(arg.Key)
		pMap := make(map[string]interface{}, len(arg.Props))
		for _, p := range arg.Props {
			var pd wjSchema.StrategyPropData
			//s, _ := strconv.Unquote(string(val))
			err := json.Unmarshal([]byte(p.Value), &pd)
			if err != nil {
				return nil
			}
			pMap[string(p.Key)] = pd
		}
		allStrategiesPropsMap[strategyKey] = pMap
	}

	newStrategies := d.Get("strategies")
	name := d.Get("name").(string)

	data := c.CreateAuthentication(name)

	log.Println("resourceAuthenticationCreate result: ", data.Authentications.Create.ResponseResult.String())

	d.SetId(strconv.Itoa(int(data.Authentications.Create.Authentication.Id)))

	tflog.Trace(ctx, fmt.Sprintf("created a resource with name %s", name))

	// Perform an update because GraphQL Create only creates a group with 'name' input, other values are the default
	updateDiags := resourceAuthenticationUpdate(ctx, d, meta)
	if updateDiags.HasError() {
		resourceAuthenticationDelete(ctx, d, meta)
		return updateDiags
	}

	tflog.Trace(ctx, fmt.Sprintf("Updated resource name %s with its values", name))

	return diags
}

func resourceAuthenticationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*Client)

	// same as create

	id := d.Id()
	name := d.Get("name").(string)
	redirectOnLogin := d.Get("redirect_on_login").(string)

	_permissions := d.Get("permissions").(*schema.Set).List()
	permissions := make([]string, len(_permissions))
	for i, arg := range _permissions {
		permissions[i] = arg.(string)
	}

	_pageRules := d.Get("page_rules").([]interface{})
	pageRules := make([]wjSchema.PageRuleInput, len(_pageRules))
	for i, arg := range _pageRules {
		var temp wjSchema.PageRuleInput
		err := mapstructure.Decode(arg, &temp)
		if err != nil {
			panic(err)
		}
		pageRules[i] = temp
	}

	for _, rule := range pageRules {
		roles := rule.Roles
		for _, role := range roles {
			if !slices.Contains(permissions, string(role)) {
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Tried to set PageRule role for unallowed global permission '%s' in the pagerule block of id: %s", role, rule.Id),
					Detail: fmt.Sprintf("In order to set a role for a page rule, that role must first be enabled under global permissions." +
						" Add the permission to wikijs_group_resource.permissions."),
				})
			}
		}
	}

	res := c.UpdateAuthentication(id, name, redirectOnLogin, permissions, pageRules)
	if err := d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
		return diag.FromErr(err)
	}

	log.Println("resourceAuthenticationUpdate result: ", res.Authentications.Update.String())
	tflog.Trace(ctx, fmt.Sprintf("Updated resource name %s with its values", name))

	return resourceAuthenticationRead(ctx, d, meta)
}

func resourceAuthenticationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*Client)
	id := d.Id()
	name := d.Get("name").(string)
	res, err := c.DeleteAuthentication(id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Println("resourceAuthenticationDelete result: ", res.Authentications.Delete.String())
	tflog.Trace(ctx, fmt.Sprintf("Deleted resource name %s", name))

	return diags
}
