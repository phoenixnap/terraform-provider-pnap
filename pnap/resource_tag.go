package pnap

import (
	"fmt"

	"github.com/PNAP/go-sdk-helper-bmc/command/tagapi/tag"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tagapiclient "github.com/phoenixnap/go-sdk-bmc/tagapi/v2"
)

func resourceTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceTagCreate,
		Read:   resourceTagRead,
		Update: resourceTagUpdate,
		Delete: resourceTagDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(pnapRetryTimeout),
			Update: schema.DefaultTimeout(pnapRetryTimeout),
			Delete: schema.DefaultTimeout(pnapDeleteRetryTimeout),
		},

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
				Optional: true,
			},
			"is_billing_tag": {
				Type:     schema.TypeBool,
				Required: true,
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

func resourceTagCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(receiver.BMCSDK)

	request := &tagapiclient.TagCreate{}
	request.Name = d.Get("name").(string)
	desc := d.Get("description").(string)
	if len(desc) > 0 {
		request.Description = &desc
	}
	request.IsBillingTag = d.Get("is_billing_tag").(bool)

	requestCommand := tag.NewCreateTagCommand(client, *request)

	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	d.SetId(resp.Id)

	return resourceTagRead(d, m)
}

func resourceTagRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	tagID := d.Id()
	requestCommand := tag.NewGetTagCommand(client, tagID)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	d.SetId(resp.Id)
	d.Set("name", resp.Name)
	if resp.Values != nil {
		readValues := resp.Values
		var values []interface{}
		for _, v := range readValues {
			values = append(values, v)
		}
		d.Set("values", values)
	}
	if resp.Description != nil {
		d.Set("description", *resp.Description)
	}
	d.Set("is_billing_tag", resp.IsBillingTag)
	if resp.ResourceAssignments != nil {
		resAssigns := resp.ResourceAssignments
		assigns := make([]interface{}, len(resAssigns))
		for i, v := range resAssigns {
			a := make(map[string]interface{})
			a["resource_name"] = v.ResourceName
			if v.Value != nil {
				a["value"] = *v.Value
			}
			assigns[i] = a
		}
		d.Set("resource_assignments", assigns)
	}
	if resp.CreatedBy != nil {
		d.Set("created_by", *resp.CreatedBy)
	}
	return nil
}

func resourceTagUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("name") || d.HasChange("is_billing_tag") || d.HasChange("description") {
		client := m.(receiver.BMCSDK)
		tagID := d.Id()

		request := &tagapiclient.TagUpdate{}
		request.Name = d.Get("name").(string)
		desc := d.Get("description").(string)
		request.Description = &desc
		request.IsBillingTag = d.Get("is_billing_tag").(bool)

		requestCommand := tag.NewUpdateTagCommand(client, tagID, *request)
		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported action")
	}
	return resourceTagRead(d, m)

}

func resourceTagDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	tagID := d.Id()

	requestCommand := tag.NewDeleteTagCommand(client, tagID)
	_, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	return nil
}
