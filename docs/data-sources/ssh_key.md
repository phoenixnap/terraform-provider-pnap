---
layout: "pnap"
page_title: "phoenixNAP: pnap_ssh_key"
sidebar_current: "docs-pnap-datasource-ssh_key"
description: |-
  Provides a phoenixNAP SSH key datasource. This can be used to read SSH keys.
---

# pnap_ssh_key Datasource

Provides a phoenixNAP SSH key datasource. This can be used to read SSH keys.



## Example Usage

Fetch a SSH key by name and show it's key 

```hcl
# Fetch a SSH key
data "pnap_ssh_key" "test" {
 name   = "test3"
}

# Show the key
output "key" {
  value = data.pnap_ssh_key.test.key
}
```

## Argument Reference

The following arguments are supported:

* `name` - Friendly SSH key name to represent an SSH key.
* `id` - The unique identifier of the SSH Key.


## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the SSH Key.
* `default` - Keys marked as default are always included on server creation and reset unless toggled off in creation/reset request.
* `name` - Friendly SSH key name to represent an SSH key.
* `key` - SSH Key value.
