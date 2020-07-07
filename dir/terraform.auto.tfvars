service_mapping = {
  west = ["web1", "web2"]
  east = ["web1"]
}

services = {
  web1: {
    name = "web1"
    address = "192.0.0.1:8000"
  },
  web2: {
    name = "web2"
    address = "192.0.0.13:5000"
  }
}
