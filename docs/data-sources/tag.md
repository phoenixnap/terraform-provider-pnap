---
layout: "pnap"
page_title: "phoenixNAP: pnap_tag"
sidebar_current: "docs-pnap-datasource-tag"
description: |-
  Provides a phoenixNAP tag datasource. This can be used to read tags.
---

# pnap_tag Datasource

Provides a phoenixNAP tag datasource. This can be used to read tags.



## Example Usage

Fetch a tag by name and show it's details.

```hcl
# Fetch a tag
data "pnap_tag" "test" {
  name   = "tag3"
}

# Show the key
output "details" {
  value = data.pnap_tag.test
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The unique name of the tag.


## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the tag.
* `name` - The name of the tag.
* `values` - The optional values of the tag..
* `description` - The description of the tag.
* `is_billing_tag `- Whether or not to show the tag as part of billing and invoices.
* `resource_assignments ` - The tag's assigned resources.
  * `resource_name` - The resource name.
  * `value` - The value of the tag assigned to the resource.
* `created_by ` - The tag's creator.
