---
layout: "pnap"
page_title: "phoenixNAP: pnap_ip_block"
sidebar_current: "docs-pnap-datasource-ip-block"
description: |-
  Provides a phoenixNAP IP Block datasource. This can be used to read IP Blocks.
---

# pnap_ip_block Datasource

Provides a phoenixNAP IP Block datasource. This can be used to read IP Blocks.



## Example Usage

Fetch an IP Block by CIDR and show it's details in alphabetical order

```hcl
# Fetch an IP Block
data "pnap_ip_block" "test" {
    cidr = "1.1.1.0/31"
}

# Show the IP Block details
output "ip-block" {
    value = data.pnap_ip_block.test
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - The IP Block in CIDR notation.
* `id` - The IP Block identifier.


## Attributes Reference

The following attributes are exported:

* `id` - The IP Block identifier.
* `location` - IP Block location ID.
* `cidr_block_size` - CIDR IP Block Size.
* `cidr` - The IP Block in CIDR notation.
* `ip_version` - The IP Version of the block.
* `status` - The status of the IP Block.
* `assigned_resource_id` - ID of the resource assigned to the IP Block.
* `assigned_resource_type `- Type of the resource assigned to the IP Block.
* `description` - Description of the IP Block.
* `tags` - The tags assigned to the IP Block.
    * `id` - The unique id of the tag.
    * `name` - The name of the tag.
    * `value` - The value of the tag assigned to the IP Block.
    * `is_billing_tag` - Whether or not to show the tag as part of billing and invoices.
    * `created_by` - Who the tag was created by.
* `is_bring_your_own` - True if the IP Block is a "bring your own" block.
* `created_on` - Date and time when the IP Block was created.