package pnap

import (
	"fmt"
	"math"

	"github.com/PNAP/go-sdk-helper-bmc/command/billingapi/reservation"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	billingapiclient "github.com/phoenixnap/go-sdk-bmc/billingapi/v3"
)

func resourceReservation() *schema.Resource {
	return &schema.Resource{
		Create: resourceReservationCreate,
		Read:   resourceReservationRead,
		Update: resourceReservationUpdate,
		Delete: resourceReservationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(pnapRetryTimeout),
			Update: schema.DefaultTimeout(pnapRetryTimeout),
			Delete: schema.DefaultTimeout(pnapDeleteRetryTimeout),
		},

		Schema: map[string]*schema.Schema{
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
				Optional: true,
				Computed: true,
			},
			"sku": {
				Type:     schema.TypeString,
				Required: true,
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
			"auto_renew_disable_reason": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceReservationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	request := &billingapiclient.ReservationRequest{}
	request.Sku = d.Get("sku").(string)
	requestCommand := reservation.NewCreateReservationCommand(client, *request)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	d.SetId(resp.Id)
	return resourceReservationRead(d, m)
}

func resourceReservationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	reservationID := d.Id()
	requestCommand := reservation.NewGetReservationCommand(client, reservationID)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	d.SetId(resp.Id)
	d.Set("product_code", resp.ProductCode)
	d.Set("product_category", resp.ProductCategory)
	d.Set("location", resp.Location)
	d.Set("reservation_model", resp.ReservationModel)
	if resp.InitialInvoiceModel != nil {
		d.Set("initial_invoice_model", *resp.InitialInvoiceModel)
	}
	d.Set("start_date_time", resp.StartDateTime.String())
	if resp.EndDateTime != nil {
		endDateTime := *resp.EndDateTime
		d.Set("end_date_time", endDateTime.String())
	}
	if resp.LastRenewalDateTime != nil {
		lastRenewalDateTime := *resp.LastRenewalDateTime
		d.Set("last_renewal_date_time", lastRenewalDateTime.String())
	}
	if resp.NextRenewalDateTime != nil {
		nextRenewalDateTime := *resp.NextRenewalDateTime
		d.Set("next_renewal_date_time", nextRenewalDateTime.String())
	}
	d.Set("auto_renew", resp.AutoRenew)
	d.Set("sku", resp.Sku)
	price := math.Round(float64(resp.Price)*100000) / 100000
	d.Set("price", price)
	d.Set("price_unit", resp.PriceUnit)
	if resp.AssignedResourceId != nil {
		d.Set("assigned_resource_id", *resp.AssignedResourceId)
	}
	if resp.NextBillingDate != nil {
		d.Set("next_billing_date", *resp.NextBillingDate)
	}
	// d.Set("auto_renew_disable_reason", "")

	return nil
}

func resourceReservationUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("sku") {
		client := m.(receiver.BMCSDK)
		reservationID := d.Id()
		request := &billingapiclient.ReservationRequest{}
		request.Sku = d.Get("sku").(string)
		requestCommand := reservation.NewConvertReservationCommand(client, reservationID, *request)
		resp, err := requestCommand.Execute()
		if err != nil {
			return err
		}
		d.SetId(resp.Id)
	} else if d.HasChange("auto_renew") {
		client := m.(receiver.BMCSDK)
		newStatus := d.Get("auto_renew").(bool)
		if !newStatus {
			reservationID := d.Id()
			request := &billingapiclient.ReservationAutoRenewDisableRequest{}
			var reason = d.Get("auto_renew_disable_reason").(string)
			if len(reason) > 0 {
				request.AutoRenewDisableReason = &reason
			}
			requestCommand := reservation.NewDisableAutoRenewReservationCommand(client, reservationID, *request)
			_, err := requestCommand.Execute()
			if err != nil {
				return err
			}
		} else if newStatus {
			reservationID := d.Id()
			requestCommand := reservation.NewEnableAutoRenewReservationCommand(client, reservationID)
			_, err := requestCommand.Execute()
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported action")
		}
	} else {
		return fmt.Errorf("unsupported action")
	}
	return resourceReservationRead(d, m)
}

func resourceReservationDelete(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("unsupported action")
}
