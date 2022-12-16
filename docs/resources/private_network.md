---
layout: "pnap"
page_title: "phoenixNAP: pnap_private_network"
sidebar_current: "docs-pnap-resource-private_network"
description: |-
  Provides a phoenixNAP Private Network resource. This can be used to create, modify, and delete private networks.
---

# pnap_private_network Resource

Provides a phoenixNAP Private Network resource. This can be used to create,
modify, and delete private networks.



## Example Usage

```hcl
# Create a private network
resource "pnap_private_network" "Test-Network-33" {
    name = "ttt"
    cidr = "10.0.0.0/24" 
    location = "PHX"
}
resource "pnap_private_network" "Test-Network-44" {
    name = "qqq"
    cidr = "172.16.0.0/24" 
    location = "PHX"
}

# Create a server
resource "pnap_server" "Test-Server-1" {
    hostname = "Test-Server-1"
    os = "ubuntu/bionic"
    type = "s1.c1.medium"
    location = "PHX"
    install_default_ssh_keys = true
    network_configuration {
      private_network_configuration {
        configuration_type = "USER_DEFINED"
        private_networks  {
          server_private_network {
              id = pnap_private_network.Test-Network-33.id
              ips=["10.0.0.12"]
          }
        }
        private_networks  {
          server_private_network {
              id = pnap_private_network.Test-Network-44.id
              ips=["172.16.0.12"]
          }
        }
      }
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The friendly name of this private network. This name should be unique.
* `description` - The description of this private network.
* `location` - (Required) The location of this private network. Supported values are `PHX`, `ASH`, `SGP`, `NLD`, `CHI`, `SEA` and `AUS`.
* `location_default` - Identifies network as the default private network for the specified location. Default value is `false`.
* `vlan_id `- The VLAN that will be assigned to this network.
* `cidr` - (Required) IP range associated with this private network in CIDR notation.

## Attributes Reference

The following attributes are exported:

* `id` - The private network identifier.
* `name` - The friendly name of this private network. This name should be unique.
* `description` - The description of this private network.
* `location` - The location of this private network.
* `location_default` - Identifies network as the default private network for the specified location. Default value is `false`.
* `cidr` - IP range associated with this private network in CIDR notation.
* `vlan_id `- The VLAN of this private network.
* `type` - The type of the private network.
* `servers ` - (Deprecated) List of server details linked to the private network.
    * `id` - The server identifier.
    * `ips` - List of private IPs associated to the server.
* `memberships` - A list of resources that are members of this private network.
    * `resource_id` - The resource identifier.
    * `resource_type` - The resource's type.
    * `ips` - List of public IPs associated to the resource.
* `status` - The status of the private network.
* `created_on` - Date and time when this private network was created.