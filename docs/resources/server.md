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
* `action` - Action to perform on server. Allowed actions are: reboot, reset, powered-on, powered-off, shutdown.


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
* `provisioned_on` - Date and time when server was provisioned.


 
