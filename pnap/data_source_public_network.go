package pnap

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	networkapiclient "github.com/phoenixnap/go-sdk-bmc/networkapi"

	"github.com/PNAP/go-sdk-helper-bmc/command/networkapi/publicnetwork"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
)

func dataSourcePublicNetwork() *schema.Resource {
	return &schema.Resource{

		Read: dataSourcePublicNetworkRead,

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_blocks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memberships": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ips": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourcePublicNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	requestCommand := publicnetwork.NewGetPublicNetworksCommand(client)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	numOfNets := 0
	for _, instance := range resp {
		if instance.Name == d.Get("name").(string) {
			numOfNets++
			d.SetId(instance.Id)
			d.Set("location", instance.Location)
			d.Set("name", instance.Name)
			desc := instance.Description
			if desc != nil {
				d.Set("description", *instance.Description)
			}
			ipBlocks := flattenDataIpBlocks(instance.IpBlocks)
			if err := d.Set("ip_blocks", ipBlocks); err != nil {
				return err
			}
			d.Set("created_on", instance.CreatedOn.String())
			d.Set("vlan_id", instance.VlanId)

			memberships := flattenMemberships(instance.Memberships)

			if err := d.Set("memberships", memberships); err != nil {
				return err
			}
		}
	}
	if numOfNets > 1 {
		return fmt.Errorf("too many public networks with name %s (found %d, expected 1)", d.Get("name").(string), numOfNets)
	}
	return nil
}

func flattenDataIpBlocks(ipBlocks []networkapiclient.PublicNetworkIpBlock) []interface{} {
	if len(ipBlocks) > 0 {
		ib := make([]interface{}, len(ipBlocks))
		for i, j := range ipBlocks {
			ibItem := make(map[string]interface{})
			ibItem["id"] = j.Id
			ib[i] = ibItem
		}
		return ib
	}
	return make([]interface{}, 0)
}
