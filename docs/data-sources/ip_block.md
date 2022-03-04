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

* `cidr` - (Required) The IP Block in CIDR notation.


## Attributes Reference

The following attributes are exported:

* `location` - IP Block location ID.
* `cidr_block_size` - CIDR IP Block Size.
* `cidr` - The IP Block in CIDR notation.
* `status` - The status of the IP Block.
* `assigned_resource_id` - ID of the resource assigned to the IP Block.
* `assigned_resource_type `- Type of the resource assigned to the IP Block.
