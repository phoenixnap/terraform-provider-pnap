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

* `hostname` - Server hostname.
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
* `primary_ip_address` - First usable public IP Address.
* `network_type` - The type of network configuration for this server.
* `esxi` - Esxi OS configuration.
    * `datastore_configuration` - Esxi data storage configuration.
        * `datastore_name` - Datastore name.
* `netris_controller` - Netris Controller configuration properties.
    * `host_os` - Host OS on which the Netris Controller is installed.
* `netris_softgate` - Netris Softgate configuration properties.
    * `host_os` - Host OS on which the Netris Softgate is installed.
* `tags` - The tags assigned to the server.
    * `id` - The unique id of the tag.
    * `name` - The name of the tag.
    * `value` - The value of the tag assigned to the server.
    * `is_billing_tag` - Whether or not to show the tag as part of billing and invoices.
    * `created_by` - Who the tag was created by.
* `network_configuration` - Entire network details of bare metal server.
    * `gateway_address` - The address of the gateway assigned to the server.
    * `private_network_configuration` - Private network details of bare metal server.
        * `configuration_type` - Determines the approach for configuring private network(s) for the server being provisioned.
        * `private_networks` - The list of private networks this server is member of.
            * `id` - The network identifier.
            * `ips` - IPs configured on the server.
            * `dhcp` - Determines whether DHCP is enabled for this server.
            * `status_description` - The status of the network.
    * `ip_blocks_configuration` - IP block details of bare metal server.
        * `configuration_type` - Determines the approach for configuring IP blocks for the server being provisioned.
        * `ip_blocks` - The IP blocks assigned to this server.
            * `id` - The IP block's ID.
            * `ips` - The VLAN on which this IP block has been configured within the network switch.
    * `public_network_configuration` - Public network details of bare metal server.
        * `public_networks` - The list of public networks this server is member of.
            * `id` - The network identifier.
            * `ips` - IPs configured on the server.
            * `status_description` - The status of the assignment to the network.
* `storage_configuration` - Storage configuration.
    * `root_partition` - Root partition configuration.
        * `raid` - Software RAID configuration.
        * `size` - The size of the root partition in GB.
* `superseded_by` - Unique identifier of the server to which the reservation has been transferred.
* `supersedes` - Unique identifier of the server from which the reservation has been transferred.
