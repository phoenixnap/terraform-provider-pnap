package pnap

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/PNAP/go-sdk-helper-bmc/command/paymentsapi/transaction"
	"github.com/PNAP/go-sdk-helper-bmc/dto"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTransactions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTransactionsRead,

		Schema: map[string]*schema.Schema{
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sort_direction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"from": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"to": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"paginated_transactions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"offset": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"results": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"details": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"amount": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"currency": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"metadata": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"invoice_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"invoice_number": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"is_auto_charge": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"card_payment_method_details": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"card_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"last_four_digits": {
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

func dataSourceTransactionsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	query := dto.Query{}
	query.Limit = int32(d.Get("limit").(int))
	query.Offset = int32(d.Get("offset").(int))
	query.SortDirection = d.Get("sort_direction").(string)
	query.SortField = d.Get("sort_field").(string)

	from := d.Get("from").(string)
	if from != "" {
		t1, err1 := time.Parse(time.RFC3339, from)
		if err1 != nil {
			return err1
		} else {
			query.From = t1
		}
	}
	to := d.Get("to").(string)
	if to != "" {
		t2, err2 := time.Parse(time.RFC3339, to)
		if err2 != nil {
			return err2
		} else {
			query.To = t2
		}
	}

	requestCommand := transaction.NewGetTransactionsCommand(client, query)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	paginatedTransactions := make([]interface{}, 1)
	paginatedResponse := make(map[string]interface{})
	paginatedResponse["limit"] = int(resp.Limit)
	paginatedResponse["offset"] = int(resp.Offset)
	paginatedResponse["total"] = int(resp.Total)

	id := d.Get("id").(string)

	if len(id) > 0 {
		numOfTransactions := 0
		for _, j := range resp.Results {
			if j.Id == id {
				numOfTransactions++

				transactionMap := make(map[string]interface{})
				transactionMap["id"] = j.Id
				transactionMap["status"] = j.Status
				if j.Details != nil {
					transactionMap["details"] = *j.Details
				}
				amount := math.Round(float64(j.Amount)*100) / 100
				transactionMap["amount"] = amount
				transactionMap["currency"] = j.Currency
				transactionMap["date"] = j.Date.String()

				metadata := make([]interface{}, 1)
				transactionMetadata := make(map[string]interface{})
				transactionMetadata["invoice_id"] = j.Metadata.InvoiceId
				if j.Metadata.InvoiceNumber != nil {
					transactionMetadata["invoice_number"] = *j.Metadata.InvoiceNumber
				}
				transactionMetadata["is_auto_charge"] = j.Metadata.IsAutoCharge
				metadata[0] = transactionMetadata
				transactionMap["metadata"] = metadata

				cardPaymentMethodDetails := make([]interface{}, 1)
				cardPaymentMethodDetailsItem := make(map[string]interface{})
				cardPaymentMethodDetailsItem["card_type"] = j.CardPaymentMethodDetails.CardType
				cardPaymentMethodDetailsItem["last_four_digits"] = j.CardPaymentMethodDetails.LastFourDigits
				cardPaymentMethodDetails[0] = cardPaymentMethodDetailsItem
				transactionMap["card_payment_method_details"] = cardPaymentMethodDetails

				result := make([]interface{}, 1)
				result[0] = transactionMap
				paginatedResponse["results"] = result
				paginatedResponse["total"] = numOfTransactions
				paginatedTransactions[0] = paginatedResponse

				d.SetId(j.Id)
				d.Set("paginated_transactions", paginatedTransactions)
			}
		}
		if numOfTransactions > 1 {
			return fmt.Errorf("too many transactions with id %s (found %d, expected 1)", id, numOfTransactions)
		}

	} else {

		results := make([]interface{}, len(resp.Results))

		for i, j := range resp.Results {

			transactionMap := make(map[string]interface{})
			transactionMap["id"] = j.Id
			transactionMap["status"] = j.Status
			if j.Details != nil {
				transactionMap["details"] = *j.Details
			}
			amount := math.Round(float64(j.Amount)*100) / 100
			transactionMap["amount"] = amount
			transactionMap["currency"] = j.Currency
			transactionMap["date"] = j.Date.String()

			metadata := make([]interface{}, 1)
			transactionMetadata := make(map[string]interface{})
			transactionMetadata["invoice_id"] = j.Metadata.InvoiceId
			if j.Metadata.InvoiceNumber != nil {
				transactionMetadata["invoice_number"] = *j.Metadata.InvoiceNumber
			}
			transactionMetadata["is_auto_charge"] = j.Metadata.IsAutoCharge
			metadata[0] = transactionMetadata
			transactionMap["metadata"] = metadata

			cardPaymentMethodDetails := make([]interface{}, 1)
			cardPaymentMethodDetailsItem := make(map[string]interface{})
			cardPaymentMethodDetailsItem["card_type"] = j.CardPaymentMethodDetails.CardType
			cardPaymentMethodDetailsItem["last_four_digits"] = j.CardPaymentMethodDetails.LastFourDigits
			cardPaymentMethodDetails[0] = cardPaymentMethodDetailsItem
			transactionMap["card_payment_method_details"] = cardPaymentMethodDetails

			results[i] = transactionMap
		}

		paginatedResponse["results"] = results
		paginatedTransactions[0] = paginatedResponse

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		d.Set("paginated_transactions", paginatedTransactions)
	}
	return nil
}
