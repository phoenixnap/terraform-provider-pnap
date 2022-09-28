package pnap

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/PNAP/go-sdk-helper-bmc/command/networkapi/publicnetwork"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"

	networkapiclient "github.com/phoenixnap/go-sdk-bmc/networkapi/v2"
)

func resourcePublicNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicNetworkCreate,
		Read:   resourcePublicNetworkRead,
		Update: resourcePublicNetworkUpdate,
		Delete: resourcePublicNetworkDelete,

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
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_blocks": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_network_ip_block": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
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

func resourcePublicNetworkCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(receiver.BMCSDK)

	request := &networkapiclient.PublicNetworkCreate{}
	request.Name = d.Get("name").(string)
	request.Location = d.Get("location").(string)
	var desc = d.Get("description").(string)
	if len(desc) > 0 {
		request.Description = &desc
	}
	ipBlocks := d.Get("ip_blocks").([]interface{})
	if len(ipBlocks) > 0 {
		ipBlocksObject := make([]networkapiclient.PublicNetworkIpBlock, len(ipBlocks))
		for i, j := range ipBlocks {
			ibItem := j.(map[string]interface{})
			pnib := ibItem["public_network_ip_block"].([]interface{})[0]
			pnibItem := pnib.(map[string]interface{})

			pnibObject := networkapiclient.PublicNetworkIpBlock{}
			pnibObject.Id = pnibItem["id"].(string)
			ipBlocksObject[i] = pnibObject
		}
		request.IpBlocks = ipBlocksObject
	}
	requestCommand := publicnetwork.NewCreatePublicNetworkCommand(client, *request)

	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	d.SetId(resp.Id)

	return resourcePublicNetworkRead(d, m)
}

func resourcePublicNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	networkID := d.Id()
	requestCommand := publicnetwork.NewGetPublicNetworkCommand(client, networkID)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	d.SetId(resp.Id)
	d.Set("name", resp.Name)
	d.Set("location", resp.Location)
	desc := resp.Description
	if desc != nil {
		d.Set("description", *resp.Description)
	} else {
		d.Set("description", "")
	}

	if len(resp.IpBlocks) > 0 {
		var ibInput = d.Get("ip_blocks").([]interface{})
		d.Set("ip_blocks", ibInput)
	}
	if len(resp.CreatedOn.String()) > 0 {
		d.Set("created_on", resp.CreatedOn.String())
	}
	d.Set("vlan_id", resp.VlanId)

	memberships := flattenMemberships(resp.Memberships)

	if err := d.Set("memberships", memberships); err != nil {
		return err
	}
	return nil
}

func resourcePublicNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("name") || d.HasChange("description") {
		client := m.(receiver.BMCSDK)
		networkID := d.Id()
		request := &networkapiclient.PublicNetworkModify{}
		var name = d.Get("name").(string)
		request.Name = &name
		var desc = d.Get("description").(string)
		request.Description = &desc
		requestCommand := publicnetwork.NewUpdatePublicNetworkCommand(client, networkID, *request)
		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported action")
	}
	return resourcePublicNetworkRead(d, m)
}

func resourcePublicNetworkDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	networkID := d.Id()

	requestCommand := publicnetwork.NewDeletePublicNetworkCommand(client, networkID)
	err := requestCommand.Execute()
	if err != nil {
		return err
	}

	return nil
}

func flattenMemberships(memberships []networkapiclient.NetworkMembership) []interface{} {
	if memberships != nil {
		mems := make([]interface{}, len(memberships))
		for i, v := range memberships {
			mem := make(map[string]interface{})
			mem["resource_id"] = v.ResourceId
			mem["resource_type"] = v.ResourceType
			ips := make([]interface{}, len(v.Ips))
			for j, k := range v.Ips {
				ips[j] = k
			}
			mem["ips"] = ips
			mems[i] = mem
		}
		return mems
	}
	return make([]interface{}, 0)
}
