package pnap

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/PNAP/go-sdk-helper-bmc/command/networkapi/privatenetwork"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
)

func dataSourcePrivateNetwork() *schema.Resource {
	return &schema.Resource{

		Read: dataSourcePrivateNetworkRead,

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
			"location_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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

func dataSourcePrivateNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	requestCommand := privatenetwork.NewGetPrivateNetworksCommand(client)
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
			d.Set("cidr", instance.Cidr)
			d.Set("description", instance.Description)
			d.Set("location_default", instance.LocationDefault)
			d.Set("type", instance.Type)
			d.Set("vlan_id", instance.VlanId)
			servers := flattenServers(instance.Servers)

			if err := d.Set("servers", servers); err != nil {
				return err
			}

		}
	}
	if numOfNets > 1 {
		return fmt.Errorf("too many private networks with name %s (found %d, expected 1)", d.Get("name").(string), numOfNets)
	}
	return nil
}
