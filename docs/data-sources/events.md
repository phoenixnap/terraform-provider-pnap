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

Fetch event logs by name and date and show their details.

```hcl
# Fetch events
data "pnap_events" "test" {
  from = "2022-04-12T00:00:00.000Z"
  to = "2023-04-12T00:00:00.000Z"
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

* `from` - From the date and time (inclusive) to filter event log records by.
* `to` - To the date and time (inclusive) to filter event log records by.
* `limit` - Limit the number of records returned.
* `order` - Ordering of the event's time. The following values are allowed: `ASC`, `DESC`. Default value is `ASC`.
* `username` - The username that did the actions.
* `verb` - The HTTP verb corresponding to the action. The following values are allowed: `POST`, `PUT`, `PATCH`, `DELETE`.
* `uri` - The request uri.
* `events` - Block `events` has field `name`.
    * `name` - Event name.


## Attributes Reference

The following attributes are exported:

* `events` - The list of events recorded.
    * `name` - The name of the event.
    * `timestamp` - The UTC time the event initiated.
    * `user_info` - Details related to the user / application performing this request.
        * `account_id` - The BMC account ID.
        * `client_id` - The client ID of the application.
        * `username` - The logged in user or owner of the client application.
