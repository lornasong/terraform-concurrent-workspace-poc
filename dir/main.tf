terraform {
    required_version = ">=0.12"
}

module "local" {
    source = "./modules/local"
    services = { for name in var.service_mapping[var.workspace] : name => var.services[name]}
}
