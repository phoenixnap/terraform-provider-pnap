package pnap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"pnap": testAccProvider,
	}
}
