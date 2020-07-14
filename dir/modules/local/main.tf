resource "local_file" "address" {
    for_each = var.services
    content = each.value.address
    filename = "../instance-${var.instance}/${each.value.name}.txt"
}
