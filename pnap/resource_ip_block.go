package pnap

import (
	"fmt"

	"github.com/PNAP/go-sdk-helper-bmc/command/ipapi/ipblock"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ipapiclient "github.com/phoenixnap/go-sdk-bmc/ipapi"
)

func resourceIpBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpBlockCreate,
		Read:   resourceIpBlockRead,
		Update: resourceIpBlockUpdate,
		Delete: resourceIpBlockDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(pnapRetryTimeout),
			Update: schema.DefaultTimeout(pnapRetryTimeout),
			Delete: schema.DefaultTimeout(pnapDeleteRetryTimeout),
		},

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr_block_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assigned_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assigned_resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIpBlockCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(receiver.BMCSDK)

	request := &ipapiclient.IpBlockCreate{}
	request.Location = d.Get("location").(string)
	request.CidrBlockSize = d.Get("cidr_block_size").(string)

	requestCommand := ipblock.NewCreateIpBlockCommand(client, *request)

	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	d.SetId(resp.Id)

	return resourceIpBlockRead(d, m)
}

func resourceIpBlockRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	ipBlockID := d.Id()
	requestCommand := ipblock.NewGetIpBlockCommand(client, ipBlockID)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	d.SetId(resp.Id)
	d.Set("location", resp.Location)
	d.Set("cidr_block_size", resp.CidrBlockSize)
	d.Set("cidr", resp.Cidr)
	d.Set("status", resp.Status)
	if resp.AssignedResourceId != nil {
		d.Set("assigned_resource_id", *resp.AssignedResourceId)
	}
	if resp.AssignedResourceType != nil {
		d.Set("assigned_resource_type", *resp.AssignedResourceType)
	}

	return nil
}

func resourceIpBlockUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("location") || d.HasChange("cidr_block_size") {
		return fmt.Errorf("unsupported action")
	}

	return resourceIpBlockRead(d, m)
}

func resourceIpBlockDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	ipBlockID := d.Id()

	requestCommand := ipblock.NewDeleteIpBlockCommand(client, ipBlockID)
	_, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	return nil
}
