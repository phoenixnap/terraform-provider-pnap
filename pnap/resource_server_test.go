package pnap

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	helperserver "github.com/PNAP/go-sdk-helper-bmc/command/bmcapi/server"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	bmcapiclient "github.com/phoenixnap/go-sdk-bmc/bmcapi"
)

func TestAccPnapServer_basic(t *testing.T) {

	var server bmcapiclient.Server
	// generate a random name for each widget test run, to avoid
	// collisions from multiple concurrent tests.
	// the acctest package includes many helpers such as RandStringFromCharSet
	// See https://godoc.org/github.com/hashicorp/terraform-plugin-sdk/helper/acctest
	rNameSuffix := acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)
	rName := "acctest-" + rNameSuffix
	rLine := "pnap_server." + rName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerResourceDestroy,
		Steps: []resource.TestStep{
			{
				// use configuration for server creation
				Config: testAccCreateServerResource(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the server object
					testAccCheckServerExists(rLine, &server),
					// verify remote values
					testAccCheckServerAttributes(rName, &server),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "hostname", rName),
					resource.TestCheckResourceAttr(rLine, "os", "ubuntu/bionic"),
					resource.TestCheckResourceAttr(rLine, "type", "s1.c1.medium"),
					resource.TestCheckResourceAttr(rLine, "location", "PHX"),
					resource.TestCheckResourceAttrSet(rLine, "location"),
					resource.TestCheckResourceAttrSet(rLine, "status"),
					resource.TestCheckResourceAttrSet(rLine, "ssh_keys.#"),
					resource.TestCheckResourceAttrSet(rLine, "private_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(rLine, "public_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(rLine, "ram"),
				),
			},
			{
				// use same configuration with power off action
				Config: testAccPowerOffServerResource(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the server object
					testAccCheckServerExists(rLine, &server),
					// verify remote values
					testAccCheckServerStatusAttribute(rName, &server),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "status", "powered-off"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				// use same configuration from above with power on action
				Config: testAccPowerOnServerResource(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the server object
					testAccCheckServerExists(rLine, &server),
					// verify remote values
					testAccCheckServerAttributes(rName, &server),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "status", "powered-on"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				// use same configuration from above with reboot action
				Config: testAccRebootServerResource(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the server object
					testAccCheckServerExists(rLine, &server),
					// verify remote values
					testAccCheckServerAttributes(rName, &server),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "status", "powered-on"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccPnapServer_shutdowntest(t *testing.T) {
	rNameSuffix := acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)
	rName := "acctest-" + rNameSuffix
	rLine := "pnap_server." + rName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerResourceDestroy,
		Steps: []resource.TestStep{
			{
				// use configuration for server creation
				Config: testAccCreateServerResource(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(

					// verify local values
					resource.TestCheckResourceAttr(rLine, "hostname", rName),
					resource.TestCheckResourceAttr(rLine, "os", "ubuntu/bionic"),
					resource.TestCheckResourceAttr(rLine, "type", "s1.c1.medium"),
					resource.TestCheckResourceAttr(rLine, "location", "PHX"),
					resource.TestCheckResourceAttrSet(rLine, "location"),
					resource.TestCheckResourceAttrSet(rLine, "status"),
					resource.TestCheckResourceAttrSet(rLine, "ssh_keys.#"),
					resource.TestCheckResourceAttrSet(rLine, "private_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(rLine, "public_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(rLine, "ram"),
				),
			},
			{
				// update previously used configuration with shutdown action
				Config: testAccShutDownServerResource(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the server object
					//testAccCheckServerExists(rName, &server),
					// verify remote values
					//testAccCheckServerStatusAttribute(&server),
					// verify local values
					resource.TestCheckResourceAttr(rLine, "status", "powered-off"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// testAccPreCheck validates the necessary test API keys exist
// in the testing environment
func testAccPreCheck(t *testing.T) {
	//client := testAccProvider.Meta().(receiver.BMCSDK)
	/* err := client..VerifyConfiguration()
	if err != nil {
		t.Fatal(err)
	} */
}

// testAccCheckServerResourceDestroy verifies the server
// has been destroyed
func testAccCheckServerResourceDestroy(s *terraform.State) error {
	// get configured client from metadata
	client := testAccProvider.Meta().(receiver.BMCSDK)
	// loop through the resources in state, verifying each server
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pnap_server" {
			continue
		}

		// Retrieve our server by referencing it's state ID for API lookup
		requestCommand := helperserver.NewGetServerCommand(client, rs.Primary.ID)

		_, err := requestCommand.Execute()
		/* code := resp.StatusCode */
		if err == nil {
			return fmt.Errorf("PNAP Server (%s) still exists", rs.Primary.ID)
		}
		/* if code != 200 && code != 404 {
			response := &dto.ErrorMessage{}
			response.FromBytes(resp)
			return fmt.Errorf("API Returned Code: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors)
		}
		if code == 200 {
			return fmt.Errorf("PNAP Server (%s) still exists", rs.Primary.ID)
		} */
	}

	return nil
}
func testAccCreateServerResource(rName string) string {
	return fmt.Sprintf(`
resource "pnap_server" "%s" {
	hostname = "%s"
    os = "ubuntu/bionic"
    type = "s1.c1.medium"
    location = "PHX"
    ssh_keys = [
        "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDF9LdAFElNCi7JoWh6KUcchrJ2Gac1aqGRPpdZNowObpRtmiRCecAMb7bUgNAaNfcmwiQi7tos9TlnFgprIcfMWb8MSs3ABYHmBgqEEt3RWYf0fAc9CsIpJdMCUG28TPGTlRXCEUVNKgLMdcseAlJoGp1CgbHWIN65fB3he3kAZcfpPn5mapV0tsl2p+ZyuAGRYdn5dJv2RZDHUZBkOeUobwsij+weHCKAFmKQKtCP7ybgVHaQjAPrj8MGnk1jBbjDt5ws+Be+9JNjQJee9zCKbAOsIo3i+GcUIkrw5jxPU/RTGlWBcemPaKHdciSzGcjWboapzIy49qypQhZe1U75 user2@172.16.1.106"
    
    ]
    #allowed actions are: reboot, reset, powered-on, powered-off, shutdown
    #action = "powered-on"
}`, rName, rName)
}

func testAccRebootServerResource(rName string) string {
	return fmt.Sprintf(`
resource "pnap_server" "%s" {
	hostname = "%s"
    os = "ubuntu/bionic"
    type = "s1.c1.medium"
    location = "PHX"
    ssh_keys = [
        "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDF9LdAFElNCi7JoWh6KUcchrJ2Gac1aqGRPpdZNowObpRtmiRCecAMb7bUgNAaNfcmwiQi7tos9TlnFgprIcfMWb8MSs3ABYHmBgqEEt3RWYf0fAc9CsIpJdMCUG28TPGTlRXCEUVNKgLMdcseAlJoGp1CgbHWIN65fB3he3kAZcfpPn5mapV0tsl2p+ZyuAGRYdn5dJv2RZDHUZBkOeUobwsij+weHCKAFmKQKtCP7ybgVHaQjAPrj8MGnk1jBbjDt5ws+Be+9JNjQJee9zCKbAOsIo3i+GcUIkrw5jxPU/RTGlWBcemPaKHdciSzGcjWboapzIy49qypQhZe1U75 user2@172.16.1.106"
    
    ]
    #allowed actions are: reboot, reset, powered-on, powered-off, shutdown
    action = "reboot"
}`, rName, rName)
}
func testAccResetServerResource(rName string) string {
	return fmt.Sprintf(`
resource "pnap_server" "%s" {
	hostname = "%s"
    os = "ubuntu/bionic"
    type = "s1.c1.medium"
    location = "PHX"
    ssh_keys = [
        "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDF9LdAFElNCi7JoWh6KUcchrJ2Gac1aqGRPpdZNowObpRtmiRCecAMb7bUgNAaNfcmwiQi7tos9TlnFgprIcfMWb8MSs3ABYHmBgqEEt3RWYf0fAc9CsIpJdMCUG28TPGTlRXCEUVNKgLMdcseAlJoGp1CgbHWIN65fB3he3kAZcfpPn5mapV0tsl2p+ZyuAGRYdn5dJv2RZDHUZBkOeUobwsij+weHCKAFmKQKtCP7ybgVHaQjAPrj8MGnk1jBbjDt5ws+Be+9JNjQJee9zCKbAOsIo3i+GcUIkrw5jxPU/RTGlWBcemPaKHdciSzGcjWboapzIy49qypQhZe1U75 user2@172.16.1.106"
    
    ]
    #allowed actions are: reboot, reset, powered-on, powered-off, shutdown
    action = "reset"
}`, rName, rName)
}
func testAccPowerOnServerResource(rName string) string {
	return fmt.Sprintf(`
resource "pnap_server" "%s" {
	hostname = "%s"
    os = "ubuntu/bionic"
    type = "s1.c1.medium"
    location = "PHX"
    ssh_keys = [
        "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDF9LdAFElNCi7JoWh6KUcchrJ2Gac1aqGRPpdZNowObpRtmiRCecAMb7bUgNAaNfcmwiQi7tos9TlnFgprIcfMWb8MSs3ABYHmBgqEEt3RWYf0fAc9CsIpJdMCUG28TPGTlRXCEUVNKgLMdcseAlJoGp1CgbHWIN65fB3he3kAZcfpPn5mapV0tsl2p+ZyuAGRYdn5dJv2RZDHUZBkOeUobwsij+weHCKAFmKQKtCP7ybgVHaQjAPrj8MGnk1jBbjDt5ws+Be+9JNjQJee9zCKbAOsIo3i+GcUIkrw5jxPU/RTGlWBcemPaKHdciSzGcjWboapzIy49qypQhZe1U75 user2@172.16.1.106"
    
    ]
    #allowed actions are: reboot, reset, powered-on, powered-off, shutdown
    action = "powered-on"
}`, rName, rName)
}
func testAccPowerOffServerResource(rName string) string {
	return fmt.Sprintf(`
resource "pnap_server" "%s" {
	hostname = "%s"
    os = "ubuntu/bionic"
    type = "s1.c1.medium"
    location = "PHX"
    ssh_keys = [
        "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDF9LdAFElNCi7JoWh6KUcchrJ2Gac1aqGRPpdZNowObpRtmiRCecAMb7bUgNAaNfcmwiQi7tos9TlnFgprIcfMWb8MSs3ABYHmBgqEEt3RWYf0fAc9CsIpJdMCUG28TPGTlRXCEUVNKgLMdcseAlJoGp1CgbHWIN65fB3he3kAZcfpPn5mapV0tsl2p+ZyuAGRYdn5dJv2RZDHUZBkOeUobwsij+weHCKAFmKQKtCP7ybgVHaQjAPrj8MGnk1jBbjDt5ws+Be+9JNjQJee9zCKbAOsIo3i+GcUIkrw5jxPU/RTGlWBcemPaKHdciSzGcjWboapzIy49qypQhZe1U75 user2@172.16.1.106"
    
    ]
    #allowed actions are: reboot, reset, powered-on, powered-off, shutdown
    action = "powered-off"
}`, rName, rName)
}
func testAccShutDownServerResource(rName string) string {
	return fmt.Sprintf(`
resource "pnap_server" "%s" {
	hostname = "%s"
    os = "ubuntu/bionic"
    type = "s1.c1.medium"
    location = "PHX"
    ssh_keys = [
        "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDF9LdAFElNCi7JoWh6KUcchrJ2Gac1aqGRPpdZNowObpRtmiRCecAMb7bUgNAaNfcmwiQi7tos9TlnFgprIcfMWb8MSs3ABYHmBgqEEt3RWYf0fAc9CsIpJdMCUG28TPGTlRXCEUVNKgLMdcseAlJoGp1CgbHWIN65fB3he3kAZcfpPn5mapV0tsl2p+ZyuAGRYdn5dJv2RZDHUZBkOeUobwsij+weHCKAFmKQKtCP7ybgVHaQjAPrj8MGnk1jBbjDt5ws+Be+9JNjQJee9zCKbAOsIo3i+GcUIkrw5jxPU/RTGlWBcemPaKHdciSzGcjWboapzIy49qypQhZe1U75 user2@172.16.1.106"
    
    ]
    #allowed actions are: reboot, reset, powered-on, powered-off, shutdown
    action = "shutdown"
}`, rName, rName)
}

// testAccCheckServerExists uses the SDK directly to retrieve
// the server, and stores it in the provided
// *dto.LongServer
func testAccCheckServerExists(resourceName string, server *bmcapiclient.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Server ID is not set")
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(receiver.BMCSDK)

		requestCommand := helperserver.NewGetServerCommand(client, rs.Primary.ID)

		resp, err := requestCommand.Execute()

		if err != nil {
			return err
		} else {
			*server = *resp
		}
		/* code := resp.StatusCode
		if code != 200 && code != 404 {
			response := &dto.ErrorMessage{}
			response.FromBytes(resp)
			return fmt.Errorf("API Returned Code: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors)
		}
		if code == 404 {
			return fmt.Errorf("PNAP Server (%s) not found ", rs.Primary.ID)
		}
		if code == 200 {
			resultServer := &dto.LongServer{}
			resultServer.FromBytes(resp)
			*server = *resultServer
		} */
		return nil
	}
}

// testAccCheckServerAttributes verifies attributes are set correctly by
// Terraform
func testAccCheckServerAttributes(resourceName string, server *bmcapiclient.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if server.Hostname != resourceName {
			return fmt.Errorf("hostname not set to %s name is %s", resourceName, server.Hostname)
		}

		if server.Os != "ubuntu/bionic" {
			return fmt.Errorf("OS is not set")
		}
		if server.Type != "s1.c1.medium" {
			return fmt.Errorf("type is not set")
		}
		if server.Location != "PHX" {
			return fmt.Errorf("location is not set")
		}
		if server.Status != "powered-on" {
			return fmt.Errorf("status is not set, should be powered-on")
		}
		if len(server.PrivateIpAddresses) < 1 {
			return fmt.Errorf("private ip is not set")
		}
		if len(server.PublicIpAddresses) < 1 {
			return fmt.Errorf("public ip is not set")
		}

		return nil
	}
}

// testAccCheckServerStatusAttribute verifies status attribute is set correctly by
// Terraform
func testAccCheckServerStatusAttribute(resourceName string, server *bmcapiclient.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if server.Hostname != resourceName {
			return fmt.Errorf("hostname not set to %s name is %s", resourceName, server.Hostname)
		}
		if server.Status != "powered-off" {
			return fmt.Errorf("status is not set, should be powered-off")
		}
		return nil
	}
}
func init() {
	resource.AddTestSweepers("server", &resource.Sweeper{
		Name: "server",
		F: func(region string) error {

			// retrieve the configured client from the test setup
			client := testAccProvider.Meta().(receiver.BMCSDK)

			requestCommand := helperserver.NewGetServersCommand(client)
			resp, err := requestCommand.Execute()
			//response := &dto.Servers{}
			if err != nil {
				return fmt.Errorf("Error getting servers: %s", err)
			} else {

				//response.FromBytes(resp)

				for _, instance := range resp {
					if strings.HasPrefix(instance.Hostname, "acctest") {
						deleteCommand := helperserver.NewDeleteServerCommand(client, instance.Id)
						_, err := deleteCommand.Execute()

						if err != nil {
							return fmt.Errorf("Error destroying %s during sweep: %s ", instance.Hostname, err)
						}
						/* code := resp.StatusCode
						if code != 200 && code != 404 {
							delresponse := &dto.ErrorMessage{}
							delresponse.FromBytes(resp)
							return fmt.Errorf("API Returned Code: %v, Message: %v, Validation Errors: %v", code, delresponse.Message, delresponse.ValidationErrors)
						} */

					}
				}
			}

			return nil
		},
	})
}
