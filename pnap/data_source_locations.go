package pnap

import (
	"strconv"
	"time"

	"github.com/PNAP/go-sdk-helper-bmc/command/locationapi/location"
	"github.com/PNAP/go-sdk-helper-bmc/dto"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	locationapiclient "github.com/phoenixnap/go-sdk-bmc/locationapi"
)

func dataSourceLocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationsRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"product_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_categories": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"product_category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"product_category_description": {
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

func dataSourceLocationsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	query := dto.Query{}

	loc := d.Get("location").(string)
	if len(loc) > 0 {
		locEnum, errorLoc := locationapiclient.NewLocationEnumFromValue(loc)
		if errorLoc != nil {
			return errorLoc
		}
		query.Location = *locEnum
	}
	productCategory := d.Get("product_category").(string)
	if len(productCategory) > 0 {
		prodCatEnum, errorProd := locationapiclient.NewProductCategoryEnumFromValue(productCategory)
		if errorProd != nil {
			return errorProd
		}
		query.ProductCategory = *prodCatEnum
	}

	requestCommand := location.NewGetLocationsCommand(client, query)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	var locations []interface{}
	for _, j := range resp {
		locationMap := make(map[string]interface{})
		locationEnum := j.Location
		locationMap["location"] = string(locationEnum)
		if j.LocationDescription != nil {
			locationMap["location_description"] = *j.LocationDescription
		}
		if len(j.ProductCategories) > 0 {
			var prodCategories []interface{}
			for _, l := range j.ProductCategories {
				prodCategory := make(map[string]interface{})
				prodCat := l.ProductCategory
				prodCategory["product_category"] = string(prodCat)
				prodCatDesc := l.ProductCategoryDescription
				if prodCatDesc != nil {
					prodCategory["product_category_description"] = *prodCatDesc
				}
				prodCategories = append(prodCategories, prodCategory)
			}
			locationMap["product_categories"] = prodCategories
		}
		locations = append(locations, locationMap)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("locations", locations)
	return nil
}
