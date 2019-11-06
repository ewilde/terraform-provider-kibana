package kibana

import (
	"fmt"
	"os"
	"strings"
	"testing"

	kibana "github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func skipIfNotXpackSecurity(t *testing.T) {
	_, useXpackSecurity := os.LookupEnv("USE_XPACK_SECURITY")
	if !useXpackSecurity {
		t.Skip("Skipping testing as we don't have xpack security")
	}
}

func TestAccKibanaRoleBasic(t *testing.T) {
	skipIfNotXpackSecurity(t)
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKibanaRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testRoleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaRoleExists("kibana_role.manager"),
					resource.TestCheckResourceAttr("kibana_role.manager", "name", "manager"),
					resource.TestCheckResourceAttr("kibana_role.manager", "elasticsearch.#", "1"),
					resource.TestCheckResourceAttr("kibana_role.manager", "elasticsearch.0.cluster.#", "2"),
					resource.TestCheckResourceAttr("kibana_role.manager", "kibana.#", "1"),
					resource.TestCheckResourceAttr("kibana_role.manager", "kibana.0.feature.#", "1"),
				),
			},
			{
				Config: testRoleConfigBasicExtraFeature,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaRoleExists("kibana_role.manager"),
					resource.TestCheckResourceAttr("kibana_role.manager", "name", "manager"),
					resource.TestCheckResourceAttr("kibana_role.manager", "elasticsearch.#", "1"),
					resource.TestCheckResourceAttr("kibana_role.manager", "elasticsearch.0.cluster.#", "2"),
					resource.TestCheckResourceAttr("kibana_role.manager", "kibana.#", "1"),
					resource.TestCheckResourceAttr("kibana_role.manager", "kibana.0.feature.#", "2"),
					resource.TestCheckResourceAttr("kibana_role.manager", "elasticsearch.0.indices.#", "1"),
				),
			},
			{
				Config: testRoleConfigBasicExtraIndices,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaRoleExists("kibana_role.manager"),
					resource.TestCheckResourceAttr("kibana_role.manager", "name", "manager"),
					resource.TestCheckResourceAttr("kibana_role.manager", "elasticsearch.#", "1"),
					resource.TestCheckResourceAttr("kibana_role.manager", "elasticsearch.0.cluster.#", "2"),
					resource.TestCheckResourceAttr("kibana_role.manager", "kibana.#", "1"),
					resource.TestCheckResourceAttr("kibana_role.manager", "kibana.0.feature.#", "2"),
					resource.TestCheckResourceAttr("kibana_role.manager", "elasticsearch.0.indices.#", "2"),
				),
			},
		},
	})
}

func testAccCheckKibanaRoleExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		api, err := testAccProvider.Meta().(*kibana.KibanaClient).Role().GetByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("role with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckKibanaRoleDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*kibana.KibanaClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "kibana_role" {
			continue
		}

		response, err := client.Role().GetByID(rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("error calling get role by id: %v", err)
		}

		if response != nil {
			return fmt.Errorf("role %s still exists, %+v", rs.Primary.ID, response)
		}
	}

	return nil

}

const testRoleConfigBasic = `
resource "kibana_role" "manager" {
  name = "manager"
  elasticsearch {
    cluster = ["a", "b"]
    indices {
      privileges = ["all"]
      names      = ["foo-*"]
    }
  }
  kibana {
    feature {
      name       = "discover"
      privileges = ["all"]
    }
	spaces     = ["default"]
  }
}
`

const testRoleConfigBasicExtraFeature = `
resource "kibana_role" "manager" {
  name = "manager"
  elasticsearch {
    cluster = ["a", "b"]
    indices {
      privileges = ["all"]
      names      = ["foo-*"]
    }
  }
  kibana {
    feature {
      name       = "discover"
      privileges = ["all"]
    }
    feature {
      name       = "dashboard"
      privileges = ["all"]
    }

	spaces     = ["default"]
  }
}
`

const testRoleConfigBasicExtraIndices = `
resource "kibana_role" "manager" {
  name = "manager"
  elasticsearch {
    cluster = ["a", "b"]
    indices {
      privileges = ["all"]
      names      = ["foo-*"]
    }
    indices {
      privileges = ["all"]
      names      = ["bar-*"]
    }
  }
  kibana {
    feature {
      name       = "discover"
      privileges = ["all"]
    }
    feature {
      name       = "dashboard"
      privileges = ["all"]
    }

	spaces     = ["default"]
  }
}
`
