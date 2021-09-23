package pnap

import (
	"strings"
	"fmt"
	"log"
	"time"
	"encoding/json"

	//"github.com/phoenixnap/go-sdk-bmc/client"
	"github.com/phoenixnap/go-sdk-bmc/command"
	"github.com/phoenixnap/go-sdk-bmc/dto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/phoenixnap/go-sdk-bmc/client/pnapClient"

)

const (
	pnapRetryTimeout       = 30 * time.Minute
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
				Optional: true,
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
			"cpu_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cores_per_cpu": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_frequency_in_ghz": &schema.Schema{
				Type:     schema.TypeInt,
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
			"network_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"install_default_ssh_keys": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssh_key_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"reservation_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pricing_model": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rdp_allowed_ips": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Sensitive: true,
			},
			"cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"management_ui_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_password": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				//Sensitive: true,
			},
			"management_access_allowed_ips": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"provisioned_on": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(client.PNAPClient)

	request := &dto.ProvisionedServer{}
	request.Name = d.Get("hostname").(string)
	request.Description = d.Get("description").(string)
	request.Os = d.Get("os").(string)
	request.Type = d.Get("type").(string)
	request.Location = d.Get("location").(string)
	request.NetworkType = d.Get("network_type").(string)

	request.ReservationId = d.Get("reservation_id").(string)
	request.PricingModel = d.Get("pricing_model").(string)

	request.InstallDefaultSshKeys = d.Get("install_default_ssh_keys").(bool)
	temp := d.Get("ssh_keys").(*schema.Set).List()
	keys := make([]string, len(temp))
	for i, v := range temp {
		keys[i] = fmt.Sprint(v)
	}
	request.SshKeys = keys

	temp1 := d.Get("ssh_key_ids").(*schema.Set).List()
	keyIds := make([]string, len(temp1))
	for i, v := range temp1 {
		keyIds[i] = fmt.Sprint(v)
	}
	request.SshKeyIds = keyIds


	dtoWindows := dto.Windows{}

	temp2 := d.Get("rdp_allowed_ips").(*schema.Set).List()
	allowedIps := make([]string, len(temp2))
	for i, v := range temp2 {
		allowedIps[i] = fmt.Sprint(v)
	}

	dtoWindows.RdpAllowedIps = allowedIps
	dtoOsConfiguration := dto.OsConfiguration{}
	dtoOsConfiguration.Windows = &dtoWindows
	request.OsConfiguration = dtoOsConfiguration

	temp3 := d.Get("management_access_allowed_ips").(*schema.Set).List()
	managementAccessAllowedIps := make([]string, len(temp3))
	for i, v := range temp3 {
		managementAccessAllowedIps[i] = fmt.Sprint(v)
	}
	request.OsConfiguration.ManagementAccessAllowedIps = managementAccessAllowedIps

	requestCommand := command.NewCreateServerCommand(client, *request)

	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	code := resp.StatusCode
	if code == 200 {
		response := &dto.LongServer{}
		response.FromBytes(resp)
		d.SetId(response.ID)
		d.Set("password", response.Password)
		if(&response.OsConfiguration != nil){
			d.Set("root_password", response.OsConfiguration.RootPassword)
			d.Set("management_ui_url", response.OsConfiguration.ManagementUiUrl)
		}

		waitResultError := resourceWaitForCreate(response.ID, &client)
		if waitResultError != nil {
			return waitResultError
		}
	} else {
		response := &dto.ErrorMessage{}
		response.FromBytes(resp)
		return fmt.Errorf("API create server Returned Code %v Message: %s Validation Errors: %s", code, response.Message, response.ValidationErrors)
	}

	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(client.PNAPClient)
	serverID := d.Id()
	requestCommand := command.NewGetServerCommand(client, serverID)
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
	d.Set("cpu_count", response.CPUCount)
	d.Set("cores_per_cpu", response.CoresPerCpu)
	d.Set("cpu_frequency_in_ghz", response.CPUFrequency)
	d.Set("ram", response.RAM)
	d.Set("storage", response.Storage)
	d.Set("network_type", response.NetworkType)
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
	d.Set("reservation_id", response.ReservationID)
	d.Set("pricing_model", response.PricingModel)
	
	d.Set("cluster_id", response.ClusterID)
	if(&response.OsConfiguration != nil && response.OsConfiguration.ManagementAccessAllowedIps != nil){
		var mgmntAccessAllowedIps []interface{}
		for _, k := range response.OsConfiguration.ManagementAccessAllowedIps {
			mgmntAccessAllowedIps = append(mgmntAccessAllowedIps, k)
		}
		d.Set("management_access_allowed_ips", mgmntAccessAllowedIps)
	}

	if(&response.OsConfiguration != nil && response.OsConfiguration.Windows != nil && response.OsConfiguration.Windows.RdpAllowedIps != nil){
		var rdpAllowedIps []interface{}
		for _, k := range response.OsConfiguration.Windows.RdpAllowedIps {
			rdpAllowedIps = append(rdpAllowedIps, k)
		}
		d.Set("rdp_allowed_ips", rdpAllowedIps)
	}

	d.Set("provisioned_on", response.ProvisionedOn)

	
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
			request.SshKeys = keys
			request.InstallDefaultSshKeys = d.Get("install_default_ssh_keys").(bool)


			temp1 := d.Get("ssh_key_ids").(*schema.Set).List()
			keyIds := make([]string, len(temp1))
			for i, v := range temp1 {
				keyIds[i] = fmt.Sprint(v)
			}
			request.SshKeyIds = keyIds


			dtoOsConfiguration := dto.OsConfiguration{}
			isWindows:= strings.Contains(d.Get("os").(string), "windows")
			isEsxi:= strings.Contains(d.Get("os").(string), "esxi")
  
			if(isWindows){
				//log.Printf("Waiting for server windows to be reseted...")
				dtoWindows := dto.Windows{}
				temp2 := d.Get("rdp_allowed_ips").(*schema.Set).List()
			    allowedIps := make([]string, len(temp2))
			    for i, v := range temp2 {
				   allowedIps[i] = fmt.Sprint(v)
			    }

			     dtoWindows.RdpAllowedIps = allowedIps
				 dtoOsConfiguration.Windows = &dtoWindows
				 dtoOsConfiguration.Esxi = nil
			}

            if(isEsxi){
				//log.Printf("Waiting for server esxi to be reseted...")
				dtoEsxi := dto.Esxi{}
				temp3 := d.Get("management_access_allowed_ips").(*schema.Set).List()
	            managementAccessAllowedIps := make([]string, len(temp3))
	            for i, v := range temp3 {
		          managementAccessAllowedIps[i] = fmt.Sprint(v)
	            }
	            dtoEsxi.ManagementAccessAllowedIps = managementAccessAllowedIps
				dtoOsConfiguration.Esxi = &dtoEsxi
				dtoOsConfiguration.Windows = nil
				
			}
			
			request.OsConfiguration = dtoOsConfiguration
			//b, err := json.MarshalIndent(request, "", "  ")
			//log.Printf("request object is" + string(b))
			request.ID = d.Id()
			requestCommand = command.NewResetCommand(client, *request)
			err, resp := runResetCommand(requestCommand)
			if err != nil {
				return err
			}
			d.Set("password", resp.Password)

			 if(&resp.OsConfiguration != nil && resp.OsConfiguration.Esxi != nil){
				d.Set("root_password", resp.OsConfiguration.Esxi.RootPassword)
				d.Set("management_ui_url", resp.OsConfiguration.Esxi.ManagementUiUrl)
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

	} else if d.HasChange("pricing_model"){
		client := m.(client.PNAPClient)
		var requestCommand command.Executor
		//reserve action
		request := &dto.ProvisionedServer{}
		request.ID = d.Id()
		request.PricingModel = d.Get("pricing_model").(string)

		requestCommand = command.NewReserveCommand(client, *request)
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
			return 0, "", fmt.Errorf("API refressh for create Returned Code: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors)
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

func runResetCommand(command command.Executor) (error, dto.ServerActionResponse) {
	resp, err := command.Execute()
	if err != nil {
		return err, dto.ServerActionResponse{}
	}
	code := resp.StatusCode
	if code != 200 {
		response := &dto.ErrorMessage{}
		response.FromBytes(resp)
		return fmt.Errorf("API Returned Code: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors), dto.ServerActionResponse{}
	}
	if code == 200 {
		response := &dto.ServerActionResponse{}
		response.FromBytes(resp)
		return nil, *response
	}
	return nil, dto.ServerActionResponse{}
}
