package main

import (
	"github.com/PNAP/bmc-api-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider inits the root of provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"pnap_server": resourceServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client, confErr := client.Create()
	if confErr != nil {
		return client, confErr
	}
	return client, nil
}
