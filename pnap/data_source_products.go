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
									"applicable_discounts": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"discounted_price": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"discount_details": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"code": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"value": {
																Type:     schema.TypeFloat,
																Computed: true,
															},
															"coupon_code": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"correlated_product_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"package_quantity": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"package_unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"package_details": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"package_quantity": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"min": {
																Type:     schema.TypeFloat,
																Computed: true,
															},
															"max": {
																Type:     schema.TypeFloat,
																Computed: true,
															},
														},
													},
												},
												"package_unit": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
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
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"cpu": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cpu_count": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"cores_per_cpu": {
										Type:     schema.TypeFloat,
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
									"gpu_configurations": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"count": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"name": {
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
	query.SkuCode = d.Get("sku_code").(string)
	query.Location = d.Get("location").(string)

	requestCommand := product.NewGetProductsCommand(client, query)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	products := make([]interface{}, 0, len(resp))
	for _, j := range resp {
		product := make(map[string]interface{})
		product["product_code"] = j.ProductCode
		product["product_category"] = j.ProductCategory
		if len(j.Plans) > 0 {
			plans := make([]interface{}, 0, len(j.Plans))
			for _, l := range j.Plans {
				pricingPlan := make(map[string]interface{})
				pricingPlan["sku"] = l.Sku
				if l.SkuDescription != nil {
					pricingPlan["sku_description"] = *l.SkuDescription
				}
				pricingPlan["location"] = l.Location
				pricingPlan["pricing_model"] = l.PricingModel
				price := customRound(float64(l.Price))
				pricingPlan["price"] = price
				pricingPlan["price_unit"] = l.PriceUnit
				if l.ApplicableDiscounts != nil {
					appliDis := *l.ApplicableDiscounts
					applicableDiscounts := make([]interface{}, 1)
					applicableDiscountsItem := make(map[string]interface{})
					if appliDis.DiscountedPrice != nil {
						discountedPrice := customRound(float64(*appliDis.DiscountedPrice))
						applicableDiscountsItem["discounted_price"] = discountedPrice
					}
					if appliDis.DiscountDetails != nil {
						discountDetails := make([]interface{}, 0, len(appliDis.DiscountDetails))
						for _, p := range appliDis.DiscountDetails {
							discountDetail := make(map[string]interface{})
							discountDetail["code"] = p.Code
							discountDetail["type"] = p.Type
							discountDetail["value"] = customRound(float64(p.Value))
							if p.CouponCode != nil {
								discountDetail["coupon_code"] = *p.CouponCode
							}
							discountDetails = append(discountDetails, discountDetail)
						}
						applicableDiscountsItem["discount_details"] = discountDetails
					}
					applicableDiscounts[0] = applicableDiscountsItem
					pricingPlan["applicable_discounts"] = applicableDiscounts
				}

				if l.CorrelatedProductCode != nil {
					pricingPlan["correlated_product_code"] = *l.CorrelatedProductCode
				}
				if l.PackageQuantity != nil {
					pricingPlan["package_quantity"] = customRound(float64(*l.PackageQuantity))
				}
				if l.PackageUnit != nil {
					pricingPlan["package_unit"] = *l.PackageUnit
				}
				if l.PackageDetails != nil {
					packDets := *l.PackageDetails
					packageDetails := make([]interface{}, 1)
					packageDetailsItem := make(map[string]interface{})
					if packDets.PackageQuantity != nil {
						packQuant := *packDets.PackageQuantity
						packageQuantity := make([]interface{}, 1)
						packageQuantityItem := make(map[string]interface{})
						packageQuantityItem["min"] = customRound(float64(packQuant.Min))
						packageQuantityItem["max"] = customRound(float64(packQuant.Max))
						packageQuantity[0] = packageQuantityItem
						packageDetailsItem["package_quantity"] = packageQuantity
					}
					if packDets.PackageUnit != nil {
						packageDetailsItem["package_unit"] = *packDets.PackageUnit
					}
					packageDetails[0] = packageDetailsItem
					pricingPlan["package_details"] = packageDetails
				}
				plans = append(plans, pricingPlan)
			}
			product["plans"] = plans
		}
		metadata := j.Metadata
		md := make([]interface{}, 1)
		mdItem := make(map[string]interface{})
		mdItem["ram_in_gb"] = customRound(float64(metadata.RamInGb))
		mdItem["cpu"] = metadata.Cpu
		mdItem["cpu_count"] = customRound(float64(metadata.CpuCount))
		mdItem["cores_per_cpu"] = customRound(float64(metadata.CoresPerCpu))
		mdItem["cpu_frequency"] = customRound(float64(metadata.CpuFrequency))
		mdItem["network"] = metadata.Network
		mdItem["storage"] = metadata.Storage
		if metadata.GpuConfigurations != nil {
			gpuConfs := make([]interface{}, 0, len(metadata.GpuConfigurations))
			for _, j := range metadata.GpuConfigurations {
				gpuConfMetadata := make(map[string]interface{})
				if j.Count != nil {
					gpuConfMetadata["count"] = customRound(float64(*j.Count))
				}
				if j.Name != nil {
					gpuConfMetadata["name"] = *j.Name
				}
				gpuConfs = append(gpuConfs, gpuConfMetadata)
			}
			mdItem["gpu_configurations"] = gpuConfs
		}
		md[0] = mdItem
		product["metadata"] = md
		products = append(products, product)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("products", products)
	return nil
}

// customRound rounds the number to different decimal places depending on its size.
func customRound(i float64) float64 {
	if i >= 1000 {
		return math.Round(i*100) / 100
	} else if i >= 100 {
		return math.Round(i*1000) / 1000
	} else if i >= 10 {
		return math.Round(i*10000) / 10000
	} else {
		return math.Round(i*100000) / 100000
	}
}
