---
layout: "pnap"
page_title: "phoenixNAP: pnap_storage_network"
sidebar_current: "docs-pnap-resource-storage_network"
description: |-
  Provides a phoenixNAP Storage Network resource. This can be used to create, modify and delete storage networks.
---

# pnap_storage_network Resource

Provides a phoenixNAP Storage Network resource. This can be used to create, modify and delete storage networks.



## Example Usage

```hcl
# Create a storage network and volume
resource "pnap_storage_network" "Storage-Network-1" {
    name = "Storage-1"
    description = "First storage network."
    location = "PHX"
    volumes {
        volume {
            name = "Volume-1"
            path_suffix = "/shared-docs"
            capacity_in_gb = 1000
            tags {
                tag_assignment {
                    name = "tag-1"
                    value = "PROD"
                }
            }
            tags {
                tag_assignment {
                    name = "tag-2"
                }
            }
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The friendly name of this storage network. This name should be unique.
* `description` - The description of this storage network.
* `location` - (Required) The location of this storage network. Currently this field should be set to `PHX` or `ASH`.
* `client_vlan` - Custom Client VLAN that the Storage Network will be set to.
* `volumes` - (Required) Volumes to be created alongside storage. Currently only 1 volume is supported (must contain exactly one item).
    * `volume` - (Required) Volume to be created alongside storage.
        * `name` - (Required) Volume friendly name.
        * `description` - Volume description.
        * `path_suffix` - Last part of volume's path.
        * `capacity_in_gb` - (Required) Capacity of volume in GB. Currently only whole numbers and multiples of 1000 GB are supported.
        * `tags` - Tags to set to the volume.
            * `tag_assignment` - Tag to set to the volume.
                * `name` - (Required) The name of the tag.
                * `value` - The value of the tag assigned to the volume.

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
* `volumes` - Volumes for the storage network.
    * `volume` - Volume for the storage network.
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
        * `tags` - The tags assigned to the volume.
            * `tag_assignment` - Tag assigned to the volume.
                * `id` - The unique id of the tag.
                * `name` - The name of the tag.
                * `value` - The value of the tag assigned to the volume.
                * `is_billing_tag` - Whether or not to show the tag as part of billing and invoices.
                * `created_by` - Who the tag was created by.
