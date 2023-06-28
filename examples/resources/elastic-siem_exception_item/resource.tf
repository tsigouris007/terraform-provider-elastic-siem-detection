resource "elastic-siem_exception_item" "my_items" {
  exception_item_content = jsonencode(
    {
      "list_id" : "hacker_list",
      "item_id" : "hacker_item",
      "name" : "Catch a hacker\n",
      "description" : "This item excepts the specified user from alerting.\n",
      "namespace_type" : "single",
      "tags" : [
        "MyTag"
      ],
      "type" : "simple",
      "entries" : [
        {
          "field" : "user.name",
          "operator" : "included",
          "type" : "match_any",
          "value" : [
            "hacker"
          ]
        }
      ],
      "expire_time" : "2024-01-01T21:00:00.000Z"
    }
  )

  # Helps syncing between objects
  depends_on = [elastic-siem_exception_container.my_containers]
}
