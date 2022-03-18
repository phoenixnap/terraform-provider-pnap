package pnap

import (
	"fmt"

	"github.com/PNAP/go-sdk-helper-bmc/command/tagapi/tag"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTagRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_billing_tag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"resource_assignments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTagRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	requestCommand := tag.NewGetTagsCommand(client)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	numOfTags := 0
	for _, instance := range resp {
		if instance.Name == d.Get("name").(string) {
			numOfTags++
			d.SetId(instance.Id)
			d.Set("name", instance.Name)
			if instance.Values != nil {
				readValues := *instance.Values
				var values []interface{}
				for _, v := range readValues {
					values = append(values, v)
				}
				d.Set("values", values)
			}
			if instance.Description != nil {
				d.Set("description", *instance.Description)
			}
			d.Set("is_billing_tag", instance.IsBillingTag)
			if instance.ResourceAssignments != nil {
				readAssigns := *instance.ResourceAssignments
				assigns := make([]interface{}, len(readAssigns))
				for i, a := range readAssigns {
					assign := make(map[string]interface{})
					assign["resource_name"] = a.ResourceName
					if a.Value != nil {
						assign["value"] = *a.Value
					}
					assigns[i] = assign
				}
				d.Set("resource_assignments", assigns)
			}
			if instance.CreatedBy != nil {
				d.Set("created_by", *instance.CreatedBy)
			}
		}
	}
	if numOfTags > 1 {
		return fmt.Errorf("too many tags with name %s (found %d, expected 1)", d.Get("name").(string), numOfTags)
	}
	return nil
}
