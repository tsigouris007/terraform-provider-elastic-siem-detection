resource "elastic-siem-detection_exception_container" "my_containers" {
  exception_container_content = jsonencode(
    {
      "list_id" : "hacker_list",
      "description" : "Catch a hacker container.\n",
      "name" : "Catch a hacker container\n",
      "tags" : [
        "MyTag"
      ],
      "type" : "detection",
      "namespace_type" : "single"
    }
  )
}
