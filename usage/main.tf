# An easy way to use this is to write each rule in a yaml file and then load it to a map
# So that you can loop it in your resources and make everything more readable

# locals.tf
locals {
  detection_rules_dir   = "${path.module}/detection-rules"
  detection_rules_files = fileset(local.detection_rules_dir, "*.yaml")
  detection_rules_data  = [for f in local.detection_rules_files : {
    filename = f
    content  = jsonencode(yamldecode(file("${local.detection_rules_dir}/${f}")))
  }]

  exception_containers_dir   = "${path.module}/exception-containers"
  exception_containers_files = fileset(local.exception_containers_dir, "*.yaml")
  exception_containers_data  = [for f in local.exception_containers_files : {
    filename = f
    content  = jsonencode(yamldecode(file("${local.exception_containers_dir}/${f}")))
  }]

  exception_items_dir   = "${path.module}/exception-items"
  exception_items_files = fileset(local.exception_items_dir, "*.yaml")
  exception_items_data  = [for f in local.exception_items_files : {
    filename = f
    content  = jsonencode(yamldecode(file("${local.exception_items_dir}/${f}")))
  }]
}

# resources.tf
resource "elastic-siem_detection_rule" "my_rules" {
  for_each     = { for rule in local.detection_rules_data : rule.filename => rule }

  rule_content = each.value.content

  depends_on = [ elastic-siem_exception_container.my_containers ]
}

resource "elastic-siem_exception_container" "my_containers" {
  for_each = { for container in local.exception_containers_data : container.filename => container }

  exception_container_content = each.value.content
}

resource "elastic-siem_exception_item" "my_items" {
  for_each = { for exception in local.exception_items_data : exception.filename => exception }

  exception_item_content = each.value.content

  depends_on = [ elastic-siem_exception_container.my_containers ]
}
