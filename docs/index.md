---
layout: "pnap"
page_title: "Provider: phoenixNAP"
sidebar_current: "docs-pnap-index"
description: |-
  The phoenixNAP provider is used to interact with the resources supported by PNAP. The provider needs to be configured with the proper credentials before it can be used.
---

# phoenixNAP Provider

The phoenixNAP provider is used to interact with the resources supported by PNAP.
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the resources available.

Be cautious when using the `pnap_server` resource. By default, phoenixNAP invoices hourly per server.

# Authentication

The following authentication methods are supported:

- Static credentials
- Environment variables
- Configuration file

Static credentials can be provided by adding an `client_id` and `client_secret`
in-line in the pnap provider block:

Usage:

```terraform
provider "pnap" {
  client_id = var.client_id
  client_secret = var.client_secret
}
```
You can provide your credentials via the `PNAP_CLIENT_ID` and
`PNAP_CLIENT_SECRET`, environment variables. Note that setting your
phoenixNAP credentials in provider block or environment variables
will override the use of `config_file_path` and configuration file for authentication.

You can use config.yaml file to specify your credentials.
The default location is in the user's home directory. File path on Linux is /.pnap/config.yaml and file path on Windows is \AppData\Roaming\pnap\config.yaml. You can optionally specify a different location in the Terraform configuration by providing the `config_file_path` argument.

Usage:

```terraform
provider "pnap" {
  config_file_path = ""C:\\config_test\\""
}
```

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
* `installDefaultSshKeys` - Whether or not to install SSH keys marked as default in addition to any SSH keys specified in this request.
* `ssh_keys` - A list of SSH Keys that will be installed on the server.
* `ssh_key_ids` - A list of SSH key IDs that will be installed on the server in addition to any SSH keys specified in this request.
* `reservation_id` - Server reservation ID.
* `pricing_model` - Server pricing model. Currently this field should be set to HOURLY, ONE_MONTH_RESERVATION, TWELVE_MONTHS_RESERVATION, TWENTY_FOUR_MONTHS_RESERVATION or THIRTY_SIX_MONTHS_RESERVATION.
* `network_type` - The type of network configuration for this server. Currently this field should be set to PUBLIC_AND_PRIVATE or PRIVATE_ONLY.
* `rdp_allowed_ips` - List of IPs allowed for RDP access to Windows OS. Supported in single IP, CIDR and range format. When undefined, RDP is disabled. To allow RDP access from any IP use 0.0.0.0/0. Must contain at least 1 item.
* `action` - Action to perform on server. Allowed actions are: reboot, reset, powered-on, powered-off, shutdown.
