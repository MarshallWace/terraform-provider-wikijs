provider "wikijs" {
  # example configuration here
}

terraform {
  required_providers {
    wikijs = {
      source  = "local/mwam/wikijs"
      version = "0.1"
    }
  }
}

#data "wikijs_site_data_source" "all" {}
#
#output "site_info" {
#    value = data.wikijs_site_data_source.all
#}

locals {
  group_page_rules_yaml = yamldecode(file("${path.module}/group_page_rules.yaml"))
  group_page_rules = { for idx, val in local.group_page_rules_yaml["groups"] :
    idx => val
  }
}

resource "wikijs_group_resource" "technology" {
  name              = "test_technology"
  permissions       = ["read:pages", "write:pages"]
  redirect_on_login = ""
  dynamic "page_rules" {
    for_each = local.group_page_rules["technology"]
    content {
      id      = page_rules.value.id
      deny    = page_rules.value.deny
      match   = page_rules.value.match
      roles   = page_rules.value.roles
      path    = page_rules.value.path
      locales = []
    }
  }
}

resource "wikijs_group_resource" "infrastructure" {
  name              = "test_infrastructure"
  permissions       = ["read:pages", "write:pages"]
  redirect_on_login = ""
  dynamic "page_rules" {
    for_each = local.group_page_rules["infrastructure"]
    content {
      id      = page_rules.value.id
      deny    = page_rules.value.deny
      match   = page_rules.value.match
      roles   = page_rules.value.roles
      path    = page_rules.value.path
      locales = []
    }
  }
}

#output "group_info" {
#    value =  wikijs_group_resource.foo
#}