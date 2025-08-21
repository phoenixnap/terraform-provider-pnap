---
layout: "pnap"
page_title: "phoenixNAP: pnap_reservation"
sidebar_current: "docs-pnap-datasource-reservation"
description: |-
  Provides a phoenixNAP reservation datasource. This can be used to read reservation details.
---

# pnap_reservation Datasource

Provides a phoenixNAP reservation datasource. This can be used to read reservation details.



## Example Usage

Fetch a reservation by ID or SKU and show it's details in alphabetical order. 

```hcl
# Fetch a reservation
data "pnap_reservation" "test" {
  id = "e6afba51-7de8-4080-83ab-0f915570659c"
  sku = "XXX-XXX-XXX"
}

# Show the reservation details
output "reservation" {
  value = data.pnap_reservation.test
}
```

## Argument Reference

The following arguments are supported:

* `id` - The reservation identifier.
* `sku` - The SKU code of product pricing plan.


## Attributes Reference

The following attributes are exported:

* `id` - The reservation identifier.
* `product_code` - The code identifying the product. This code has significance across all locations.
* `product_category` - The product category.
* `location` - The location code.
* `reservation_model` - The reservation model.
* `reservation_state` - Reservation state.
* `initial_invoice_model` - Reservations created with initial invoice model ON_CREATION will be invoiced on same date when reservation is created. Reservation created with CALENDAR_MONTH initial invoice model will be invoiced at the begining of next month.
* `quantity` - Represents the quantity.
  * `quantity` - Quantity size.
  * `unit` - Quantity unit.
* `start_date_time` - The point in time (in UTC) when the reservation starts.
* `end_date_time` - The point in time (in UTC) when the reservation ends.
* `last_renewal_date_time` - The point in time (in UTC) when the reservation was renewed last.
* `next_renewal_date_time` - The point in time (in UTC) when the reservation will be renewed if auto renew is set to true.
* `auto_renew` - A flag indicating whether the reservation will auto-renew (default is true, it can only be modified after the creation of resource).
* `sku` - The SKU applied to this reservation.
* `price` - Reservation price.
* `price_unit` - The unit to which the price applies.
* `assigned_resource_id` - The resource ID currently being assigned to reservation.
* `next_billing_date` - Next billing date for reservation.
