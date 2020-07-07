variable "service_mapping" {
  description = "Map of workspace to service"
  type        = map(list(string))
  default = {
    east = []
    west = []
  }
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

variable "workspace" {
  description = "workspace name"
  type = string
  default = "east"
}
