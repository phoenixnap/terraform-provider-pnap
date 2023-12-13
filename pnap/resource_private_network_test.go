package pnap

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	helperprivatenetwork "github.com/PNAP/go-sdk-helper-bmc/command/networkapi/privatenetwork"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	networkapiclient "github.com/phoenixnap/go-sdk-bmc/networkapi/v2"
)

func TestAccPnapPrivateNetwork_basic(t *testing.T) {

	var privateNetwork networkapiclient.PrivateNetwork
	// generate a random name for each widget test run, to avoid
	// collisions from multiple concurrent tests.
	// the acctest package includes many helpers such as RandStringFromCharSet
	// See https://godoc.org/github.com/hashicorp/terraform-plugin-sdk/helper/acctest
	rNameSuffix := acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)
	rName := "acctest-" + rNameSuffix
	rName2 := "acctest-" + rNameSuffix + "-basic"
	rLine := "pnap_private_network." + rName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrivateNetworkResourceDestroy,
		Steps: []resource.TestStep{
			{
				// use configuration for private network creation
				Config: testAccCreatePrivateNetworkResource_basic(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the private network object
					testAccCheckPrivateNetworkExists(rLine, &privateNetwork),
					// verify remote values
					testAccCheckPrivateNetworkAttributes_basic(rName, &privateNetwork),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "location_default", "false"),
					resource.TestCheckResourceAttr(rLine, "cidr", "10.0.1.0/31"),
					resource.TestCheckResourceAttr(rLine, "status", "READY"),
					resource.TestCheckResourceAttrSet(rLine, "vlan_id"),
					resource.TestCheckResourceAttrSet(rLine, "type"),
					resource.TestCheckResourceAttrSet(rLine, "created_on"),
				),
			},
			{
				// update private network details
				Config: testAccUpdatePrivateNetworkResource_basic(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the private network object
					testAccCheckPrivateNetworkExists(rLine, &privateNetwork),
					// verify remote values
					testAccCheckPrivateNetworkNewAttributes_basic(rName2, &privateNetwork),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "location_default", "false"),
					resource.TestCheckResourceAttr(rLine, "cidr", "10.0.1.0/31"),
					resource.TestCheckResourceAttr(rLine, "status", "READY"),
					resource.TestCheckResourceAttrSet(rLine, "vlan_id"),
					resource.TestCheckResourceAttrSet(rLine, "type"),
					resource.TestCheckResourceAttrSet(rLine, "created_on"),
				),
			},
		},
	})
}

func TestAccPnapPrivateNetwork_force(t *testing.T) {

	var privateNetwork networkapiclient.PrivateNetwork
	rNameSuffix := acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)
	rName := "acctest-" + rNameSuffix
	rName2 := "acctest-" + rNameSuffix + "-force"
	rLine := "pnap_private_network." + rName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrivateNetworkResourceDestroy,
		Steps: []resource.TestStep{
			{
				// use configuration for private network creation
				Config: testAccCreatePrivateNetworkResource_force(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the private network object
					testAccCheckPrivateNetworkExists(rLine, &privateNetwork),
					// verify remote values
					testAccCheckPrivateNetworkAttributes_force(rName, &privateNetwork),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "location_default", "false"),
					resource.TestCheckResourceAttr(rLine, "cidr", ""),
					resource.TestCheckResourceAttr(rLine, "status", "READY"),
					resource.TestCheckResourceAttrSet(rLine, "vlan_id"),
					resource.TestCheckResourceAttrSet(rLine, "type"),
					resource.TestCheckResourceAttrSet(rLine, "created_on"),
				),
			},
			{
				// update private network details
				Config: testAccUpdatePrivateNetworkResource_force(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the private network object
					testAccCheckPrivateNetworkExists(rLine, &privateNetwork),
					// verify remote values
					testAccCheckPrivateNetworkNewAttributes_force(rName2, &privateNetwork),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "location_default", "false"),
					resource.TestCheckResourceAttr(rLine, "cidr", ""),
					resource.TestCheckResourceAttr(rLine, "status", "READY"),
					resource.TestCheckResourceAttrSet(rLine, "vlan_id"),
					resource.TestCheckResourceAttrSet(rLine, "type"),
					resource.TestCheckResourceAttrSet(rLine, "created_on"),
				),
			},
		},
	})
}

// testAccCheckPrivateNetworkResourceDestroy verifies the private network
// has been destroyed
func testAccCheckPrivateNetworkResourceDestroy(s *terraform.State) error {
	// get configured client from metadata
	client := testAccProvider.Meta().(receiver.BMCSDK)
	// loop through the resources in state, verifying each private network
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pnap_private_network" {
			continue
		}

		// Retrieve our private network by referencing its state ID for API lookup
		requestCommand := helperprivatenetwork.NewGetPrivateNetworkCommand(client, rs.Primary.ID)

		_, err := requestCommand.Execute()

		if err == nil {
			return fmt.Errorf("PNAP private network (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCreatePrivateNetworkResource_basic(rName string) string {
	return fmt.Sprintf(`
resource "pnap_private_network" "%s" {
	name = "%s"
    location = "PHX"
	description = "acctest"
	cidr = "10.0.1.0/31"
}`, rName, rName)
}

func testAccUpdatePrivateNetworkResource_basic(rName string) string {
	return fmt.Sprintf(`
resource "pnap_private_network" "%s" {
	name = "%s-basic"
    location = "PHX"
	description = "acctest-basic"
	cidr = "10.0.1.0/31"
}`, rName, rName)
}

func testAccCreatePrivateNetworkResource_force(rName string) string {
	return fmt.Sprintf(`
resource "pnap_private_network" "%s" {
	name = "%s"
    location = "PHX"
	description = "acctest"
	location_default = false
	force = true
}`, rName, rName)
}

func testAccUpdatePrivateNetworkResource_force(rName string) string {
	return fmt.Sprintf(`
resource "pnap_private_network" "%s" {
	name = "%s-force"
    location = "PHX"
	description = "acctest-force"
	location_default = false
	force = true
}`, rName, rName)
}

// testAccCheckPrivateNetworkExists retrieves the private network
func testAccCheckPrivateNetworkExists(resourceName string, privateNetwork *networkapiclient.PrivateNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("private network ID is not set")
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(receiver.BMCSDK)

		requestCommand := helperprivatenetwork.NewGetPrivateNetworkCommand(client, rs.Primary.ID)

		resp, err := requestCommand.Execute()

		if err != nil {
			return err
		} else {
			*privateNetwork = *resp
		}

		return nil
	}
}

// testAccCheckPrivateNetworkAttributes_basic verifies attributes are set correctly by
// Terraform
func testAccCheckPrivateNetworkAttributes_basic(resourceName string, privateNetwork *networkapiclient.PrivateNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if privateNetwork.Name != resourceName {
			return fmt.Errorf("name not set to %s name is %s", resourceName, privateNetwork.Name)
		}
		if privateNetwork.Location != "PHX" {
			return fmt.Errorf("location is not set")
		}
		if privateNetwork.Cidr != nil {
			if *privateNetwork.Cidr != "10.0.1.0/31" {
				return fmt.Errorf("cidr is not set")
			}
		} else {
			return fmt.Errorf("cidr is not set")
		}
		if privateNetwork.LocationDefault != false {
			return fmt.Errorf("location default is not set")
		}
		if privateNetwork.Description != nil {
			if *privateNetwork.Description != "acctest" {
				return fmt.Errorf("description is not set")
			}
		} else {
			return fmt.Errorf("description is not set")
		}

		return nil
	}
}

// testAccCheckPrivateNetworkNewAttributes_basic verifies attributes are updated correctly by
// Terraform
func testAccCheckPrivateNetworkNewAttributes_basic(resourceName string, privateNetwork *networkapiclient.PrivateNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if privateNetwork.Name != resourceName {
			return fmt.Errorf("name not updated to %s name is %s", resourceName, privateNetwork.Name)
		}
		if privateNetwork.Location != "PHX" {
			return fmt.Errorf("location is not set")
		}
		if privateNetwork.Cidr != nil {
			if *privateNetwork.Cidr != "10.0.1.0/31" {
				return fmt.Errorf("cidr is not set")
			}
		} else {
			return fmt.Errorf("cidr is not set")
		}
		if privateNetwork.LocationDefault != false {
			return fmt.Errorf("location default is not set")
		}
		if privateNetwork.Description != nil {
			if *privateNetwork.Description != "acctest-basic" {
				return fmt.Errorf("description is not updated")
			}
		} else {
			return fmt.Errorf("description is not updated")
		}

		return nil
	}
}

// testAccCheckPrivateNetworkAttributes_force verifies attributes are set correctly by
// Terraform
func testAccCheckPrivateNetworkAttributes_force(resourceName string, privateNetwork *networkapiclient.PrivateNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if privateNetwork.Name != resourceName {
			return fmt.Errorf("name not set to %s name is %s", resourceName, privateNetwork.Name)
		}
		if privateNetwork.Location != "PHX" {
			return fmt.Errorf("location is not set")
		}
		if privateNetwork.Cidr != nil {
			return fmt.Errorf("cidr is not set")
		}
		if privateNetwork.LocationDefault != false {
			return fmt.Errorf("location default is not set")
		}
		if privateNetwork.Description != nil {
			if *privateNetwork.Description != "acctest" {
				return fmt.Errorf("description is not set")
			}
		} else {
			return fmt.Errorf("description is not set")
		}

		return nil
	}
}

// testAccCheckPrivateNetworkNewAttributes_force verifies attributes are updated correctly by
// Terraform
func testAccCheckPrivateNetworkNewAttributes_force(resourceName string, privateNetwork *networkapiclient.PrivateNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if privateNetwork.Name != resourceName {
			return fmt.Errorf("name not updated to %s name is %s", resourceName, privateNetwork.Name)
		}
		if privateNetwork.Location != "PHX" {
			return fmt.Errorf("location is not set")
		}
		if privateNetwork.Cidr != nil {
			return fmt.Errorf("cidr is not set")
		}
		if privateNetwork.LocationDefault != false {
			return fmt.Errorf("location default is not set")
		}
		if privateNetwork.Description != nil {
			if *privateNetwork.Description != "acctest-force" {
				return fmt.Errorf("description is not updated")
			}
		} else {
			return fmt.Errorf("description is not updated")
		}

		return nil
	}
}

func init() {
	resource.AddTestSweepers("private-network", &resource.Sweeper{
		Name: "private-network",
		F: func(region string) error {

			// retrieve the configured client from the test setup
			client := testAccProvider.Meta().(receiver.BMCSDK)

			requestCommand := helperprivatenetwork.NewGetPrivateNetworksCommand(client)
			resp, err := requestCommand.Execute()

			if err != nil {
				return fmt.Errorf("Error getting private networks: %s", err)
			} else {

				for _, instance := range resp {
					if strings.HasPrefix(instance.Name, "acctest") {
						deleteCommand := helperprivatenetwork.NewDeletePrivateNetworkCommand(client, instance.Id)
						err := deleteCommand.Execute()

						if err != nil {
							return fmt.Errorf("Error destroying %s during sweep: %s ", instance.Name, err)
						}
					}
				}
			}

			return nil
		},
	})
}
