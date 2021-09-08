---
layout: "pnap"
page_title: "phoenixNAP: pnap_ssh_key"
sidebar_current: "docs-pnap-resource-ssh_key"
description: |-
  Provides a phoenixNAP SSH key resource. This can be used to create, modify, and delete ssh keys.
---

# pnap_ssh_key Resource

Provides a phoenixNAP SSH key resource. This can be used to create,
modify, and delete ssh keys.



## Example Usage

Create a SSH key 

```hcl
# Create a ssh key
resource "pnap_ssh_key" "ssh-key-1" {
    name = "sshkey-1"
    default = false
    key = "ssh-rsa                                                              ABeeeABbb3NzaC1yc2EAAAADAQABAAABAQDF9LdAFElNCi7JoWh6KUcchrJ2Gac1aqGRPpdZNowObpRtmiRCecAMb7bUgNAaNfcmwiQi7tos9TlnFgprIcfMWb8MSs3ABYHmBgqEEt3RWYf0fAc9CsIpJdMCUG28TPGTlRXCEUVNKgLMdcseAlJoGp1CgbHWIN65fB3he3kAZcfpPn5mapV0tsl2p+ZyuAGRYdn5dJv2RZDHUZBkOeUobwsij+weHCKAFmKQKtCP7ybgVHaQjAPrj8MGnk1jBbjDt5ws+Be+9JNjQJre9zCKbAOsIo3i+GcUIkrw5jxPU/RTGlWBcemPaKHdciSzGcjWboapzIy49qypQhZe1U75 user3@172.16.1.106"

}

# Create a server
resource "pnap_server" "Test-Server-1" {
    hostname = "Test-Server-1"
    os = "ubuntu/bionic"
    type = "s1.c1.medium"
    location = "PHX"
    ssh_key_ids = [pnap_ssh_key.ssh-key-1.id]
    ]
    #allowed actions are: reboot, reset, powered-on, powered-off, shutdown
    #action = "powered-on"
}
```

## Argument Reference

The following arguments are supported:

* `default` - (Required) Keys marked as default are always included on server creation and reset unless toggled off in creation/reset request.
* `name` - (Required) Friendly SSH key name to represent an SSH key.
* `key` - (Required) SSH key actual key value.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the SSH Key.
* `default` - Keys marked as default are always included on server creation and reset unless toggled off in creation/reset request.
* `name` - Friendly SSH key name to represent an SSH key.
* `key` - SSH Key value.
* `fingerprint` - SSH key auto-generated SHA-256 fingerprint.
* `createdOn `- Date and time of creation.
* `lastUpdatedOn ` - Date and time of last update.
