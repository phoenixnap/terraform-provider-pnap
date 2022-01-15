package pnap

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/PNAP/go-sdk-helper-bmc/command/networkapi/privatenetwork"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"

	networkapiclient "github.com/phoenixnap/go-sdk-bmc/networkapi"
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
				Default:  false,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			/* "servers": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"ips": &schema.Schema{
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			}, */
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

	request.Cidr = d.Get("cidr").(string)

	requestCommand := privatenetwork.NewCreatePrivateNetworkCommand(client, *request)

	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	d.SetId(resp.Id)

	return resourcePrivateNetworkRead(d, m)
}

func resourcePrivateNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	keyID := d.Id()
	requestCommand := privatenetwork.NewGetPrivateNetworkCommand(client, keyID)
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
	return nil
}

func resourcePrivateNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("name") || d.HasChange("default") {
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
		return fmt.Errorf("unsuported action")
	}
	return resourcePrivateNetworkRead(d, m)

}

func resourcePrivateNetworkDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	sshKeyID := d.Id()

	requestCommand := privatenetwork.NewDeletePrivateNetworkCommand(client, sshKeyID)
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