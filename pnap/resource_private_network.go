package pnap

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/PNAP/go-sdk-helper-bmc/command/networkapi/privatenetwork"
	"github.com/PNAP/go-sdk-helper-bmc/dto"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"

	networkapiclient "github.com/phoenixnap/go-sdk-bmc/networkapi/v2"
)

const (
	pnapPrivateNetworkRetryDelay   = 10 * time.Second
	pnapPrivateNetworkRetryTimeout = 7 * time.Minute
)

func resourcePrivateNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrivateNetworkCreate,
		Read:   resourcePrivateNetworkRead,
		Update: resourcePrivateNetworkUpdate,
		Delete: resourcePrivateNetworkDelete,

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
			"location_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"servers": { // Deprecated
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
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePrivateNetworkCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(receiver.BMCSDK)

	request := &networkapiclient.PrivateNetworkCreate{}
	request.Name = d.Get("name").(string)
	request.Location = d.Get("location").(string)
	var locDefault = d.Get("location_default").(bool)

	request.LocationDefault = &locDefault

	var desc = d.Get("description").(string)
	if len(desc) > 0 {
		request.Description = &desc
	}

	cidr := d.Get("cidr").(string)
	request.Cidr = &cidr
	var vlanId = d.Get("vlan_id").(int)
	if vlanId > 0 {
		vlanId32 := int32(vlanId)
		request.VlanId = &vlanId32
	}

	query := &dto.Query{}
	query.Force = d.Get("force").(bool)

	requestCommand := privatenetwork.NewCreatePrivateNetworkCommand(client, *request, *query)

	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	d.SetId(resp.Id)

	return resourcePrivateNetworkRead(d, m)
}

func resourcePrivateNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	networkID := d.Id()
	requestCommand := privatenetwork.NewGetPrivateNetworkCommand(client, networkID)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	d.SetId(resp.Id)
	d.Set("location", resp.Location)
	d.Set("name", resp.Name)
	d.Set("cidr", resp.Cidr)
	d.Set("description", resp.Description)
	d.Set("location_default", resp.LocationDefault)
	d.Set("type", resp.Type)
	d.Set("vlan_id", resp.VlanId)

	servers := flattenServers(resp.Servers)
	if err := d.Set("servers", servers); err != nil {
		return err
	}
	memberships := flattenMemberships(resp.Memberships)
	if err := d.Set("memberships", memberships); err != nil {
		return err
	}
	d.Set("status", resp.Status)

	if len(resp.CreatedOn.String()) > 0 {
		d.Set("created_on", resp.CreatedOn.String())
	}

	return nil
}

func resourcePrivateNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("name") || d.HasChange("location_default") || d.HasChange("description") {
		client := m.(receiver.BMCSDK)

		request := &networkapiclient.PrivateNetworkModify{}
		request.Name = d.Get("name").(string)
		request.LocationDefault = d.Get("location_default").(bool)
		var desc = d.Get("description").(string)
		if len(desc) > 0 {
			request.Description = &desc
		}
		requestCommand := privatenetwork.NewUpdatePrivateNetworkCommand(client, d.Id(), *request)

		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("unsupported action")
	}
	return resourcePrivateNetworkRead(d, m)

}

func resourcePrivateNetworkDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	networkID := d.Id()

	waitResultError := privateNetworkWaitForUnassign(networkID, &client)
	if waitResultError != nil {
		return waitResultError
	}

	requestCommand := privatenetwork.NewDeletePrivateNetworkCommand(client, networkID)
	err := requestCommand.Execute()
	if err != nil {
		return err
	}

	return nil
}

func flattenServers(servers []networkapiclient.PrivateNetworkServer) []interface{} {
	if servers != nil {
		ss := make([]interface{}, len(servers))
		for i, v := range servers {
			s := make(map[string]interface{})

			privateIPs := make([]interface{}, len(v.Ips))
			for j, k := range v.Ips {
				privateIPs[j] = k
			}
			s["ips"] = privateIPs
			s["id"] = v.Id
			ss[i] = s
		}
		return ss
	}
	return make([]interface{}, 0)
}

func privateNetworkWaitForUnassign(id string, client *receiver.BMCSDK) error {
	log.Printf("Waiting for private network %s to be unassigned...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"assigned"},
		Target:     []string{"unassigned"},
		Refresh:    refreshForPrivateNetworkMembershipStatus(client, id),
		Timeout:    pnapPrivateNetworkRetryTimeout,
		Delay:      pnapPrivateNetworkRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for private network (%s) to be unassigned: %v", id, err)
	}

	return nil
}

func refreshForPrivateNetworkMembershipStatus(client *receiver.BMCSDK, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		requestCommand := privatenetwork.NewGetPrivateNetworkCommand(*client, id)

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
