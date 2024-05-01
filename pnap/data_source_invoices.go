package pnap

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/PNAP/go-sdk-helper-bmc/command/invoicingapi/invoice"
	"github.com/PNAP/go-sdk-helper-bmc/dto"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInvoices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInvoicesRead,

		Schema: map[string]*schema.Schema{
			"number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sent_on_from": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sent_on_to": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sort_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_direction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pdf_folder_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"paginated_invoices": {
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
									"number": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"currency": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"amount": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"outstanding_amount": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sent_on": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"due_date": {
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

func dataSourceInvoicesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	query := dto.Query{}
	query.Number = d.Get("number").(string)
	query.Status = d.Get("status").(string)

	sentOnFrom := d.Get("sent_on_from").(string)
	if sentOnFrom != "" {
		t1, err1 := time.Parse(time.RFC3339, sentOnFrom)
		if err1 != nil {
			return err1
		} else {
			query.SentOnFrom = t1
		}
	}
	sentOnTo := d.Get("sent_on_to").(string)
	if sentOnTo != "" {
		t2, err2 := time.Parse(time.RFC3339, sentOnTo)
		if err2 != nil {
			return err2
		} else {
			query.SentOnTo = t2
		}
	}
	query.Limit = int32(d.Get("limit").(int))
	query.Offset = int32(d.Get("offset").(int))
	query.SortField = d.Get("sort_field").(string)
	query.SortDirection = d.Get("sort_direction").(string)

	requestCommand := invoice.NewGetInvoicesCommand(client, query)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	paginatedInvoices := make([]interface{}, 1)
	paginatedResponse := make(map[string]interface{})
	paginatedResponse["limit"] = int(resp.Limit)
	paginatedResponse["offset"] = int(resp.Offset)
	paginatedResponse["total"] = int(resp.Total)

	id := d.Get("id").(string)
	path := d.Get("pdf_folder_path").(string)

	if len(id) > 0 {
		numOfInvoices := 0
		for _, j := range resp.Results {
			if j.Id == id {
				numOfInvoices++

				invoiceMap := make(map[string]interface{})
				invoiceMap["id"] = j.Id
				invoiceMap["number"] = j.Number
				invoiceMap["currency"] = j.Currency

				amount := math.Round(float64(j.Amount)*100) / 100
				invoiceMap["amount"] = amount

				outstandingAmount := math.Round(float64(j.OutstandingAmount)*100) / 100
				invoiceMap["outstanding_amount"] = outstandingAmount

				invoiceMap["status"] = j.Status
				invoiceMap["sent_on"] = j.SentOn.String()
				invoiceMap["due_date"] = j.DueDate.String()

				if len(path) > 0 {
					pdfRequestCommand := invoice.NewGenerateInvoicePdfCommand(client, id)
					pdf, err := pdfRequestCommand.Execute()
					if err != nil {
						return err
					}
					data, err := io.ReadAll(pdf)
					if err != nil {
						return err
					}
					invoicePdf, err := os.Create(path + j.Number + ".pdf")
					if err != nil {
						return err
					}
					defer invoicePdf.Close()

					if _, err := invoicePdf.Write(data); err != nil {
						return (err)
					}
				}
				result := make([]interface{}, 1)
				result[0] = invoiceMap
				paginatedResponse["results"] = result
				paginatedResponse["total"] = numOfInvoices
				paginatedInvoices[0] = paginatedResponse

				d.SetId(j.Id)
				d.Set("paginated_invoices", paginatedInvoices)
			}
		}
		if numOfInvoices > 1 {
			return fmt.Errorf("too many invoices with id %s (found %d, expected 1)", id, numOfInvoices)
		}

	} else {

		results := make([]interface{}, len(resp.Results))

		for i, j := range resp.Results {
			id = j.Id
			invoiceMap := make(map[string]interface{})
			invoiceMap["id"] = j.Id
			invoiceMap["number"] = j.Number
			invoiceMap["currency"] = j.Currency

			amount := math.Round(float64(j.Amount)*100) / 100
			invoiceMap["amount"] = amount

			outstandingAmount := math.Round(float64(j.OutstandingAmount)*100) / 100
			invoiceMap["outstanding_amount"] = outstandingAmount

			invoiceMap["status"] = j.Status
			invoiceMap["sent_on"] = j.SentOn.String()
			invoiceMap["due_date"] = j.DueDate.String()
			if len(path) > 0 {
				pdfRequestCommand := invoice.NewGenerateInvoicePdfCommand(client, id)
				pdf, err := pdfRequestCommand.Execute()
				if err != nil {
					return err
				}
				data, err := io.ReadAll(pdf)
				if err != nil {
					return err
				}
				invoicePdf, err := os.Create(path + j.Number + ".pdf")
				if err != nil {
					return err
				}
				defer invoicePdf.Close()

				if _, err := invoicePdf.Write(data); err != nil {
					return (err)
				}
			}

			results[i] = invoiceMap
		}

		paginatedResponse["results"] = results
		paginatedInvoices[0] = paginatedResponse

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		d.Set("paginated_invoices", paginatedInvoices)
	}
	return nil
}
