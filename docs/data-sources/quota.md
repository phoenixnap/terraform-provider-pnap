---
layout: "pnap"
page_title: "phoenixNAP: pnap_quota"
sidebar_current: "docs-pnap-datasource-quota"
description: |-
  Provides a phoenixNAP Quota datasource. This can be used to read Quotas.
---

# pnap_quota Datasource

Provides a phoenixNAP Quota datasource. This can be used to read Quotas.



## Example Usage

Fetch a Quota by name and show it's details in alphabetical order

```hcl
# Fetch a Quota
data "pnap_quota" "test" {
    name = "Public IPs"
}

# Show the Quota details
output "quota" {
    value = data.pnap_quota.test
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Quota.
* `id` - The ID of the Quota.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Quota.
* `name` - The name of the Quota.
* `description` - The Quota description.
* `status` - The status of the Quota.
* `limit` - The limit set for the Quota.
* `unit`- Unit of the Quota type.
* `used` - The Quota used expressed as a number.
* `quota_edit_limit_request_details` - List of requests to change the limit on a Quota.
    * `limit` - The new limit that is requested.
    * `reason` - The reason for changing the limit.
    * `requested_on` - The point in time the request was submitted.
