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

data "wikijs_site_data_source" "all" {}

output "site_info" {
    value = data.wikijs_site_data_source.all
}
