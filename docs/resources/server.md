---
layout: "pnap"
page_title: "PhoenixNAP: pnap_server"
sidebar_current: "docs-pnap-resource-server"
description: |-
  Provides a PhoenixNAP server resource. This can be used to create, modify, and delete servers.
---

# pnap_server Resource

Provides a PhoenixNAP server resource. This can be used to create,
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
    ssh_keys = [
        "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDF9LdAFElNCi7JoWh6KUcchrJ2Gac1aqGRPpdZNowObpRtmiRCecAMb7bUgNAaNfcmwiQi7tos9TlnFgprIcfMWb8MSs3ABYHmBgqEEt3RWYf0fAc9CsIpJdMCUG28TPGTlRXCEUVNKgLMdcseAlJoGp1CgbHWIN65fB3he3kAZcfpPn5mapV0tsl2p+ZyuAGRYdn5dJv2RZDHUZBkOeUobwsij+weHCKAFmKQKtCP7ybgVHaQjAPrj8MGnk1jBbjDt5ws+Be+9JNjQJee9zCKbAOsIo3i+GcUIkrw5jxPU/RTGlWBcemPaKHdciSzGcjWboapzIy49qypQhZe1U75 user2@122.16.1.126"
    
    ]
    #allowed actions are: reboot, reset, powered-on, powered-off, shutdown
    #action = "powered-on"
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Required) Server hostname.
* `description` - Server description.
* `os` - (Required) The server’s OS ID used when the server was created (e.g., ubuntu/bionic, centos/centos7). For a full list of available operating systems visit [API docs](https://developers.phoenixnap.com/docs/bmc/1).
* `type` - (Required) Server type ID. Cannot be changed once a server is created (e.g., s1.c1.small, s1.c1.medium). 
* `location` - (Required) Server Location ID. Cannot be changed once a server is created (e.g., PHX).
* `ssh_keys` - (Required) A list of SSH Keys that will be installed on the Linux server. Must contain at least 1 item.
* `action` - Action to perform on server. Allowed actions are: reboot, reset, powered-on, powered-off, shutdown.

## Attributes Reference

The following attributes are exported:

* `cpu` - A description of the machine's CPU.
* `description` - Server description.
* `hostname ` - Server hostname.
* `id` - The unique identifier of the server.
* `location` - Server Location ID. Cannot be changed once a server is created.
* `os` - The server’s OS ID used when the server was created. 
* `ram` - A description of the machine's RAM.
* `status` - The status of the server.
* `storage`- A description of the machine's storage.
* `type` - Server type ID. Cannot be changed once a server is created. 
* `private_ip_addresses` - Private IP Addresses assigned to server. Must contain at least 1 item. 
* `public_ip_addresses` - Public IP Addresses assigned to server. Must contain at least 1 item 

 
