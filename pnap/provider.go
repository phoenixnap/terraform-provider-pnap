package pnap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/PNAP/go-sdk-helper-bmc/dto"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
)

// Provider inits the root of provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PNAP_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("PNAP_CLIENT_SECRET", nil),
			},
			"config_file_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"pnap_ssh_key":         resourceSshKey(),
			"pnap_server":          resourceServer(),
			"pnap_private_network": resourcePrivateNetwork(),
			"pnap_reservation":     resourceReservation(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"pnap_ssh_key":         dataSourceSshKey(),
			"pnap_server":          dataSourceServer(),
			"pnap_private_network": dataSourcePrivateNetwork(),
			"pnap_reservation":     dataSourceReservation(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	configFilePath := d.Get("config_file_path").(string)

	configuration := dto.Configuration{}
	configuration.UserAgent = "terraform-provider-pnap"
	if (clientId != "") && (clientSecret != "") {
		configuration.ClientID = clientId
		configuration.ClientSecret = clientSecret
		configuration.TokenURL = "https://auth.phoenixnap.com/auth/realms/BMC/protocol/openid-connect/token"
		configuration.ApiHostName = "https://api.phoenixnap.com/"
		/* auth := dto.Authentication{ClientID : clientId,
		ClientSecret: clientSecret,
		TokenURL: "https://auth.phoenixnap.com/auth/realms/BMC/protocol/openid-connect/token",
		ApiHostName:"https://api.phoenixnap.com/",
		PoweredBy: "terraform-provider-pnap"}
		cl := newClient.NewPNAPClient(auth) */
		cl := receiver.NewBMCSDK(configuration)
		return cl, nil
	}

	if configFilePath != "" {

		cl, confErr := receiver.NewBMCSDKWithCustomConfig(configFilePath, configuration)
		/* if confErr == nil {
			auth := dto.Authentication{ClientID : "",
			ClientSecret: "",
			TokenURL: "",
			ApiHostName:"",
			PoweredBy: "terraform-provider-pnap"}
			cl.SetAuthentication(auth)
		} */
		return cl, confErr
	}

	client, confErr := receiver.NewBMCSDKWithDefaultConfig(configuration)
	/* if confErr == nil {
		auth := dto.Authentication{ClientID : "",
		ClientSecret: "",
		TokenURL: "",
		ApiHostName:"",
		PoweredBy: "terraform-provider-pnap"}
		client.SetAuthentication(auth)
	} */
	return client, confErr
}
