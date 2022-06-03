---
layout: "pnap"
page_title: "phoenixNAP: pnap_server"
sidebar_current: "docs-pnap-datasource-server"
description: |-
  Provides a phoenixNAP server datasource. This can be used to read servers.
---

# pnap_server Datasource

Provides a phoenixNAP server datasource. This can be used to read servers.



## Example Usage

Fetch a server data by hostname and show it's primary public IP address

```hcl
# Fetch a server
data "pnap_server" "server_ds" {
 #id   = "60ef31a84bf50b11fc50d7af"
 hostname = "demo-server"
}

# Show IP address
output "server_id" {
  value = data.pnap_server.server_ds.primary_ip_address
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Required) Server hostname.
* `id` - The unique identifier of the server.


## Attributes Reference

The following attributes are exported:



* `hostname ` - Server hostname.
* `id` - The unique identifier of the server.
* `location` - Server Location ID. Cannot be changed once a server is created.
* `os` - The server’s OS ID used when the server was created. 
* `status` - The status of the server.
* `type` - Server type ID. Cannot be changed once a server is created. 
* `private_ip_addresses` - Private IP Addresses assigned to server. Must contain at least 1 item. 
* `public_ip_addresses` - Public IP Addresses assigned to server. Must contain at least 1 item.
* `primary_ip_address` - First usable public IP Addresses.



 
