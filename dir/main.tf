terraform {
    required_version = ">=0.12"
}

module "local" {
    source = "./modules/local"
    services = { for name in var.service_mapping[terraform.workspace] : name => var.services[name]}
}
