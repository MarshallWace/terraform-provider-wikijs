resource "wikijs_group_resource" "my_group" {
  name              = "my_group_name"
  permissions       = ["read:pages", "write:pages"]
  redirect_on_login = ""
  page_rules {
    id      = "page_rule_1"
    deny    = false
    match   = "START"
    path    = "my_path"
    roles   = ["read:pages"]
    locales = []
  }
}