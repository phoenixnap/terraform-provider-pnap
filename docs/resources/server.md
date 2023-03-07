---
layout: "pnap"
page_title: "phoenixNAP: pnap_server"
sidebar_current: "docs-pnap-resource-server"
description: |-
  Provides a phoenixNAP server resource. This can be used to create, modify, and delete servers.
---

# pnap_server Resource

Provides a phoenixNAP server resource. This can be used to create,
modify, and delete servers.



## Example Usage

Create a server

```hcl
# Create a server
resource "pnap_server" "Test-Server-1" {
    hostname = "Test-Server-1"
    os = "ubuntu/bionic"
    type = "s1.c1.medium"
    location = "PHX"
    install_default_ssh_keys = true
    ssh_keys = [
        "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDF9LdAFElNCi7JoWh6KUcchrJ2Gac1aqGRPpdZNowObpRtmiRCecAMb7bUgNAaNfcmwiQi7tos9TlnFgprIcfMWb8MSs3ABYHmBgqEEt3RWYf0fAc9CsIpJdMCUG28TPGTlRXCEUVNKgLMdcseAlJoGp1CgbHWIN65fB3he3kAZcfpPn5mapV0tsl2p+ZyuAGRYdn5dJv2RZDHUZBkOeUobwsij+weHCKAFmKQKtCP7ybgVHaQjAPrj8MGnk1jBbjDt5ws+Be+9JNjQJee9zCKbAOsIo3i+GcUIkrw5jxPU/RTGlWBcemPaKHdciSzGcjWboapzIy49qypQhZe1U75 user2@122.16.1.126"
    
    ]
    cloud_init {
        user_data = filebase64("~/terraform-provider-pnap/create-folder.txt")
    }
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
    #pricing_model = "ONE_MONTH_RESERVATION"
    #allowed actions are: reboot, reset, powered-on, powered-off, shutdown
    #action = "powered-on"
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Required) Server hostname.
* `description` - Server description.
* `os` - (Required) The server’s OS ID used when the server was created (e.g., ubuntu/bionic, centos/centos7). For a full list of available operating systems visit [API docs](https://developers.phoenixnap.com/docs/bmc/1).
* `type` - (Required) Server type ID. Cannot be changed once a server is created (e.g., s1.c1.small, s1.c1.medium). For a full list of available types visit [API docs](https://developers.phoenixnap.com/docs/bmc/1). 
* `location` - (Required) Server Location ID. Cannot be changed once a server is created (e.g., PHX). For a full list of available locations visit [API docs](https://developers.phoenixnap.com/docs/bmc/1)
* `installDefaultSshKeys` - Whether or not to install SSH keys marked as default in addition to any SSH keys specified in this request.
* `ssh_keys` - A list of SSH Keys that will be installed on the server.
* `ssh_key_ids` - A list of SSH key IDs that will be installed on the server in addition to any SSH keys specified in this request.
* `reservation_id` - Server reservation ID.
* `pricing_model` - Server pricing model. Currently this field should be set to HOURLY, ONE_MONTH_RESERVATION, TWELVE_MONTHS_RESERVATION, TWENTY_FOUR_MONTHS_RESERVATION or THIRTY_SIX_MONTHS_RESERVATION.
* `network_type` - The type of network configuration for this server. Currently this field should be set to PUBLIC_AND_PRIVATE or PRIVATE_ONLY.
* `rdp_allowed_ips` - List of IPs allowed for RDP access to Windows OS. Supported in single IP, CIDR and range format. When undefined, RDP is disabled. To allow RDP access from any IP use 0.0.0.0/0. Must contain at least 1 item.
* `management_access_allowed_ips` - Define list of IPs allowed to access the Management UI. Supported in single IP, CIDR and range format. When undefined, Management UI is disabled.Must contain at least 1 item.
* `install_os_to_ram` - If true, OS will be installed to and booted from the server's RAM. On restart RAM OS will be lost and the server will not be reachable unless a custom bootable OS has been deployed. Only supported for ubuntu/focal. Default value is `false`.
* `cloud_init` - Cloud-init configuration details.
* `tags` - Tags to set to server, if any.
* `network_configuration` - Entire network details of bare metal server. Structure is documented below.
* `action` - Action to perform on server. Allowed actions are: reboot, reset (deprecated), powered-on, powered-off, shutdown.
* `force` - Query parameter controlling advanced features availability. Currently applicable for networking. It is advised to use with caution since it might lead to unhealthy setups.

The `cloud_init` block has one field:

* `user_data` - User data for the [cloud-init](https://cloudinit.readthedocs.io/en/latest/) configuration in base64 encoding. NoCloud format is supported. Follow the [instructions](https://phoenixnap.com/kb/bmc-cloud-init) on how to provision a server using cloud-init. Only ubuntu/bionic and ubuntu/focal are supported.


The `tags` block has field `tag_assignment`.
The `tag_assignment` block has 2 fields:

* `name` - (Required) The name of the tag.
* `value` - The value of the tag assigned to the IP Block.


The `network_configuration` block has 4 fields: `gateway_address`, `private_network_configuration`, `ip_blocks_configuration` and `public_network_configuration`.

* `gateway_address` -The address of the gateway assigned / to assign to the server. When used as part of request body, IP address has to be part of private/public network assigned to this server.

The `private_network_configuration` is the second field of the `network_configuration` block. 
The `private_network_configuration` block has 3 fields:

* `gateway_address` - (Deprecated) The address of the gateway assigned / to assign to the server. When used as part of request body, it has to match one of the IP addresses used in the existing assigned private networks for the relevant location. Deprecated in favour of a common gateway address across all networks available under `network_configuration`.
* `configuration_type` - Determines the approach for configuring IP blocks for the server being provisioned. Currently this field should be set to `USE_OR_CREATE_DEFAULT` or `USER_DEFINED`. Default value is `USE_OR_CREATE_DEFAULT`.
* `private_networks` - The list of private networks this server is member of. When this field is part of request body, it'll be used to specify the private networks to assign to this server upon provisioning. Used alongside the `USER_DEFINED` configurationType.

The `private_networks` block has field `server_private_network`.
The `server_private_network` block has 3 fields:

* `id` - (Required) The network identifier.
* `ips` - IPs to configure/configured on the server. Should be null or empty list if DHCP is true. Must contain at most 10 items.
* `dhcp` - Determines whether DHCP is enabled for this server. Should be false if ips is not an empty list. Not supported for proxmox OS. Default value is `false`.

The `ip_blocks_configuration` is the third field of the `network_configuration` block.
The `ip_blocks_configuration` block has 2 fields:

* `configuration_type` - Determines the approach for configuring IP blocks for the server being provisioned. If `PURCHASE_NEW` is selected, the smallest supported range, depending on the operating system, is allocated to the server. The following values are allowed: `PURCHASE_NEW`, `USER_DEFINED`, `NONE`. Default value is `PURCHASE_NEW`.
* `ip_blocks` - Used to specify the previously purchased IP blocks to assign to this server upon provisioning. Used alongside the `USER_DEFINED` configurationType. Must contain at most 1 item.

The `ip_blocks` block has field `server_ip_block`.
The `server_ip_block` block has 2 fields:

* `id` - (Required) The IP Block's ID.
* `vlan_id` - The VLAN on which this IP block has been configured within the network switch.

The `public_network_configuration` is the fourth field of the `network_configuration` block. 
The `public_network_configuration` block has field `public_networks`:

The `public_networks` block has field `server_public_network`.
The `server_public_network` block has 2 fields:

* `id` - (Required) The network identifier.
* `ips` - (Required) IPs to configure on the server. IPs must be within the network's range. Must contain at least 1 item.

## Attributes Reference

The following attributes are exported:

* `cpu` - A description of the machine CPU.
* `cpu_count` - The number of CPUs available in the system.
* `cores_per_cpu` - The number of physical cores present on each CPU.
* `cpu_frequency_in_ghz` - The CPU frequency in GHz.
* `description` - Server description.
* `hostname ` - Server hostname.
* `id` - The unique identifier of the server.
* `location` - Server Location ID. Cannot be changed once a server is created.
* `os` - The server’s OS ID used when the server was created.
* `ram` - A description of the machine RAM.
* `status` - The status of the server.
* `storage`- A description of the machine storage.
* `type` - Server type ID. Cannot be changed once a server is created. 
* `private_ip_addresses` - Private IP Addresses assigned to server. Must contain at least 1 item. 
* `public_ip_addresses` - Public IP Addresses assigned to server. Must contain at least 1 item.
* `reservation_id` - The reservation reference id if any.
* `pricing_model` - The pricing model this server is being billed.
* `password` - Password set for user Admin on Windows server which will only be returned in response to provisioning a server.
* `network_type` - The type of network configuration for this server. 
* `cluster_id` - The cluster reference id if any.
* `management_ui_url` - The URL of the management UI which will only be returned in response to provisioning a server.
* `root_password` - Password set for user root on an ESXi server which will only be returned in response to provisioning a server.
* `management_access_allowed_ips` - A list of IPs allowed to access the Management UI. Supported in single IP, CIDR and range format. When undefined, Management UI is disabled.
* `install_os_to_ram` - If true, OS will be installed to and booted from the server's RAM. On restart RAM OS will be lost and the server will not be reachable unless a custom bootable OS has been deployed. Only supported for ubuntu/focal. Default value is `false`.
* `cloud_init` - Cloud-init configuration details.
* `tags` - The tags assigned if any.
* `network_configuration` - Entire network details of bare metal server.
* `provisioned_on` - Date and time when server was provisioned.

The `cloud_init` block has one field:
* `user_data` - User data for the cloud-init configuration in base64 encoding.

The `tags` block has field `tag_assignment`.
The `tag_assignment` block has 5 fields:
* `id` - The unique id of the tag.
* `name` - The name of the tag.
* `value` - The value of the tag assigned to the server.
* `is_billing_tag` - Whether or not to show the tag as part of billing and invoices.
* `created_by` - Who the tag was created by.
