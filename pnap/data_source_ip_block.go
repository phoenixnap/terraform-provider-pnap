package pnap

import (
	"fmt"

	"github.com/PNAP/go-sdk-helper-bmc/command/ipapi/ipblock"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIpBlock() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIpBlockRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr_block_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
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

func dataSourceIpBlockRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	requestCommand := ipblock.NewGetIpBlocksCommand(client)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	numOfKeys := 0
	for _, instance := range resp {
		if instance.Cidr == d.Get("cidr").(string) {
			numOfKeys++
			d.SetId(instance.Id)
			d.Set("location", instance.Location)
			d.Set("cidr_block_size", instance.CidrBlockSize)
			d.Set("cidr", instance.Cidr)
			d.Set("status", instance.Status)
			if instance.AssignedResourceId != nil {
				d.Set("assigned_resource_id", *instance.AssignedResourceId)
			}
			if instance.AssignedResourceType != nil {
				d.Set("assigned_resource_type", *instance.AssignedResourceType)
			}
		}
	}
	if numOfKeys > 1 {
		return fmt.Errorf("too many IP Blocks with CIDR %s (found %d, expected 1)", d.Get("cidr").(string), numOfKeys)
	}

	return nil
}
