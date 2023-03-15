package pnap

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PNAP/go-sdk-helper-bmc/command/billingapi/product"
	"github.com/PNAP/go-sdk-helper-bmc/dto"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProductAvailability() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProductAvailabilityRead,

		Schema: map[string]*schema.Schema{
			"product_category": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"product_code": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"show_only_min_quantity_available": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"location": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"solution": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"min_quantity": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  1,
			},
			"product_availabilities": {
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
						"location_availability_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"location": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"min_quantity_requested": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"min_quantity_available": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"available_quantity": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"solutions": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceProductAvailabilityRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	query := dto.ProductAvailabilityQuery{}
	proCatTemp := d.Get("product_category").(*schema.Set).List()
	proCat := make([]string, len(proCatTemp))
	for i, v := range proCatTemp {
		proCat[i] = fmt.Sprint(v)
	}
	query.ProductCategory = proCat
	proCodTemp := d.Get("product_code").(*schema.Set).List()
	proCod := make([]string, len(proCodTemp))
	for i, v := range proCodTemp {
		proCod[i] = fmt.Sprint(v)
	}
	query.ProductCode = proCod
	somqa := d.Get("show_only_min_quantity_available").(bool)
	query.ShowOnlyMinQuantityAvailable = somqa
	locationTemp := d.Get("location").(*schema.Set).List()
	location := make([]string, len(locationTemp))
	for i, v := range locationTemp {
		location[i] = fmt.Sprint(v)
	}
	query.Location = location
	solutionTemp := d.Get("solution").(*schema.Set).List()
	solution := make([]string, len(solutionTemp))
	for i, v := range solutionTemp {
		solution[i] = fmt.Sprint(v)
	}
	query.Solution = solution
	minQua := d.Get("min_quantity").(float64)
	if minQua > 0 {
		query.MinQuantity = float32(minQua)
	}
	queryA := &query
	b, _ := json.MarshalIndent(queryA, "", "  ")
	log.Printf("request object is" + string(b))

	requestCommand := product.NewGetProductAvailabilityCommand(client, query)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	var productAvailabilities []interface{}
	for _, j := range resp {
		productA := make(map[string]interface{})
		productA["product_code"] = j.ProductCode
		productA["product_category"] = j.ProductCategory
		if len(j.LocationAvailabilityDetails) > 0 {
			var lad []interface{}
			for _, l := range j.LocationAvailabilityDetails {
				lad1 := make(map[string]interface{})
				lad1["location"] = l.Location
				lad1["min_quantity_requested"] = int(l.MinQuantityRequested)
				lad1["min_quantity_available"] = l.MinQuantityAvailable
				lad1["available_quantity"] = int(l.AvailableQuantity)
				var solutions []interface{}
				for _, v := range l.Solutions {
					solutions = append(solutions, v)
				}
				lad1["solutions"] = solutions
				lad = append(lad, lad1)
			}
			productA["location_availability_details"] = lad
		}
		productAvailabilities = append(productAvailabilities, productA)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("product_availabilities", productAvailabilities)
	return nil
}
