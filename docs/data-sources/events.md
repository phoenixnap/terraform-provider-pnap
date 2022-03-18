---
layout: "pnap"
page_title: "phoenixNAP: pnap_events"
sidebar_current: "docs-pnap-datasource-events"
description: |-
  Provides a phoenixNAP events datasource. This can be used to read event logs.
---

# pnap_events Datasource

Provides a phoenixNAP events datasource. This can be used to read event logs.



## Example Usage

Fetch event logs by name and show their details.

```hcl
# Fetch events
data "pnap_events" "test" {
  events {
    name = "API.SshKeysUpdate"
  }
}

# Show events
output "logs" {
  value = data.pnap_events.test
}
```

## Argument Reference

The following arguments are supported:

* `events` - (Required) Block `events` has field `name`.
    * `name` - (Required) Event name.


## Attributes Reference

The following attributes are exported:

* `events` - The list of events recorded.
    * `name` - The name of the event.
    * `timestamp` - The UTC time the event initiated.
    * `user_info` - Details related to the user / application.
        * `account_id` - The BMC account ID.
        * `client_id` - The client ID of the application.
        * `username` - The logged in user or owner of the client application.
