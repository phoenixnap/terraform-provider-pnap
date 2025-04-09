---
layout: "pnap"
page_title: "phoenixNAP: pnap_public_network"
sidebar_current: "docs-pnap-resource-public_network"
description: |-
  Provides a phoenixNAP Public Network resource. This can be used to create, modify, and delete public networks.
---

# pnap_public_network Resource

Provides a phoenixNAP Public Network resource. This can be used to create,
modify, and delete public networks.



## Example Usage

```hcl
# Create a public network
resource "pnap_public_network" "Public-Network-1" {
    name = "PubNet1"
    description = "First public network."
    location = "PHX"
    ip_blocks {
        public_network_ip_block {
            id = "60473a6115e34466c9f8f083"
        }
    }
    ip_blocks {
        public_network_ip_block {
            id = "616e6ec6d66b406a45ab8797"
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The friendly name of this public network. This name should be unique.
* `description` - The description of this public network.
* `location` - (Required) The location of this public network. Supported values are `PHX`, `ASH`, `SGP`, `NLD`, `CHI`, `SEA` and `AUS`.
* `vlan_id `- The VLAN that will be assigned to this network.
* `ip_blocks` - A list of IP Blocks that will be associated with this public network (10 items at most).
    * `public_network_ip_block` - The assigned IP Block to the public network.
        * `id` - The IP Block identifier.
* `ra_enabled` - Boolean indicating whether Router Advertisement is enabled. Only applicable for Network with IPv6 Blocks.

## Attributes Reference

The following attributes are exported:

* `id` - The public network identifier.
* `vlan_id `- The VLAN of this public network.
* `memberships` - A list of resources that are members of this public network.
    * `resource_id` - The resource identifier.
    * `resource_type` - The resource's type.
    * `ips` - List of public IPs associated to the resource.
* `name` - The friendly name of this public network.
* `location` - The location of this public network.
* `description` - The description of this public network.
* `status` - The status of the public network.
* `created_on` - Date and time when this public network was created.
* `ip_blocks` - A list of IP Blocks that are associated with this public network.
    * `public_network_ip_block` - The assigned IP Block to the public network.
        * `id` - The IP Block identifier.
        * `cidr` - The CIDR notation of the IP block.
        * `used_ips_count` - The number of IPs used in the IP block.
* `ra_enabled` - Boolean indicating whether Router Advertisement is enabled. Only applicable for Network with IPv6 Blocks.
