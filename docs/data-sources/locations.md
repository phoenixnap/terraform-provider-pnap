---
layout: "pnap"
page_title: "phoenixNAP: pnap_locations"
sidebar_current: "docs-pnap-datasource-locations"
description: |-
  Provides a phoenixNAP locations datasource. This can be used to retrieve locations info.
---

# pnap_locations Datasource

Provides a phoenixNAP locations datasource. This can be used to retrieve locations info.



## Example Usage

Fetch locations by product category and show their details.

```hcl
# Fetch locations
data "pnap_locations" "Query-C" {
  product_category = "BANDWIDTH"
}

# Show locations
output "Locations" {
  value = data.pnap_locations.Query-C.locations
}
```

## Argument Reference

The following arguments are supported:

* `location` - The location code. Currently the following values are allowed: `PHX`, `ASH`, `NLD`, `SGP`, `CHI`, `SEA` and `AUS`.
* `product_category` - The product category. Currently the following values are allowed: `SERVER`, `BANDWIDTH`, `OPERATING_SYSTEM`, `PUBLIC_IP` and `STORAGE`.


## Attributes Reference

The following attributes are exported:

* `locations` - The list of locations found.
    * `location` - The location code.
    * `location_description` - Description of the location.
    * `product_categories` - The list of product categories.
        * `product_category` - The product category.
        * `product_category_description` - Description of the product category.
