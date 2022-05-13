provider "wikijs" {
  #    host = wikjs_host_url  # Oor pass as env var WIKIJS_HOST
  #    token = wikijs_api_token # or pass as env var WIKIJS_TOKEN
}

terraform {
  required_providers {
    wikijs = {
      source  = "wikijs"
      version = "0.1"
    }
  }
}