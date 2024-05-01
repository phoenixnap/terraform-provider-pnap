---
layout: "pnap"
page_title: "phoenixNAP: pnap_transactions"
sidebar_current: "docs-pnap-datasource-transactions"
description: |-
  Provides a phoenixNAP transactions datasource. This can be used to read transactions.
---

# pnap_transactions Datasource

Provides a phoenixNAP transactions datasource. This can be used to read transactions.



## Example Usage

Fetch transactions by date  and show their details.

```hcl
# Fetch transactions
data "pnap_transactions" "Query-C" {
  limit = 25
  from = "2021-04-27T16:24:57.123Z"
  to = "2021-04-29T16:24:57.123Z"
}

# Show transactions
output "transactions" {
  value = data.pnap_transactions.Query-C.paginated_transactions
}
```

## Argument Reference

The following arguments are supported:

* `limit` - The limit of the number of results returned. Default value is `100`.
* `offset` - The number of items to skip in the results. Default value is `0`.
* `sort_direction` - Sort given field depending on the desired direction. The following values are allowed: `ASC`, `DESC`. Default sorting is descending.
* `sort_field` - If a sort field is requested, pagination will be done after sorting. The following values are allowed: `date`, `amount`, `status`, `cardPaymentMethodDetails.cardType`, `cardPaymentMethodDetails.lastFourDigits`, `metadata.invoiceId`, `metadata.isAutoCharge`. Default sorting is by date.
* `from` - From the date and time (inclusive) to filter transactions by.
* `to` - To the date and time (inclusive) to filter transactions by.
* `id` - The unique identifier of the transaction.


## Attributes Reference

The following attributes are exported:

* `paginated_transactions` - The paginated list of transactions.
    * `limit` - Maximum number of items in the page (actual returned length can be less).
    * `offset` - The number of returned items skipped.
    * `total` - The total number of records available for retrieval.
    * `results` - The list of transactions.
        * `id` - The transaction ID.
        * `status` - The status of the transaction.
        * `details` - Details about the transaction. Contains failure reason in case of failed transactions.
        * `amount` - The transaction amount.
        * `currency` - The transaction currency.
        * `date` - Date and time when transaction was created.
        * `metadata` - Transaction's metadata.
            * `invoice_id` - The invoice ID that this transaction pertains to.
            * `invoice_number` - A user-friendly reference number assigned to the invoice that this transaction pertains to.
            * `is_auto_charge` - Whether this transaction was triggered by an auto charge or not.
        * `card_payment_method_details` - Card payment details of a transaction.
            * `card_type` - The Card Type. Supported Card Types include: VISA, MASTERCARD, DISCOVER, JCB & AMEX.
            * `last_four_digits` - The last four digits of the card number.
