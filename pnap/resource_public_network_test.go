package pnap

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	helperpublicnetwork "github.com/PNAP/go-sdk-helper-bmc/command/networkapi/publicnetwork"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	networkapiclient "github.com/phoenixnap/go-sdk-bmc/networkapi/v4"
)

func TestAccPnapPublicNetwork_basic(t *testing.T) {

	var publicNetwork networkapiclient.PublicNetwork
	// generate a random name for each widget test run, to avoid
	// collisions from multiple concurrent tests.
	// the acctest package includes many helpers such as RandStringFromCharSet
	// See https://godoc.org/github.com/hashicorp/terraform-plugin-sdk/helper/acctest
	rNameSuffix := acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)
	rName := "acctest-" + rNameSuffix
	rName2 := "acctest-" + rNameSuffix + "-basic"
	rLine := "pnap_public_network." + rName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPublicNetworkResourceDestroy,
		Steps: []resource.TestStep{
			{
				// use configuration for public network creation
				Config: testAccCreatePublicNetworkResource(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the public network object
					testAccCheckPublicNetworkExists(rLine, &publicNetwork),
					// verify remote values
					testAccCheckPublicNetworkAttributes(rName, &publicNetwork),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "status", "READY"),
					resource.TestCheckResourceAttrSet(rLine, "vlan_id"),
					resource.TestCheckResourceAttrSet(rLine, "created_on"),
				),
			},
			{
				// update public network's details
				Config: testAccUpdatePublicNetworkResource(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the public network object
					testAccCheckPublicNetworkExists(rLine, &publicNetwork),
					// verify remote values
					testAccCheckPublicNetworkNewAttributes(rName2, &publicNetwork),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "status", "READY"),
					resource.TestCheckResourceAttrSet(rLine, "vlan_id"),
					resource.TestCheckResourceAttrSet(rLine, "created_on"),
				),
			},
		},
	})
}

// testAccCheckPublicNetworkResourceDestroy verifies the public network
// has been destroyed
func testAccCheckPublicNetworkResourceDestroy(s *terraform.State) error {
	// get configured client from metadata
	client := testAccProvider.Meta().(receiver.BMCSDK)
	// loop through the resources in state, verifying each public network
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pnap_public_network" {
			continue
		}

		// Retrieve our public network by referencing its state ID for API lookup
		requestCommand := helperpublicnetwork.NewGetPublicNetworkCommand(client, rs.Primary.ID)

		_, err := requestCommand.Execute()

		if err == nil {
			return fmt.Errorf("PNAP public network (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCreatePublicNetworkResource(rName string) string {
	return fmt.Sprintf(`
resource "pnap_public_network" "%s" {
	name = "%s"
	location = "PHX"
	description = "acctest"
}`, rName, rName)
}

func testAccUpdatePublicNetworkResource(rName string) string {
	return fmt.Sprintf(`
resource "pnap_public_network" "%s" {
	name = "%s-basic"
	location = "PHX"
	description = "acctest-basic"
}`, rName, rName)
}

// testAccCheckPublicNetworkExists retrieves the public network
func testAccCheckPublicNetworkExists(resourceName string, publicNetwork *networkapiclient.PublicNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("public network ID is not set")
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(receiver.BMCSDK)

		requestCommand := helperpublicnetwork.NewGetPublicNetworkCommand(client, rs.Primary.ID)

		resp, err := requestCommand.Execute()

		if err != nil {
			return err
		} else {
			*publicNetwork = *resp
		}

		return nil
	}
}

// testAccCheckPublicNetworkAttributes verifies attributes are set correctly by
// Terraform
func testAccCheckPublicNetworkAttributes(resourceName string, publicNetwork *networkapiclient.PublicNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if publicNetwork.Name != resourceName {
			return fmt.Errorf("name not set to %s name is %s", resourceName, publicNetwork.Name)
		}
		if publicNetwork.Location != "PHX" {
			return fmt.Errorf("location is not set")
		}
		if publicNetwork.Description != nil {
			if *publicNetwork.Description != "acctest" {
				return fmt.Errorf("description is not set")
			}
		} else {
			return fmt.Errorf("description is not set")
		}

		return nil
	}
}

// testAccCheckPublicNetworkNewAttributes verifies attributes are updated correctly by
// Terraform
func testAccCheckPublicNetworkNewAttributes(resourceName string, publicNetwork *networkapiclient.PublicNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if publicNetwork.Name != resourceName {
			return fmt.Errorf("name not updated to %s name is %s", resourceName, publicNetwork.Name)
		}
		if publicNetwork.Location != "PHX" {
			return fmt.Errorf("location is not set")
		}
		if publicNetwork.Description != nil {
			if *publicNetwork.Description != "acctest-basic" {
				return fmt.Errorf("description is not updated")
			}
		} else {
			return fmt.Errorf("description is not updated")
		}

		return nil
	}
}

func init() {
	resource.AddTestSweepers("public-network", &resource.Sweeper{
		Name: "public-network",
		F: func(region string) error {

			// retrieve the configured client from the test setup
			client := testAccProvider.Meta().(receiver.BMCSDK)

			requestCommand := helperpublicnetwork.NewGetPublicNetworksCommand(client)
			resp, err := requestCommand.Execute()

			if err != nil {
				return fmt.Errorf("Error getting public networks: %s", err)
			} else {

				for _, instance := range resp {
					if strings.HasPrefix(instance.Name, "acctest") {
						deleteCommand := helperpublicnetwork.NewDeletePublicNetworkCommand(client, instance.Id)
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
