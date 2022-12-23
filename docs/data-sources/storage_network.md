---
layout: "pnap"
page_title: "phoenixNAP: pnap_storage_network"
sidebar_current: "docs-pnap-datasource-storage_network"
description: |-
  Provides a phoenixNAP Storage Network datasource. This can be used to read storage networks.
---

# pnap_storage_network Datasource

Provides a phoenixNAP Storage Network datasource. This can be used to read storage networks.



## Example Usage

Fetch a storage network by name and show it's volumes

```hcl
# Fetch a storage network
data "pnap_storage_network" "Storage-Network-1" {
    name   = "Storage-1"
}

# Show volumes
output "Volumes" {
    value = data.pnap_storage_network.Storage-Network-1.volumes
}
```

## Argument Reference

The following arguments are supported:

* `name` - The friendly name of this storage network.
* `id` - The storage network identifier.

## Attributes Reference

The following attributes are exported:

* `id` - The storage network identifier.
* `name` - The friendly name of this storage network.
* `description` - The description of this storage network.
* `status` - Storage network's status.
* `location` - The location of this storage network.
* `network_id `- ID of network the storage belongs to.
* `ips` - IP of the storage network
* `created_on` - Date and time when this storage network was created.
* `volumes` - Volume for the storage network.
    * `id` - Volume ID.
    * `name` - Volume friendly name.
    * `description` - Volume description.
    * `path` - Volume's full path. It is in form of `/{volumeId}/pathSuffix`.
    * `path_suffix` - Last part of volume's path.
    * `capacity_in_gb` - Maximum capacity in GB.
    * `used_capacity_in_gb` - Used capacity in GB, updated periodically.
    * `protocol` - File system protocol.
    * `status` - Volume's status.
    * `created_on` - Date and time when this volume was created.
    * `permissions` - Permissions for the volume.
        * `nfs` - NFS specific permissions on the volume.
            * `read_write` - Read/Write access.
            * `read_only` - Read only access.
            * `root_squash` - Root squash permission.
            * `no_squash` - No squash permission.
            * `all_squash` - All squash permission.