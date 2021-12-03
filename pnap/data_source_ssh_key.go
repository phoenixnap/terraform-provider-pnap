package pnap

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	/* client "github.com/phoenixnap/go-sdk-bmc/client/pnapClient"
	"github.com/phoenixnap/go-sdk-bmc/command"
	"github.com/phoenixnap/go-sdk-bmc/dto" */

	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/PNAP/go-sdk-helper-bmc/command/bmcapi/sshkey"
	
)

func dataSourceSshKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSshKeyRead,

		Schema: map[string]*schema.Schema{
			"default": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSshKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	requestCommand := sshkey.NewGetSshKeysCommand(client)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	/* code := resp.StatusCode
	if code != 200 {
		response := &dto.ErrorMessage{}
		response.FromBytes(resp)
		return fmt.Errorf("API Returned Code from read method: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors)
	}
	response := &dto.SshKeys{}
	response.FromBytes(resp) */

	numOfKeys := 0
	for _, instance := range resp {
		if instance.Name == d.Get("name").(string) {
			numOfKeys++
			d.SetId(instance.Id)
			d.Set("default", instance.Default)
			d.Set("name", instance.Name)
			d.Set("key", instance.Key)

		}
	}
	if numOfKeys > 1 {
		return fmt.Errorf("Too many ssh keys with name %s (found %d, expected 1)", d.Get("name").(string), numOfKeys)
	}

	return nil
}
