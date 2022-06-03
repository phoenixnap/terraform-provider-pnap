---
layout: "pnap"
page_title: "phoenixNAP: pnap_public_network"
sidebar_current: "docs-pnap-datasource-public_network"
description: |-
  Provides a phoenixNAP Public Network datasource. This can be used to read public networks.
---

# pnap_public_network Datasource

Provides a phoenixNAP Public Network datasource. This can be used to read public networks.



## Example Usage

Fetch a public network by name and show it's IP Blocks 

```hcl
# Fetch a public network
data "pnap_public_network" "Public-Network-1" {
    name   = "PubNet1"
}

# Show IP Blocks
output "IP-Blocks" {
    value = data.pnap_public_network.Public-Network-1.ip_blocks
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The friendly name of this public network.

## Attributes Reference

The following attributes are exported:

* `id` - The public network identifier.
* `vlan_id `- The VLAN of this public network.
* `memberships` - A list of resources that are members in this public network.
    * `resource_id` - The resource identifier.
    * `resource_type` - The resource's type.
    * `ips` - List of public IPs associated to the resource.
* `name` - The friendly name of this public network.
* `location` - The location of this public network.
* `description` - The description of this public network.
* `created_on` - Date and time when this public network was created.
* `ip_blocks` - A list of IP Blocks that are associated with this public network.
    * `id` - The IP Block identifier.