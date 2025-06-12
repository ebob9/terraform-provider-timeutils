terraform {
  required_providers {
    timeutils = {
      source = "ebob9/timeutils"
      version = "~> 0.1"
    }
  }
}

provider "timeutils" {}
