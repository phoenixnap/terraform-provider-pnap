package pnap

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/PNAP/go-sdk-helper-bmc/command/networkapi/publicnetwork"
	"github.com/PNAP/go-sdk-helper-bmc/dto"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"

	networkapiclient "github.com/phoenixnap/go-sdk-bmc/networkapi/v4"
)

const (
	pnapPublicNetworkRetryDelay   = 10 * time.Second
	pnapPublicNetworkRetryTimeout = 7 * time.Minute
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
					},
				},
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
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
				Optional: true,
				Computed: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
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
	var vlanId = d.Get("vlan_id").(int)
	if vlanId > 0 {
		vlanId32 := int32(vlanId)
		request.VlanId = &vlanId32
	}
	ipBlocks := d.Get("ip_blocks").([]interface{})
	if len(ipBlocks) > 0 {
		ipBlocksObject := make([]networkapiclient.PublicNetworkIpBlockCreate, len(ipBlocks))
		for i, j := range ipBlocks {
			ibItem := j.(map[string]interface{})
			pnib := ibItem["public_network_ip_block"].([]interface{})[0]
			pnibItem := pnib.(map[string]interface{})

			pnibObject := networkapiclient.PublicNetworkIpBlockCreate{}
			pnibObject.Id = pnibItem["id"].(string)
			ipBlocksObject[i] = pnibObject
		}
		request.IpBlocks = ipBlocksObject
	}
	raEnabledInterface, exists := d.GetOkExists("ra_enabled")
	if exists {
		raEnabled := raEnabledInterface.(bool)
		request.RaEnabled = &raEnabled
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
	var ipBlocksInput = d.Get("ip_blocks").([]interface{})
	ipBlocks := flattenIpBlocks(resp.IpBlocks, ipBlocksInput)

	if err := d.Set("ip_blocks", ipBlocks); err != nil {
		return err
	}
	if len(resp.CreatedOn.String()) > 0 {
		d.Set("created_on", resp.CreatedOn.String())
	}
	d.Set("vlan_id", resp.VlanId)

	memberships := flattenMemberships(resp.Memberships)

	if err := d.Set("memberships", memberships); err != nil {
		return err
	}
	d.Set("status", resp.Status)
	if resp.RaEnabled != nil {
		d.Set("ra_enabled", *resp.RaEnabled)
	} else {
		d.Set("ra_enabled", nil)
	}

	return nil
}

func resourcePublicNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("ip_blocks") {
		client := m.(receiver.BMCSDK)
		networkID := d.Id()
		query := &dto.Query{}
		var force = d.Get("force").(bool)
		query.Force = force
		oldInterface, newInterface := d.GetChange("ip_blocks")
		old := oldInterface.([]interface{})
		new := newInterface.([]interface{})

		var newIds []string
		if len(new) > 0 {
			for _, j := range new {
				ibItem := j.(map[string]interface{})
				pnib := ibItem["public_network_ip_block"].([]interface{})[0]
				pnibItem := pnib.(map[string]interface{})
				newId := pnibItem["id"].(string)
				newIds = append(newIds, newId)
			}
		}
		var oldIds []string
		if len(old) > 0 {
			for _, j := range old {
				ibItem := j.(map[string]interface{})
				pnib := ibItem["public_network_ip_block"].([]interface{})[0]
				pnibItem := pnib.(map[string]interface{})
				oldId := pnibItem["id"].(string)
				oldIds = append(oldIds, oldId)
			}
		}
		var sameIds []string
		var idExists bool
		for _, l := range newIds {
			idExists = false
			for _, n := range oldIds {
				if n == l {
					idExists = true
				}
			}
			if idExists {
				sameIds = append(sameIds, l)
			}
		}
		for _, p := range newIds {
			idExists = false
			for _, r := range sameIds {
				if p == r {
					idExists = true
				}
			}
			if !idExists {
				request := &networkapiclient.PublicNetworkIpBlockCreate{}
				request.Id = p
				requestCommand := publicnetwork.NewAddIpBlock2PublicNetworkCommand(client, networkID, *request)
				_, err := requestCommand.Execute()
				if err != nil {
					return err
				}
				waitResultError := ipBlockWaitForUnassign(p, &client)
				if waitResultError != nil {
					return waitResultError
				}
			}
		}
		for _, t := range oldIds {
			idExists = false
			for _, v := range sameIds {
				if t == v {
					idExists = true
				}
			}
			if !idExists {
				requestCommand := publicnetwork.NewRemoveIpBlockFromPublicNetworkCommandWithQuery(client, networkID, t, query)
				_, err := requestCommand.Execute()
				if err != nil {
					return err
				}
				waitResultError := ipBlockWaitForUnassign(t, &client)
				if waitResultError != nil {
					return waitResultError
				}
			}
		}
	} else if d.HasChange("name") || d.HasChange("description") {
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
	} else if d.HasChange("ra_enabled") {
		client := m.(receiver.BMCSDK)
		networkID := d.Id()
		request := &networkapiclient.PublicNetworkModify{}
		raEnabled := d.Get("ra_enabled").(bool)
		request.RaEnabled = &raEnabled

		requestCommand := publicnetwork.NewUpdatePublicNetworkCommand(client, networkID, *request)
		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
	} else if d.HasChange("force") {
		// Do nothing
	} else {
		return fmt.Errorf("unsupported action")
	}
	return resourcePublicNetworkRead(d, m)
}

func resourcePublicNetworkDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	networkID := d.Id()

	waitResultError := publicNetworkWaitForUnassign(networkID, &client)
	if waitResultError != nil {
		return waitResultError
	}

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

func publicNetworkWaitForUnassign(id string, client *receiver.BMCSDK) error {
	log.Printf("Waiting for public network %s to be unassigned...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"assigned"},
		Target:     []string{"unassigned"},
		Refresh:    refreshForPublicNetworkMembershipStatus(client, id),
		Timeout:    pnapPublicNetworkRetryTimeout,
		Delay:      pnapPublicNetworkRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for public network (%s) to be unassigned: %v", id, err)
	}

	return nil
}

func refreshForPublicNetworkMembershipStatus(client *receiver.BMCSDK, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		requestCommand := publicnetwork.NewGetPublicNetworkCommand(*client, id)

		resp, err := requestCommand.Execute()
		if err != nil {
			return 0, "", err
		} else if len(resp.Memberships) > 0 {
			return 0, "assigned", nil
		} else {
			return 0, "unassigned", nil
		}
	}
}

func flattenIpBlocks(pubNetIpBlock []networkapiclient.PublicNetworkIpBlock, ipBlocksInput []interface{}) []interface{} {
	if pubNetIpBlock != nil {
		var ib []interface{}
		var ipBlocksExists = false
		if ipBlocksInput != nil {
			ib = ipBlocksInput
			ipBlocksExists = true
		} else {
			ib = make([]interface{}, len(pubNetIpBlock))
			ipBlocksExists = false
		}
		for i, j := range pubNetIpBlock {
			for k := range ib {
				if !ipBlocksExists || ib[k].(map[string]interface{})["public_network_ip_block"].([]interface{})[0].(map[string]interface{})["id"] == j.Id {

					var ibItem map[string]interface{}
					var pnib []interface{}
					var pnibItem map[string]interface{}

					if !ipBlocksExists {
						ibItem = make(map[string]interface{})
						pnib = make([]interface{}, 1)
						pnibItem = make(map[string]interface{})
					} else {
						ibItem = ib[k].(map[string]interface{})
						pnib = ibItem["public_network_ip_block"].([]interface{})
						pnibItem = pnib[0].(map[string]interface{})
					}

					pnibItem["id"] = j.Id
					if len(j.Cidr) > 0 {
						pnibItem["cidr"] = j.Cidr
					}
					if len(j.UsedIpsCount) > 0 {
						pnibItem["used_ips_count"] = j.UsedIpsCount
					}
					if !ipBlocksExists {
						pnib[0] = pnibItem
						ibItem["public_network_ip_block"] = pnib
						ib[i] = ibItem
					}
				}
				if !ipBlocksExists {
					break
				}
			}
		}
		ipBlocksInput = ib
	}
	return ipBlocksInput
}
