terraform {
    required_version = ">=0.12"

    # Configuration block to use Consul as Terraform Backend
    # backend "consul" {
    #     address = "127.0.0.1:8500"
    #     scheme  = "http"
    #     path    = "network/terraform"
    # }
}

module "local" {
  source = "./modules/local"
  services = var.services
  instance = var.instance
}

provider "local" {
  # Example of how provider instance variables would be consumed
  # if the provider supported them. (local does not)
  # address = var.provider_address
}
