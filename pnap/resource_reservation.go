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
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quantity": {
							Type:     schema.TypeFloat,
							Required: true,
						},
						"unit": {
							Type:     schema.TypeString,
							Required: true,
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
			"utilization": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"percentage": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceReservationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	request := &billingapiclient.ReservationRequest{}
	request.Sku = d.Get("sku").(string)
	if d.Get("quantity") != nil && len(d.Get("quantity").([]interface{})) > 0 {
		quantity := d.Get("quantity").([]interface{})[0]
		quantityItem := quantity.(map[string]interface{})
		quan := quantityItem["quantity"].(float64)
		unit := quantityItem["unit"].(string)
		quantityObject := billingapiclient.Quantity{}
		quantityObject.Quantity = float32(quan)

		unitEnum, errorUnit := billingapiclient.NewQuantityUnitEnumFromValue(unit)
		if errorUnit != nil {
			return errorUnit
		}
		quantityObject.Unit = *unitEnum
		request.Quantity = &quantityObject
	}
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
	d.Set("reservation_state", resp.ReservationState)
	if resp.InitialInvoiceModel != nil {
		d.Set("initial_invoice_model", *resp.InitialInvoiceModel)
	}
	quant := flattenQuantity(&resp.Quantity)
	d.Set("quantity", quant)
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
	utilization := flattenUtilization(resp.Utilization)
	d.Set("utilization", utilization)

	return nil
}

func resourceReservationUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("sku") || d.HasChange("quantity") {
		client := m.(receiver.BMCSDK)
		reservationID := d.Id()
		request := &billingapiclient.ReservationRequest{}
		request.Sku = d.Get("sku").(string)
		if d.Get("quantity") != nil && len(d.Get("quantity").([]interface{})) > 0 {
			quantity := d.Get("quantity").([]interface{})[0]
			quantityItem := quantity.(map[string]interface{})
			quan := quantityItem["quantity"].(float64)
			unit := quantityItem["unit"].(string)
			quantityObject := billingapiclient.Quantity{}
			quantityObject.Quantity = float32(quan)

			unitEnum, errorUnit := billingapiclient.NewQuantityUnitEnumFromValue(unit)
			if errorUnit != nil {
				return errorUnit
			}
			quantityObject.Unit = *unitEnum
			request.Quantity = &quantityObject
		}
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

func flattenQuantity(quantity *billingapiclient.Quantity) []interface{} {
	quant := make([]interface{}, 1)
	quantItem := make(map[string]interface{})
	size := quantity.Quantity
	if size > 0 {
		quantItem["quantity"] = float64(size)
	}
	unit := quantity.Unit
	if len(unit) > 0 {
		quantItem["unit"] = string(unit)
	}
	quant[0] = quantItem
	return quant
}

func flattenUtilization(utilization *billingapiclient.Utilization) []interface{} {
	util := make([]interface{}, 1)
	utilItem := make(map[string]interface{})
	if utilization != nil {
		quantity := utilization.Quantity
		quant := flattenQuantity(&quantity)
		utilItem["quantity"] = quant
		percentage := utilization.Percentage
		if percentage >= 0 {
			utilItem["percentage"] = float64(percentage)
		}
	}
	util[0] = utilItem
	return util
}
