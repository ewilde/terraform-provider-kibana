package kibana

import (
	"fmt"
	"strings"
	"testing"

	kibana "github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccKibanaSpaceBasic(t *testing.T) {
	skipIfNotXpackSecurity(t)
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKibanaSpaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSpaceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaSpaceExists("kibana_space.blue"),
					resource.TestCheckResourceAttr("kibana_space.blue", "name", "blue"),
					resource.TestCheckResourceAttr("kibana_space.blue", "title", "Blue space"),
				),
			},
			{
				Config: testSpaceConfigBasicDisabledFeatures,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaSpaceExists("kibana_space.blue"),
					resource.TestCheckResourceAttr("kibana_space.blue", "name", "blue"),
					resource.TestCheckResourceAttr("kibana_space.blue", "title", "Blue space"),
					resource.TestCheckResourceAttr("kibana_space.blue", "disabled_features.#", "1"),
				),
			},
		},
	})
}

func testAccCheckKibanaSpaceExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		api, err := testAccProvider.Meta().(*kibana.KibanaClient).Space().GetByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("space with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckKibanaSpaceDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*kibana.KibanaClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "kibana_space" {
			continue
		}

		response, err := client.Space().GetByID(rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("error calling get space by id: %v", err)
		}

		if response != nil {
			return fmt.Errorf("space %s still exists, %+v", rs.Primary.ID, response)
		}
	}

	return nil

}

const testSpaceConfigBasic = `
resource "kibana_space" "blue" {
  name = "blue"
  title = "Blue space"
}
`

const testSpaceConfigBasicDisabledFeatures = `
resource "kibana_space" "blue" {
  name = "blue"
  title = "Blue space"
  disabled_features = ["dev_tools"]
}
`
