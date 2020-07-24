---
layout: "pnap"
page_title: "Provider: PhoenixNAP"
sidebar_current: "docs-pnap-index"
description: |-
  The PhoenixNAP provider is used to interact with the resources supported by PNAP. The provider needs to be configured with the proper credentials before it can be used.
---

# PhoenixNAP Provider

The PhoenixNAP provider is used to interact with the resources supported by PNAP.
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the resources available.

Be cautious when using the `pnap_server` resource. PhoenixNAP invoices hourly per server.

# Authentication

The config.yaml configuration file is required for authentication. The file should be located in the user's home directory. File path on Linux is /.pnap/config.yaml and file path on Windows is \AppData\Roaming\pnap\config.yaml

The following shows a sample config file:

			# ===================================================== 
			#Sample yaml config file 
			# =====================================================
			# Authentication
			clientId: <enter your client id>
			clientSecret: <enter your client secret>


## Example Usage

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
