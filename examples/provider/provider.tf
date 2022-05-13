provider "wikijs" {
  host = "https://your-wiki-url.com" # Or pass as env var WIKIJS_HOST
  #    token = wikijs_api_token # or pass as env var WIKIJS_TOKEN
}

terraform {
  required_providers {
    wikijs = {
      source  = "MarshallWace/wikijs"
      version = "0.0.1"
    }
  }
}