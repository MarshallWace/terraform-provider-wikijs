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
    "log"
    "strconv"
    "time"
)

func resourceGroup() *schema.Resource {
    return &schema.Resource{
        // This description is used by the documentation generator and the language server.
        Description: "Sample resource in the Terraform provider group.",

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

func resourceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    c := meta.(*Client)
    id := d.Id()
    data, err := c.GetGroup(id)
    if err != nil {
        return diag.FromErr(err)
    }
    if err := d.Set("name", data.Groups.Single.Name); err != nil {
        return diag.FromErr(err)
    }
    if err := d.Set("is_system", data.Groups.Single.IsSystem); err != nil {
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
    log.Println("resourceGroupCreate result: ", data.Groups.Create.ResponseResult.String())

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
    var diags diag.Diagnostics
    c := meta.(*Client)
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

    data, err := c.UpdateGroup(id, name, redirectOnLogin, permissions, pageRules)
    if err != nil {
        return diag.FromErr(err)
    }
    if err := d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
        return diag.FromErr(err)
    }

    log.Println("resourceGroupUpdate result: ", data.Groups.Update.String())
    tflog.Trace(ctx, fmt.Sprintf("Updated resource name %s with its values", name))

    return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    c := meta.(*Client)
    id := d.Id()
    name := d.Get("name").(string)
    data, err := c.DeleteGroup(id)
    if err != nil {
        return diag.FromErr(err)
    }
    d.SetId("")
    log.Println("resourceGroupDelete result: ", data.Groups.Delete.String())
    tflog.Trace(ctx, fmt.Sprintf("Deleted resource name %s", name))

    return diags
}
