terraform {
  cloud {
    hostname     = "app.terraform.io"
    organization = "amiasea"

    workspaces {
      name = "amiasea-tfcred-e2e"
    }
  }
}