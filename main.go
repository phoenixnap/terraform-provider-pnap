package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/phoenixnap/terraform-provider-pnap/pnap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return pnap.Provider()
		},
	})
}
