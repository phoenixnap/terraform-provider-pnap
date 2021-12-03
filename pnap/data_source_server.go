package pnap

import (
	"fmt"

	//"github.com/phoenixnap/go-sdk-bmc/client"
	//"github.com/phoenixnap/go-sdk-bmc/command"
	//"github.com/phoenixnap/go-sdk-bmc/dto"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//client "github.com/phoenixnap/go-sdk-bmc/client/pnapClient"

	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/PNAP/go-sdk-helper-bmc/command/bmcapi/server"
	//bmcapiclient "github.com/phoenixnap/go-sdk-bmc/bmcapi"
)

func dataSourceServer() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceServerRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"hostname"},
			},
			"hostname": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ip_addresses": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public_ip_addresses": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"os": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	//serverID := d.Id()
	requestCommand := server.NewGetServersCommand(client)
	//requestCommand.SetRequester(client)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	/* code := resp.StatusCode
	if code != 200 {
		response := &dto.ErrorMessage{}
		response.FromBytes(resp)
		return fmt.Errorf("API Returned Code: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors)
	}
	response := &dto.Servers{}
	response.FromBytes(resp) */

	numOfServers := 0
	for _, instance := range resp {
		if instance.Hostname == d.Get("hostname").(string) || instance.Id == d.Get("id").(string) {
			numOfServers++
			d.SetId(instance.Id)
			d.Set("status", instance.Status)
			d.Set("hostname", instance.Hostname)

			d.Set("os", instance.Os)
			d.Set("type", instance.Type)
			d.Set("location", instance.Location)

			var privateIPs []interface{}
			for _, v := range instance.PrivateIpAddresses {
				privateIPs = append(privateIPs, v)
			}
			d.Set("private_ip_addresses", privateIPs)
			var publicIPs []interface{}
			for _, k := range instance.PublicIpAddresses {
				publicIPs = append(publicIPs, k)
			}
			d.Set("public_ip_addresses", publicIPs)
			d.Set("primary_ip_address", instance.PublicIpAddresses[0])

		}
	}

	if numOfServers > 1 {
		return fmt.Errorf("Too many devices found with hostname %s (found %d, expected 1)", d.Get("hostname").(string), numOfServers)
	}

	return nil
}
