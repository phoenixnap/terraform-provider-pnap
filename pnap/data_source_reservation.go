package pnap

import (
	"fmt"

	"github.com/PNAP/go-sdk-helper-bmc/command/billingapi/reservation"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceReservation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceReservationRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reservation_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initial_invoice_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_date_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_date_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_renewal_date_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_renewal_date_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sku": {
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
			"assigned_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceReservationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	requestCommand := reservation.NewGetReservationsCommand(client)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	numOfKeys := 0
	for _, instance := range resp {
		if instance.ID == d.Get("id").(string) {
			numOfKeys++
			d.SetId(instance.ID)
			d.Set("id", instance.ID)
			d.Set("product_code", instance.ProductCode)
			d.Set("product_category", instance.ProductCategory)
			d.Set("location", instance.Location)
			d.Set("reservation_model", instance.ReservationModel)
			d.Set("initial_invoice_model", instance.InitialInvoiceModel)
			d.Set("start_date_time", instance.StartDateTime.String())
			d.Set("end_date_time", instance.EndDateTime.String())
			d.Set("last_renewal_date_time", instance.LastRenewalDateTime.String())
			d.Set("next_renewal_date_time", instance.NextRenewalDateTime.String())
			d.Set("auto_renew", instance.AutoRenew)
			d.Set("sku", instance.SKU)
			d.Set("price", instance.Price)
			d.Set("price_unit", instance.PriceUnit)
			d.Set("assigned_resource_id", instance.AssignedResourceID)
		}
	}
	if numOfKeys > 1 {
		return fmt.Errorf("too many reservations with id %s (found %d, expected 1)", d.Get("id").(string), numOfKeys)
	}

	return nil
}
