---
layout: "pnap"
page_title: "phoenixNAP: pnap_private_network"
sidebar_current: "docs-pnap-datasource-private_network"
description: |-
  Provides a phoenixNAP Private Network datasource. This can be used to read private networks.
---

# pnap_private_network Datasource

Provides a phoenixNAP Private Network datasource. This can be used to read private networks.



## Example Usage

Fetch a private network by name and show it's servers 

```hcl
# Fetch a private network
data "pnap_private_network" "Test-Network-44" {
    name   = "qqq"
}

# Show servers
output "servers" {
    value = data.pnap_private_network.Test-Network-44.servers
}
```

## Argument Reference

The following arguments are supported:

* `name` - The friendly name of this private network.
* `id` - The private network identifier.

## Attributes Reference

The following attributes are exported:

* `id` - The private network identifier.
* `name` - The friendly name of this private network.
* `description` - The description of this private network.
* `location` - The location of this private network.
* `location_default` - Identifies network as the default private network for the specified location. Default value is `false`.
* `cidr` - IP range associated with this private network in CIDR notation.
* `vlan_id `- The VLAN of this private network.
* `type` - The type of the private network.
* `servers ` - (Deprecated) List of servers' details linked to the private network.
    * `id` - The server identifier.
    * `ips` - List of private IPs associated to the server.
* `memberships` - A list of resources that are members of this private network.
    * `resource_id` - The resource identifier.
    * `resource_type` - The resource's type.
    * `ips` - List of public IPs associated to the resource.
* `status` - The status of the private network.
* `created_on` - Date and time when this private network was created.