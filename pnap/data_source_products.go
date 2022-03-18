package pnap

import (
	"math"
	"strconv"
	"time"

	"github.com/PNAP/go-sdk-helper-bmc/command/billingapi/product"
	"github.com/PNAP/go-sdk-helper-bmc/dto"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProductsRead,

		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"product_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sku_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"products": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plans": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sku": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sku_description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"location": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pricing_model": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"price": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"price_unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"correlated_product_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"package_quantity": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"package_unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"metadata": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ram_in_gb": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"cpu": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cpu_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"cores_per_cpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"cpu_frequency": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"network": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"storage": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceProductsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	query := dto.ProductQuery{}
	query.ProductCode = d.Get("product_code").(string)
	query.ProductCategory = d.Get("product_category").(string)
	query.SKUCode = d.Get("sku_code").(string)
	query.Location = d.Get("location").(string)

	requestCommand := product.NewGetProductsCommand(client, query)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	var products []interface{}
	for _, j := range resp {
		product := make(map[string]interface{})
		product["product_code"] = j.ProductCode
		product["product_category"] = j.ProductCategory
		if len(j.Plans) > 0 {
			var plans []interface{}
			for _, l := range j.Plans {
				pricingPlan := make(map[string]interface{})
				pricingPlan["sku"] = l.SKU
				if l.SKUDescription != "" {
					pricingPlan["sku_description"] = l.SKUDescription
				}
				pricingPlan["location"] = l.Location
				pricingPlan["pricing_model"] = l.PricingModel
				price := math.Round(float64(l.Price)*100) / 100
				pricingPlan["price"] = price
				pricingPlan["price_unit"] = l.PriceUnit
				if l.CorrelatedProductCode != "" {
					pricingPlan["correlated_product_code"] = l.CorrelatedProductCode
				}
				if l.PackageQuantity != 0 {
					pricingPlan["package_quantity"] = int(l.PackageQuantity)
				}
				if l.PackageUnit != "" {
					pricingPlan["package_unit"] = l.PackageUnit
				}
				plans = append(plans, pricingPlan)
			}
			product["plans"] = plans
		}
		if j.Metadata != nil {
			metadata := *j.Metadata
			md := make([]interface{}, 1)
			mdItem := make(map[string]interface{})
			mdItem["ram_in_gb"] = int(metadata.RamInGb)
			mdItem["cpu"] = metadata.CPU
			mdItem["cpu_count"] = int(metadata.CPUCount)
			mdItem["cores_per_cpu"] = int(metadata.CoresPerCPU)
			cpuFreq := math.Round(float64(metadata.CPUFrequency)*100) / 100
			mdItem["cpu_frequency"] = cpuFreq
			mdItem["network"] = metadata.Network
			mdItem["storage"] = metadata.Storage
			md[0] = mdItem
			product["metadata"] = md
		}
		products = append(products, product)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("products", products)
	return nil
}
