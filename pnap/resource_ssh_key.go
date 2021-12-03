package pnap

import (
	"fmt"

	//"github.com/phoenixnap/go-sdk-bmc/command"
	//"github.com/phoenixnap/go-sdk-bmc/dto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//client "github.com/phoenixnap/go-sdk-bmc/client/pnapClient"

	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/PNAP/go-sdk-helper-bmc/command/bmcapi/sshkey"
	//helpercommand "github.com/PNAP/go-sdk-helper-bmc/command"
	bmcapiclient "github.com/phoenixnap/go-sdk-bmc/bmcapi"

)



func resourceSshKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSshKeyCreate,
		Read:   resourceSshKeyRead,
		Update: resourceSshKeyUpdate,
		Delete: resourceSshKeyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(pnapRetryTimeout),
			Update: schema.DefaultTimeout(pnapRetryTimeout),
			Delete: schema.DefaultTimeout(pnapDeleteRetryTimeout),
		},

		Schema: map[string]*schema.Schema{
			"default": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_on": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_on": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSshKeyCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(receiver.BMCSDK)

	request := &bmcapiclient.SshKeyCreate{}
	request.Name = d.Get("name").(string)
	request.Default = d.Get("default").(bool)
	request.Key = d.Get("key").(string)
	

	requestCommand := sshkey.NewCreateSshKeyCommand(client, *request)

	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	//code := resp.StatusCode
	//if code == 201 {
		//response := &dto.SshKey{}
		//response.FromBytes(resp)
		d.SetId(resp.Id)
		
	/* } else {
		response := &dto.ErrorMessage{}
		response.FromBytes(resp)
		return fmt.Errorf("API Returned Code %v Message: %s Validation Errors: %s", code, response.Message, response.ValidationErrors)
	} */

	return resourceSshKeyRead(d, m)
}

func resourceSshKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	keyID := d.Id()
	requestCommand := sshkey.NewGetSshKeyCommand(client, keyID)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	/* code := resp.StatusCode
	if code != 200 {
		response := &dto.ErrorMessage{}
		response.FromBytes(resp)
		return fmt.Errorf("API Returned Code from read method: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors)
	} */
	//response := &dto.SshKey{}
	//response.FromBytes(resp)
	d.SetId(resp.Id)
	d.Set("default", resp.Default)
	d.Set("name", resp.Name)
	d.Set("key", resp.Key)
	d.Set("fingerprint", resp.Fingerprint)
	d.Set("created_on", resp.CreatedOn.String())
	d.Set("last_updated_on", resp.LastUpdatedOn.String())
	
	return nil
}

func resourceSshKeyUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("name") || d.HasChange("default") {
		client := m.(receiver.BMCSDK)
		//var requestCommand command.Executor
		
		request := &bmcapiclient.SshKeyUpdate{}
		request.Name = d.Get("name").(string)
		request.Default = d.Get("default").(bool)
		//request.Id = d.Id()
		requestCommand := sshkey.NewUpdateSshKeyCommand(client, d.Id(), *request)

		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
		/* code := resp.StatusCode
		if code != 200 {
			response := &dto.ErrorMessage{}
			response.FromBytes(resp)
			return fmt.Errorf("API Returned Code %v Message: %s Validation Errors: %s", code, response.Message, response.ValidationErrors)
			
		} */
	}  else {
		return fmt.Errorf("Unsuported action")
	}
	return resourceSshKeyRead(d, m)

}

func resourceSshKeyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	sshKeyID := d.Id()

	requestCommand := sshkey.NewDeleteSshKeyCommand(client, sshKeyID)
	_, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	/* code := resp.StatusCode
	if code != 200 && code != 404 {
		response := &dto.ErrorMessage{}
		response.FromBytes(resp)
		return fmt.Errorf("API Returned Code: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors)
	} */
	return nil
}
