---
layout: "pnap"
page_title: "phoenixNAP: pnap_ip_block"
sidebar_current: "docs-pnap-resource-ip-block"
description: |-
  Provides a phoenixNAP IP Block resource. This can be used to create and delete IP Blocks.
---

# pnap_ip_block Resource

Provides a phoenixNAP IP Block resource. This can be used to create and delete IP Blocks.



## Example Usage

Create an IP Block 

```hcl
# Create an IP Block
resource "pnap_ip_block" "ip-block-1" {
    location = "PHX"
    cidr_block_size = "/30"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) IP Block location ID. Currently this field should be set to PHX, ASH, SGP, NLD, CHI or SEA.
* `cidr_block_size` - (Required) CIDR IP Block Size. Currently this field should be set to either /31, /30, /29 or /28.

## Attributes Reference

The following attributes are exported:

* `location` - IP Block location ID.
* `cidr_block_size` - CIDR IP Block Size.
* `cidr` - The IP Block in CIDR notation.
* `status` - The status of the IP Block.
* `assigned_resource_id` - ID of the resource assigned to the IP Block.
* `assigned_resource_type `- Type of the resource assigned to the IP Block.
