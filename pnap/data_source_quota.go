package pnap

import (
	"fmt"

	"github.com/PNAP/go-sdk-helper-bmc/command/bmcapi/quota"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceQuotaRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"quota_edit_limit_request_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"requested_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceQuotaRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	requestCommand := quota.NewGetQuotasCommand(client)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	numOfQuotas := 0
	for _, instance := range resp {
		if instance.Name == d.Get("name").(string) || instance.Id == d.Get("id").(string) {
			numOfQuotas++
			d.SetId(instance.Id)
			d.Set("name", instance.Name)
			d.Set("description", instance.Description)
			d.Set("status", instance.Status)
			d.Set("limit", int(instance.Limit))
			d.Set("unit", instance.Unit)
			d.Set("used", int(instance.Used))
			quotaRequests := instance.QuotaEditLimitRequestDetails
			qelrd := make([]interface{}, len(quotaRequests))
			for i, j := range quotaRequests {
				qelrdItem := make(map[string]interface{})
				qelrdItem["limit"] = int(j.Limit)
				qelrdItem["reason"] = j.Reason
				qelrdItem["requested_on"] = j.RequestedOn.String()
				qelrd[i] = qelrdItem
			}
			d.Set("quota_edit_limit_request_details", qelrd)
		}
	}
	if numOfQuotas > 1 {
		return fmt.Errorf("too many Quotas with name %s (found %d, expected 1)", d.Get("name").(string), numOfQuotas)
	}

	return nil
}
