package pnap

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	helperipblock "github.com/PNAP/go-sdk-helper-bmc/command/ipapi/ipblock"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	ipapiclient "github.com/phoenixnap/go-sdk-bmc/ipapi/v2"
)

func TestAccPnapIpBlock_basic(t *testing.T) {

	var ipBlock ipapiclient.IpBlock
	// generate a random name for each widget test run, to avoid
	// collisions from multiple concurrent tests.
	// the acctest package includes many helpers such as RandStringFromCharSet
	// See https://godoc.org/github.com/hashicorp/terraform-plugin-sdk/helper/acctest
	rNameSuffix := acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)
	rName := "acctest-" + rNameSuffix
	rLine := "pnap_ip_block." + rName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIpBlockResourceDestroy,
		Steps: []resource.TestStep{
			{
				// use configuration for ip block creation
				Config: testAccCreateIpBlockResource(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the ip block object
					testAccCheckIpBlockExists(rLine, &ipBlock),
					// verify remote values
					testAccCheckIpBlockAttributes(rName, &ipBlock),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "status", "unassigned"),
					resource.TestCheckResourceAttr(rLine, "is_bring_your_own", "false"),
					resource.TestCheckResourceAttrSet(rLine, "cidr"),
					resource.TestCheckResourceAttrSet(rLine, "created_on"),
				),
			},
			{
				// update ip block details
				Config: testAccUpdateIpBlockResource(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the ip block object
					testAccCheckIpBlockExists(rLine, &ipBlock),
					// verify remote values
					testAccCheckIpBlockNewAttributes(rName, &ipBlock),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "status", "unassigned"),
					resource.TestCheckResourceAttr(rLine, "is_bring_your_own", "false"),
					resource.TestCheckResourceAttrSet(rLine, "cidr"),
					resource.TestCheckResourceAttrSet(rLine, "created_on"),
				),
			},
		},
	})
}

// testAccCheckIpBlockResourceDestroy verifies the ip block
// has been destroyed
func testAccCheckIpBlockResourceDestroy(s *terraform.State) error {
	// get configured client from metadata
	client := testAccProvider.Meta().(receiver.BMCSDK)
	// loop through the resources in state, verifying each ip block
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pnap_ip_block" {
			continue
		}

		// Retrieve our ip block by referencing its state ID for API lookup
		requestCommand := helperipblock.NewGetIpBlockCommand(client, rs.Primary.ID)

		_, err := requestCommand.Execute()

		if err == nil {
			return fmt.Errorf("PNAP Ip Block (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCreateIpBlockResource(rName string) string {
	return fmt.Sprintf(`
resource "pnap_ip_block" "%s" {
    location = "PHX"
    cidr_block_size = "/31"
    description = "acctest"
}`, rName)
}

func testAccUpdateIpBlockResource(rName string) string {
	return fmt.Sprintf(`
resource "pnap_ip_block" "%s" {
    location = "PHX"
    cidr_block_size = "/31"
    description = "acctest-basic"
}`, rName)
}

// testAccCheckIpBlockExists retrieves the ip block
func testAccCheckIpBlockExists(resourceName string, ipBlock *ipapiclient.IpBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Ip Block ID is not set")
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(receiver.BMCSDK)

		requestCommand := helperipblock.NewGetIpBlockCommand(client, rs.Primary.ID)

		resp, err := requestCommand.Execute()

		if err != nil {
			return err
		} else {
			*ipBlock = *resp
		}

		return nil
	}
}

// testAccCheckIpBlockAttributes verifies attributes are set correctly by
// Terraform
func testAccCheckIpBlockAttributes(resourceName string, ipBlock *ipapiclient.IpBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ipBlock.Location != "PHX" {
			return fmt.Errorf("location is not set")
		}
		if ipBlock.CidrBlockSize != "/31" {
			return fmt.Errorf("cidr block size is not set")
		}
		if ipBlock.Description != nil {
			if *ipBlock.Description != "acctest" {
				return fmt.Errorf("description is not set")
			}
		} else {
			return fmt.Errorf("description is not set")
		}
		if ipBlock.Status != "unassigned" {
			return fmt.Errorf("status is not set")
		}
		if ipBlock.IsBringYourOwn != false {
			return fmt.Errorf("is bring your own is not set")
		}

		return nil
	}
}

// testAccCheckIpBlockNewAttributes verifies attributes are updated correctly by
// Terraform
func testAccCheckIpBlockNewAttributes(resourceName string, ipBlock *ipapiclient.IpBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ipBlock.Location != "PHX" {
			return fmt.Errorf("location is not set")
		}
		if ipBlock.CidrBlockSize != "/31" {
			return fmt.Errorf("cidr block size is not set")
		}
		if ipBlock.Description != nil {
			if *ipBlock.Description != "acctest-basic" {
				return fmt.Errorf("description is not updated")
			}
		} else {
			return fmt.Errorf("description is not updated")
		}
		if ipBlock.Status != "unassigned" {
			return fmt.Errorf("status is not set")
		}
		if ipBlock.IsBringYourOwn != false {
			return fmt.Errorf("is bring your own is not set")
		}

		return nil
	}
}

func init() {
	resource.AddTestSweepers("ip-block", &resource.Sweeper{
		Name: "ip-block",
		F: func(region string) error {

			// retrieve the configured client from the test setup
			client := testAccProvider.Meta().(receiver.BMCSDK)

			requestCommand := helperipblock.NewGetIpBlocksCommand(client)
			resp, err := requestCommand.Execute()

			if err != nil {
				return fmt.Errorf("Error getting ip blocks: %s", err)
			} else {

				for _, instance := range resp {
					if instance.Description != nil && strings.HasPrefix(*instance.Description, "acctest") {
						deleteCommand := helperipblock.NewDeleteIpBlockCommand(client, instance.Id)
						_, err := deleteCommand.Execute()

						if err != nil {
							return fmt.Errorf("Error destroying ip block %s during sweep: %s ", instance.Id, err)
						}
					}
				}
			}

			return nil
		},
	})
}
