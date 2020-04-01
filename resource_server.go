package main

import (
	"fmt"
	"log"
	"time"

	"github.com/PNAP/bmc-api-sdk/client"
	"github.com/PNAP/bmc-api-sdk/command"
	"github.com/PNAP/bmc-api-sdk/dto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	pnapRetryTimeout       = 15 * time.Minute
	pnapDeleteRetryTimeout = 15 * time.Minute
	pnapRetryDelay         = 5 * time.Second
	pnapRetryMinTimeout    = 3 * time.Second
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(pnapRetryTimeout),
			Update: schema.DefaultTimeout(pnapRetryTimeout),
			Delete: schema.DefaultTimeout(pnapDeleteRetryTimeout),
		},

		Schema: map[string]*schema.Schema{
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"public": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
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
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ssh_keys": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cpu": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {

	//client, confErr := client.Create()
	client := m.(client.PNAPClient)
	/* if confErr != nil {
		return confErr
	} */
	request := &dto.ProvisionedServer{}
	request.Name = d.Get("hostname").(string)
	request.Description = d.Get("description").(string)
	request.Os = d.Get("os").(string)
	request.Type = d.Get("type").(string)
	request.Location = d.Get("location").(string)
	request.Public = d.Get("public").(bool)
	temp := d.Get("ssh_keys").(*schema.Set).List()
	keys := make([]string, len(temp))
	for i, v := range temp {
		keys[i] = fmt.Sprint(v)
	}
	request.SSHKeys = keys

	requestCommand := &command.CreateServerCommand{}
	requestCommand.SetRequester(client)
	requestCommand.SetServer(*request)

	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	code := resp.StatusCode
	if code == 200 {
		response := &dto.LongServer{}
		response.FromBytes(resp)
		d.SetId(response.ID)
		waitResultError := resourceWaitForCreate(response.ID, &client)
		if waitResultError != nil {
			return waitResultError
		}
	} else {
		response := &dto.ErrorMessage{}
		response.FromBytes(resp)
		return fmt.Errorf("API Returned Code %v Message: %s Validation Errors: %s", code, response.Message, response.ValidationErrors)
	}

	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(client.PNAPClient)

	requestCommand := &command.GetServerCommand{}
	requestCommand.SetRequester(client)
	serverID := d.Id()
	requestCommand.SetServerID(serverID)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	code := resp.StatusCode
	if code != 200 {
		response := &dto.ErrorMessage{}
		response.FromBytes(resp)
		return fmt.Errorf("API Returned Code: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors)
	}
	response := &dto.LongServer{}
	response.FromBytes(resp)
	d.SetId(response.ID)
	d.Set("status", response.Status)
	d.Set("hostname", response.Name)
	d.Set("description", response.Description)
	d.Set("os", response.Os)
	d.Set("type", response.Type)
	d.Set("location", response.Location)
	d.Set("cpu", response.CPU)
	d.Set("ram", response.RAM)
	d.Set("action", "")
	var privateIPs []interface{}
	for _, v := range response.PrivateIPAddresses {
		privateIPs = append(privateIPs, v)
	}
	d.Set("private_ip_addresses", privateIPs)
	var publicIPs []interface{}
	for _, k := range response.PublicIPAddresses {
		publicIPs = append(publicIPs, k)
	}
	d.Set("public_ip_addresses", publicIPs)
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("action") {
		client := m.(client.PNAPClient)
		var requestCommand command.Executor
		newStatus := d.Get("action").(string)

		switch newStatus {
		case "powered-on":
			//do power-on request
			serverID := d.Id()
			requestCommand = command.NewPowerOnCommand(client, serverID)
			err := run(requestCommand)
			if err != nil {
				return err
			}
			waitResultError := resourceWaitForPowerON(d.Id(), &client)
			if waitResultError != nil {
				return waitResultError
			}
		case "powered-off":
			//power off request

			serverID := d.Id()

			requestCommand = command.NewPowerOffCommand(client, serverID)
			err := run(requestCommand)
			if err != nil {
				return err
			}
			waitResultError := resourceWaitForPowerOff(d.Id(), &client)
			if waitResultError != nil {
				return waitResultError
			}
		case "reboot":
			//reboot

			serverID := d.Id()

			requestCommand = command.NewRebootCommand(client, serverID)
			err := run(requestCommand)
			if err != nil {
				return err
			}
			waitResultError := resourceWaitForCreate(d.Id(), &client)
			if waitResultError != nil {
				return waitResultError
			}
		case "reset":
			//reset
			request := &dto.ProvisionedServer{}
			temp := d.Get("ssh_keys").(*schema.Set).List()
			keys := make([]string, len(temp))
			for i, v := range temp {
				keys[i] = fmt.Sprint(v)
			}
			request.SSHKeys = keys
			request.ID = d.Id()

			requestCommand = command.NewResetCommand(client, *request)
			err := run(requestCommand)
			if err != nil {
				return err
			}
			waitResultError := resourceWaitForCreate(d.Id(), &client)
			if waitResultError != nil {
				return waitResultError
			}

		case "shutdown":

			serverID := d.Id()

			requestCommand = command.NewShutDownCommand(client, serverID)
			err := run(requestCommand)
			if err != nil {
				return err
			}
			waitResultError := resourceWaitForPowerOff(d.Id(), &client)
			if waitResultError != nil {
				return waitResultError
			}

		case "default":
			return fmt.Errorf("Unsuported action")
		}

	} else {
		return fmt.Errorf("Unsuported action")
	}
	return resourceServerRead(d, m)

}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(client.PNAPClient)

	serverID := d.Id()

	requestCommand := command.NewDeleteServerCommand(client, serverID)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	code := resp.StatusCode
	if code != 200 && code != 404 {
		response := &dto.ErrorMessage{}
		response.FromBytes(resp)
		return fmt.Errorf("API Returned Code: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors)
	}
	return nil
}

func resourceWaitForCreate(id string, client *client.PNAPClient) error {
	log.Printf("Waiting for server %s to be created...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "resetting", "rebooting"},
		Target:     []string{"powered-on", "powered-off"},
		Refresh:    refreshForCreate(client, id),
		Timeout:    pnapRetryTimeout,
		Delay:      pnapRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for server (%s) to switch to target state: %v", id, err)
	}

	return nil
}

func resourceWaitForPowerON(id string, client *client.PNAPClient) error {
	log.Printf("Waiting for server %s to power on...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"powered-off"},
		Target:     []string{"powered-on"},
		Refresh:    refreshForCreate(client, id),
		Timeout:    pnapRetryTimeout,
		Delay:      pnapRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for server (%s) to power on: %v", id, err)
	}

	return nil
}

func resourceWaitForPowerOff(id string, client *client.PNAPClient) error {
	log.Printf("Waiting for server %s to power off...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"powered-on"},
		Target:     []string{"powered-off"},
		Refresh:    refreshForCreate(client, id),
		Timeout:    pnapRetryTimeout,
		Delay:      pnapRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for server (%s) to power off: %v", id, err)
	}

	return nil
}

func refreshForCreate(client *client.PNAPClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		requestCommand := &command.GetServerCommand{}
		requestCommand.SetRequester(client)
		serverID := id
		requestCommand.SetServerID(serverID)
		resp, err := requestCommand.Execute()
		if err != nil {
			return 0, "", err
		}
		code := resp.StatusCode
		if code != 200 {
			response := &dto.ErrorMessage{}
			response.FromBytes(resp)
			return 0, "", fmt.Errorf("API Returned Code: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors)
		}
		response := &dto.LongServer{}
		response.FromBytes(resp)
		return 0, response.Status, nil
	}
}

func run(command command.Executor) error {
	resp, err := command.Execute()
	if err != nil {
		return err
	}
	code := resp.StatusCode
	if code != 200 {
		response := &dto.ErrorMessage{}
		response.FromBytes(resp)
		return fmt.Errorf("API Returned Code: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors)
	}
	return nil
}
