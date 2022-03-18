---
layout: "pnap"
page_title: "phoenixNAP: pnap_tag"
sidebar_current: "docs-pnap-resource-tag"
description: |-
  Provides a phoenixNAP tag resource. This can be used to create, modify, and delete tags.
---

# pnap_tag Resource

Provides a phoenixNAP tag resource. This can be used to create, modify, and delete tags.



## Example Usage

Create a tag 

```hcl
# Create a tag
resource "pnap_tag" "tag-1" {
    name = "tag-1"
    is_billing_tag = false    
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The unique name of the tag.
* `description` - The description of the tag.
* `is_billing_tag` - (Required) Whether or not to show the tag as part of billing and invoices.


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
