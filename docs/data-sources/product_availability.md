---
layout: "pnap"
page_title: "phoenixNAP: pnap_product_availability"
sidebar_current: "docs-pnap-datasource-product_availability"
description: |-
  Provides a phoenixNAP product availability datasource. This can be used to read product availabilities.
---

# pnap_product_availability Datasource

Provides a phoenixNAP product availability datasource. This can be used to read product availabilities.



## Example Usage

Fetch product availabilities by product category, product codes and locations.

```hcl
# Fetch product availabilities
data "pnap_product_availability" "Query-1" {
  product_category = ["SERVER"]
  product_code = ["s1.c1.small", "s1.c1.medium"]
  location = ["PHX", "ASH"]
}

# Show product availabilities
output "Availabilities" {
  value = data.pnap_product_availability.Query-1.product_availabilities
}
```

## Argument Reference

The following arguments are supported:

* `product_category` - Product category. Currently only `SERVER` category is supported.
* `product_code` - The code identifying the product. This code has significance across all locations.
* `show_only_min_quantity_available` - Show only locations where product with requested quantity is available or all locations where product is offered. Default value is `true`.
* `location` - The location code. Currently the following values are allowed: `PHX`, `ASH`, `NLD`, `SGP`, `CHI`, `SEA` and `AUS`.
* `solution` - Currently only the following value is allowed: `SERVER_RANCHER`.
* `min_quantity` - Minimal quantity of product needed. Minimum, maximum and default values might differ for different products. For servers, they are 1, 10 and 1 respectively.


## Attributes Reference

The following attributes are exported:

* `product_availabilities` - List of product availabilities.
    * `product_code` - The code identifying the product.
    * `product_category` - The product category.
    * `location_availability_details` - Infos about location, solutions and availability for a product.
        * `location` - The code identifying the location.
        * `min_quantity_requested` - Requested quantity.
        * `min_quantity_available` - Is product available in specific location for requested quantity.
        * `available_quantity` - Total available quantity of product in specific location. Max value is 10.
        * `solutions` - Solutions supported in specific location for a product.
       