---
layout: "pnap"
page_title: "phoenixNAP: pnap_invoices"
sidebar_current: "docs-pnap-datasource-invoices"
description: |-
  Provides a phoenixNAP invoices datasource. This can be used to read invoices.
---

# pnap_invoices Datasource

Provides a phoenixNAP invoices datasource. This can be used to read invoices.



## Example Usage

Fetch invoices by status and date sent and show their details.

```hcl
# Fetch invoices
data "pnap_invoices" "Query-C" {
  status = "PAID"
  sent_on_from = "2020-04-13T00:00:00.000Z"
  sent_on_to = "2022-04-13T00:00:00.000Z"
}

# Show invoices
output "invoices" {
  value = data.pnap_invoices.Query-C.paginated_invoices
}
```

## Argument Reference

The following arguments are supported:

* `number` - A user-friendly reference number assigned to the invoice.
* `status` - Payment status of the invoice. The following values are allowed: `PAID`, `UNPAID`, `OVERDUE`, `PAYMENT_PROCESSING`
* `sent_on_from` - Minimum value to filter invoices by sent on date.
* `sent_on_to` - Maximum value to filter invoices by sent on date.
* `limit` - The limit of the number of results returned. The number of records returned may be smaller than the limit.
* `offset` - The number of items to skip in the results.
* `sort_field` - If a sort field is requested, pagination will be done after sorting. The following values are allowed: `number`, `sentOn`, `dueDate`, `amount`, `outstandingAmount`.
* `sort_direction` - Sort given field depending on the desired direction. The following values are allowed: `ASC`, `DESC`.


## Attributes Reference

The following attributes are exported:

* `paginated_invoices` - The paginated list of invoices.
    * `limit` - Maximum number of items in the page (actual returned length can be less).
    * `offset` - The number of returned items skipped.
    * `plans` - The total number of records available for retrieval.
    * `results` - The list of invoices.
        * `id` - The unique resource identifier of the invoice.
        * `number` - A user-friendly reference number assigned to the invoice.
        * `currency` - The currency of the invoice.
        * `amount` - The invoice amount.
        * `outstanding_amount` - The invoice outstanding amount.
        * `status` - The status of the invoice.
        * `sent_on` - Date and time when the invoice was sent.
        * `due_date` - Date and time when the invoice payment is due.

