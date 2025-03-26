package pnap

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	networkapiclient "github.com/phoenixnap/go-sdk-bmc/networkapi/v4"

	"github.com/PNAP/go-sdk-helper-bmc/command/networkapi/publicnetwork"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
)

func dataSourcePublicNetwork() *schema.Resource {
	return &schema.Resource{

		Read: dataSourcePublicNetworkRead,

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
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used_ips_count": {
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ra_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
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
	name := d.Get("name").(string)
	id := d.Get("id").(string)
	for _, instance := range resp {
		if instance.Name == name || instance.Id == id {
			numOfNets++
			d.SetId(instance.Id)
			d.Set("location", instance.Location)
			d.Set("name", instance.Name)
			if instance.Description != nil {
				d.Set("description", *instance.Description)
			} else {
				d.Set("description", "")
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
			d.Set("status", instance.Status)
			if instance.RaEnabled != nil {
				d.Set("ra_enabled", *instance.RaEnabled)
			} else {
				d.Set("ra_enabled", nil)
			}
		}
	}
	if numOfNets > 1 && len(name) > 0 {
		return fmt.Errorf("too many public networks with name %s (found %d, expected 1)", name, numOfNets)
	} else if numOfNets > 1 && len(id) > 0 {
		return fmt.Errorf("too many public networks with ID %s (found %d, expected 1)", id, numOfNets)
	}

	return nil
}

func flattenDataIpBlocks(ipBlocks []networkapiclient.PublicNetworkIpBlock) []interface{} {
	if len(ipBlocks) > 0 {
		ib := make([]interface{}, len(ipBlocks))
		for i, j := range ipBlocks {
			ibItem := make(map[string]interface{})
			ibItem["id"] = j.Id
			ibItem["cidr"] = j.Cidr
			ibItem["used_ips_count"] = j.UsedIpsCount
			ib[i] = ibItem
		}
		return ib
	}
	return make([]interface{}, 0)
}
