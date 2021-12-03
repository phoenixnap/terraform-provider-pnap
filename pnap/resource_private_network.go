package pnap

import (
	"fmt"

	//"github.com/phoenixnap/go-sdk-bmc/command"
	//"github.com/phoenixnap/go-sdk-bmc/dto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//client "github.com/phoenixnap/go-sdk-bmc/client/pnapClient"

	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/PNAP/go-sdk-helper-bmc/command/networkapi/privatenetwork"
	//helpercommand "github.com/PNAP/go-sdk-helper-bmc/command"
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
			
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"location_default": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			/* "servers": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
				  Schema: map[string]*schema.Schema{
					"id": &schema.Schema{
					  Type:     schema.TypeInt,
					  Computed: true,
					},
					"ips": &schema.Schema{
						Type:     schema.TypeSet,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				  },
				}, 
			},*/	
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
	if (len(desc) > 0){
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
	//resp.Servers
	
	return nil
}

func resourcePrivateNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("name") || d.HasChange("default") {
		client := m.(receiver.BMCSDK)
		
		request := &networkapiclient.PrivateNetworkModify{}
		request.Name = d.Get("name").(string)
		request.LocationDefault = d.Get("location_default").(bool)
		var desc = d.Get("description").(string)
		if (len(desc) > 0){
			request.Description = &desc
		}
		requestCommand := privatenetwork.NewUpdatePrivateNetworkCommand(client, d.Id(), *request)

		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
		
	}  else {
		return fmt.Errorf("Unsuported action")
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
