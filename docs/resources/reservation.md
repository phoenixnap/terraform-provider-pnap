---
layout: "pnap"
page_title: "phoenixNAP: pnap_reservation"
sidebar_current: "docs-pnap-resource-reservation"
description: |-
  Provides a phoenixNAP reservation resource. This can be used to create and modify reservations.
---

# pnap_reservation Resource

Provides a phoenixNAP reservation resource. This can be used to create and modify reservations.



## Example Usage

Create a reservation 

```hcl
# Create a reservation
resource "pnap_reservation" "Test-Reservation-1" {
    sku = "XXX-XXX-XXX"    
}
```

## Argument Reference

The following arguments are supported:

* `sku` - (Required) The SKU code of product pricing plan.
* `auto_renew` - A flag indicating whether the reservation will auto-renew (default is true, it can only be modified after the creation of resource).
* `auto_renew_disable_reason` - The reason for disabling auto-renewal.


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
* `sku` - The SKU that will be applied to this reservation.
* `price` - Reservation price.
* `price_unit` - The unit to which the price applies.
* `assigned_resource_id` - The resource ID currently being assigned to reservation.
* `next_billing_date` - Next billing date for reservation.
