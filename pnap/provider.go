package pnap

import (
	//"github.com/phoenixnap/go-sdk-bmc/client"
	"github.com/phoenixnap/go-sdk-bmc/dto"
	newClient "github.com/phoenixnap/go-sdk-bmc/client/pnapClient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider inits the root of provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": &schema.Schema{
			  Type:        schema.TypeString,
			  Optional:    true,
			  DefaultFunc: schema.EnvDefaultFunc("PNAP_CLIENT_ID", nil),
			},
			"client_secret": &schema.Schema{
			  Type:        schema.TypeString,
			  Optional:    true,
			  Sensitive:   true,
			  DefaultFunc: schema.EnvDefaultFunc("PNAP_CLIENT_SECRET", nil),
			},
			"config_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
		  },
		ResourcesMap: map[string]*schema.Resource{
			"pnap_ssh_key": resourceSshKey(),
			"pnap_server": resourceServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	configFilePath := d.Get("config_file_path").(string)

	if (clientId != "") && (clientSecret != ""){
		auth := dto.Authentication{ClientID : clientId,
		ClientSecret: clientSecret,
		TokenURL: "https://auth.phoenixnap.com/auth/realms/BMC/protocol/openid-connect/token",
		ApiHostName:"https://api.phoenixnap.com/",
		PoweredBy: "terraform-provider-pnap"}
		cl := newClient.NewPNAPClient(auth)
		return cl, nil
	}

	if (configFilePath != ""){
	
		cl, confErr := newClient.NewPNAPClientWithCustomConfig(configFilePath)
		if confErr == nil {
			auth := dto.Authentication{ClientID : "",
			ClientSecret: "",
			TokenURL: "",
			ApiHostName:"",
			PoweredBy: "terraform-provider-pnap"}
			cl.SetAuthentication(auth)
		}
		return cl, confErr
	}

	client, confErr := newClient.NewPNAPClientWithDefaultConfig()
	if confErr == nil {
		auth := dto.Authentication{ClientID : "",
		ClientSecret: "",
		TokenURL: "",
		ApiHostName:"",
		PoweredBy: "terraform-provider-pnap"}
		client.SetAuthentication(auth)
	}
	return client, confErr
}
