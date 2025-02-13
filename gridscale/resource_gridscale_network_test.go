package gridscale

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gridscale/gsclient-go/v3"
)

func TestAccResourceGridscaleNetwork_Basic(t *testing.T) {
	var object gsclient.Network
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleNetworkDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleNetworkConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleNetworkExists("gridscale_network.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "name", name),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "dhcp_active", "true"),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "dhcp_gateway", "192.168.121.1"),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "dhcp_dns", "192.168.121.2"),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "dhcp_reserved_subnet.#", "1"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleNetworkConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleNetworkExists("gridscale_network.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "l2security", "true"),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "dhcp_active", "true"),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "dhcp_gateway", "192.168.122.1"),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "dhcp_dns", "192.168.122.2"),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "dhcp_reserved_subnet.#", "1"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleNetworkExists(n string, object *gsclient.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No object UUID is set")
		}

		client := testAccProvider.Meta().(*gsclient.Client)

		id := rs.Primary.ID

		foundObject, err := client.GetNetwork(context.Background(), id)

		if err != nil {
			return err
		}

		if foundObject.Properties.ObjectUUID != id {
			return fmt.Errorf("Object not found")
		}

		*object = foundObject

		return nil
	}
}

func testAccCheckGridscaleNetworkDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_network" {
			continue
		}

		_, err := client.GetNetwork(context.Background(), rs.Primary.ID)
		if err != nil {
			if requestError, ok := err.(gsclient.RequestError); ok {
				if requestError.StatusCode != 404 {
					return fmt.Errorf("Object %s still exists", rs.Primary.ID)
				}
			} else {
				return fmt.Errorf("Unable to fetching object %s", rs.Primary.ID)
			}
		} else {
			return fmt.Errorf("Object %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckResourceGridscaleNetworkConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_network" "foo" {
  name   = "%s"
  dhcp_active = true
  dhcp_gateway = "192.168.121.1"
  dhcp_dns = "192.168.121.2"
  dhcp_range = "192.168.121.0/27"
  dhcp_reserved_subnet = ["192.168.121.0/31"]
}
`, name)
}

func testAccCheckResourceGridscaleNetworkConfig_basic_update() string {
	return fmt.Sprint(`
resource "gridscale_network" "foo" {
  name   = "newname"
  l2security = true
  dhcp_active = true
  dhcp_gateway = "192.168.122.1"
  dhcp_dns = "192.168.122.2"
  dhcp_range = "192.168.122.0/27"
  dhcp_reserved_subnet = ["192.168.122.0/31"]
}
`)
}
