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