terraform {
    required_version = ">=0.12"
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
