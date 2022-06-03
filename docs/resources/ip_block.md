---
layout: "pnap"
page_title: "phoenixNAP: pnap_ip_block"
sidebar_current: "docs-pnap-resource-ip-block"
description: |-
  Provides a phoenixNAP IP Block resource. This can be used to create, modify and delete IP Blocks.
---

# pnap_ip_block Resource

Provides a phoenixNAP IP Block resource. This can be used to create, modify and delete IP Blocks.



## Example Usage

Create an IP Block 

```hcl
# Create an IP Block
resource "pnap_ip_block" "ip-block-1" {
    location = "PHX"
    cidr_block_size = "/28"
    description = "IP Block #1 used for publicly accessing server #1."
    tags {
        tag_assignment {
            name = "tag-1"
            value = "PROD"
        }
    }
    tags {
        tag_assignment {
            name = "tag-2"
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) IP Block location ID. Currently this field should be set to `PHX`, `ASH`, `SGP`, `NLD`, `CHI`, `SEA` or `AUS`.
* `cidr_block_size` - (Required) CIDR IP Block Size. Currently this field should be set to either `/31`, `/30`, `/29` or `/28`. For a larger Block Size contact support.
* `description` - Description of the IP Block.
* `tags` - Tags to set to IP Block, if any.
    * `tag_assignment` - Tag request to assign to the IP Block.
        * `name` - (Required) The name of the tag.
        * `value` - The value of the tag assigned to the IP Block.

## Attributes Reference

The following attributes are exported:

* `id` - IP Block identifier.
* `location` - IP Block location ID.
* `cidr_block_size` - CIDR IP Block Size.
* `cidr` - The IP Block in CIDR notation.
* `status` - The status of the IP Block.
* `assigned_resource_id` - ID of the resource assigned to the IP Block.
* `assigned_resource_type `- Type of the resource assigned to the IP Block.
* `description` - Description of the IP Block.
* `tags` - The tags assigned to the IP Block.
    * `tag_assignment` - Tag assigned to the IP Block.
        * `id` - The unique id of the tag.
        * `name` - The name of the tag.
        * `value` - The value of the tag assigned to the IP Block.
        * `is_billing_tag` - Whether or not to show the tag as part of billing and invoices.
        * `created_by` - Who the tag was created by.
* `is_bring_your_own` - True if the IP Block is a "bring your own" block.
* `created_on` - Date and time when the IP Block was created.
