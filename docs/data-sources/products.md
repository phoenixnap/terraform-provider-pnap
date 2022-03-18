---
layout: "pnap"
page_title: "phoenixNAP: pnap_products"
sidebar_current: "docs-pnap-datasource-products"
description: |-
  Provides a phoenixNAP products datasource. This can be used to read products.
---

# pnap_products Datasource

Provides a phoenixNAP products datasource. This can be used to read products.



## Example Usage

Fetch products by product category and show their details.

```hcl
# Fetch products
data "pnap_products" "Query-B" {
  product_category = "BANDWIDTH"
}

# Show products
output "Products" {
  value = data.pnap_products.Query-B.products
}
```

## Argument Reference

The following arguments are supported:

* `product_code` - The code identifying the product. This code has significance across all locations.
* `product_category` - The product category.
* `sku_code` - The SKU identifier.
* `location` - The location code. Currently the following values are allowed: `PHX`, `ASH`, `NLD`, `SGP`, `CHI`, `SEA`, `AUS` and `GLOBAL`.


## Attributes Reference

The following attributes are exported:

* `products` - The list of products recorded.
    * `product_code` - The code identifying the product.
    * `product_category` - The product category.
    * `plans` - The pricing plans available for this product.
        * `sku` - The SKU identifying the pricing plan.
        * `sku_description` - Description of the pricing plan.
        * `location` - The code identifying the location.
        * `pricing_model` - The pricing model.
        * `price` - Price per unit.
        * `price_unit` - The unit to which the price applies.
        * `correlated_product_code` - Product code of the correlated product.
        * `package_quantity` - Package size per month.
        * `package_unit` - Package size unit.
    * `metadata` - Details of the server product.
        * `ram_in_gb` - RAM in GB.
        * `cpu` - CPU name.
        * `cpu_count` - Number of CPUs.
        * `cores_per_cpu` - The number of physical cores present on each CPU.
        * `cpu_frequency` - CPU frequency in GHz.
        * `network` - Server network.
        * `storage` - Server storage.
