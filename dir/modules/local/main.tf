resource "local_file" "address" {
    for_each = var.services
    content = each.value.address
    filename = "../instance-${terraform.workspace}/${each.value.name}.txt"
}

variable "services" {
  description = "Simplified version of a monitored service"
  type = map(object({
    # Name of the service
    name = string
    # List of addresses for instances of the service
    address = string
  }))
}
