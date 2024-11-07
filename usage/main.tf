# An easy way to use this is to write each rule in a yaml file and then load it in your terraform code
# So that you can loop it in your resources and make everything more readable
# Update on the resources to load the yaml files into a fileset and then use the data.local_file to load the content of the files
# This helps avoid issues where your rule might not fit into a huge variable

# locals.tf
locals {
  detection_rules_dir = "./"
  exception_containers_dir = "./"
  exception_items_dir = "./"
}

# detection_rules.tf
# Detection Rules
locals {
  detection_rules_files = fileset(local.detection_rules_dir, "*.yaml")
}

data "local_file" "detection_rules" {
  for_each = toset(local.detection_rules_files)
  filename = "${local.detection_rules_dir}/${each.value}"
}

resource "elastic-siem-detection_detection_rule" "elastic_detection_rule" {
  for_each     = data.local_file.detection_rules
  rule_content = jsonencode(yamldecode(each.value.content))
  depends_on   = [ elastic-siem-detection_exception_container.my_containers ]
}

# exception_containers.tf
# Exception Containers
locals {
  exception_containers_files = fileset(local.exception_containers_dir, "*.yaml")
}

data "local_file" "exception_containers" {
  for_each = toset(local.exception_containers_files)
  filename = "${local.exception_containers_dir}/${each.value}"
}

resource "elastic-siem-detection_exception_container" "elastic_exception_containers" {
  for_each                    = data.local_file.exception_containers
  exception_container_content = jsonencode(yamldecode(each.value.content))
}

# exception_items.tf
# Exception Items
locals {
  exception_items_files = fileset(local.exception_items_dir, "*.yaml")
}

data "local_file" "exception_items" {
  for_each = toset(local.exception_items_files)
  filename = "${local.exception_items_dir}/${each.value}"
}

resource "elastic-siem-detection_exception_item" "elastic_exception_items" {
  for_each               = data.local_file.exception_items
  exception_item_content = jsonencode(yamldecode(each.value.content))
}
