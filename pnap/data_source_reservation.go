package pnap

import (
	"fmt"
	"math"

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
				Optional: true,
				Computed: true,
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
			"reservation_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initial_invoice_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quantity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quantity": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
				Optional: true,
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
			"next_billing_date": {
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
	if len(d.Get("id").(string)) > 0 && len(d.Get("sku").(string)) > 0 {
		numOfKeys := 0
		for _, instance := range resp {
			if instance.Id == d.Get("id").(string) && instance.Sku == d.Get("sku").(string) {
				numOfKeys++
				d.SetId(instance.Id)
				d.Set("id", instance.Id)
				d.Set("product_code", instance.ProductCode)
				d.Set("product_category", instance.ProductCategory)
				d.Set("location", instance.Location)
				d.Set("reservation_model", instance.ReservationModel)
				d.Set("reservation_state", instance.ReservationState)
				if instance.InitialInvoiceModel != nil {
					d.Set("initial_invoice_model", *instance.InitialInvoiceModel)
				}
				quant := flattenQuantity(instance.Quantity)
				d.Set("quantity", quant)
				d.Set("start_date_time", instance.StartDateTime.String())
				if instance.EndDateTime != nil {
					endDateTime := *instance.EndDateTime
					d.Set("end_date_time", endDateTime.String())
				}
				if instance.LastRenewalDateTime != nil {
					lastRenewalDateTime := *instance.LastRenewalDateTime
					d.Set("last_renewal_date_time", lastRenewalDateTime.String())
				}
				if instance.NextRenewalDateTime != nil {
					nextRenewalDateTime := *instance.NextRenewalDateTime
					d.Set("next_renewal_date_time", nextRenewalDateTime.String())
				}
				d.Set("auto_renew", instance.AutoRenew)
				d.Set("sku", instance.Sku)
				price := math.Round(float64(instance.Price)*100000) / 100000
				d.Set("price", price)
				d.Set("price_unit", instance.PriceUnit)
				if instance.AssignedResourceId != nil {
					d.Set("assigned_resource_id", *instance.AssignedResourceId)
				}
				if instance.NextBillingDate != nil {
					d.Set("next_billing_date", *instance.NextBillingDate)
				}
			}
		}
		if numOfKeys > 1 {
			return fmt.Errorf("too many reservations with id %s and sku %s (found %d, expected 1)", d.Get("id").(string), d.Get("sku").(string), numOfKeys)
		}
	} else if len(d.Get("sku").(string)) > 0 {
		numOfKeys := 0
		for _, instance := range resp {
			if instance.Sku == d.Get("sku").(string) {
				numOfKeys++
				d.SetId(instance.Id)
				d.Set("id", instance.Id)
				d.Set("product_code", instance.ProductCode)
				d.Set("product_category", instance.ProductCategory)
				d.Set("location", instance.Location)
				d.Set("reservation_model", instance.ReservationModel)
				d.Set("reservation_state", instance.ReservationState)
				if instance.InitialInvoiceModel != nil {
					d.Set("initial_invoice_model", *instance.InitialInvoiceModel)
				}
				quant := flattenQuantity(instance.Quantity)
				d.Set("quantity", quant)
				d.Set("start_date_time", instance.StartDateTime.String())
				if instance.EndDateTime != nil {
					endDateTime := *instance.EndDateTime
					d.Set("end_date_time", endDateTime.String())
				}
				if instance.LastRenewalDateTime != nil {
					lastRenewalDateTime := *instance.LastRenewalDateTime
					d.Set("last_renewal_date_time", lastRenewalDateTime.String())
				}
				if instance.NextRenewalDateTime != nil {
					nextRenewalDateTime := *instance.NextRenewalDateTime
					d.Set("next_renewal_date_time", nextRenewalDateTime.String())
				}
				d.Set("auto_renew", instance.AutoRenew)
				d.Set("sku", instance.Sku)
				price := math.Round(float64(instance.Price)*100000) / 100000
				d.Set("price", price)
				d.Set("price_unit", instance.PriceUnit)
				if instance.AssignedResourceId != nil {
					d.Set("assigned_resource_id", *instance.AssignedResourceId)
				}
				if instance.NextBillingDate != nil {
					d.Set("next_billing_date", *instance.NextBillingDate)
				}
			}
		}
		if numOfKeys > 1 {
			return fmt.Errorf("too many reservations with sku %s (found %d, expected 1)", d.Get("sku").(string), numOfKeys)
		}
	} else {
		numOfKeys := 0
		for _, instance := range resp {
			if instance.Id == d.Get("id").(string) {
				numOfKeys++
				d.SetId(instance.Id)
				d.Set("id", instance.Id)
				d.Set("product_code", instance.ProductCode)
				d.Set("product_category", instance.ProductCategory)
				d.Set("location", instance.Location)
				d.Set("reservation_model", instance.ReservationModel)
				d.Set("reservation_state", instance.ReservationState)
				if instance.InitialInvoiceModel != nil {
					d.Set("initial_invoice_model", *instance.InitialInvoiceModel)
				}
				quant := flattenQuantity(instance.Quantity)
				d.Set("quantity", quant)
				d.Set("start_date_time", instance.StartDateTime.String())
				if instance.EndDateTime != nil {
					endDateTime := *instance.EndDateTime
					d.Set("end_date_time", endDateTime.String())
				}
				if instance.LastRenewalDateTime != nil {
					lastRenewalDateTime := *instance.LastRenewalDateTime
					d.Set("last_renewal_date_time", lastRenewalDateTime.String())
				}
				if instance.NextRenewalDateTime != nil {
					nextRenewalDateTime := *instance.NextRenewalDateTime
					d.Set("next_renewal_date_time", nextRenewalDateTime.String())
				}
				d.Set("auto_renew", instance.AutoRenew)
				d.Set("sku", instance.Sku)
				price := math.Round(float64(instance.Price)*100000) / 100000
				d.Set("price", price)
				d.Set("price_unit", instance.PriceUnit)
				if instance.AssignedResourceId != nil {
					d.Set("assigned_resource_id", *instance.AssignedResourceId)
				}
				if instance.NextBillingDate != nil {
					d.Set("next_billing_date", *instance.NextBillingDate)
				}
			}
		}
		if numOfKeys > 1 {
			return fmt.Errorf("too many reservations with id %s (found %d, expected 1)", d.Get("id").(string), numOfKeys)
		}
	}
	return nil
}
