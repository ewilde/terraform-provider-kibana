package kibana

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ewilde/go-kibana"

)

func TestAccKibanaSearchApi(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKibanaSearchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateSearchConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaSearchExists("kibana_search.china"),
					resource.TestCheckResourceAttr("kibana_search.china", "name", "Chinese search"),
				),
			},
			{
				Config: testUpdateSearchConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaSearchExists("kibana_search.china"),
					resource.TestCheckResourceAttr("kibana_search.china", "name", "Chinese search - errors"),
				),
			},
		},
	})
}

func testAccCheckKibanaSearchDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*kibana.KibanaClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "kibana_search" {
			continue
		}

		response, err := client.Search().GetById(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error calling get search by id: %v", err)
		}

		if response != nil {
			return fmt.Errorf("search %s still exists, %+v", rs.Primary.ID, response)
		}
	}

	return nil
}

func testAccCheckKibanaSearchExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		api, err := testAccProvider.Meta().(*kibana.KibanaClient).Search().GetById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("search with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateSearchConfig = `
resource "kibana_search" "china" {
	name 	        = "Chinese search"
	display_columns = ["_source"]
	sort_by_columns = ["@timestamp"]
	search = {
		index   = "${data.kibana_index.main.id}"
		filters = [
			{
				match = {
					field_name = "geo.src"
					query      = "CN"
					type       = "phrase"
				}
			}
		]
	}
}

data "kibana_index" "main" {
	filter = {
		name = "title"
		values = ["logstash-*"]
	}
}
`
const testUpdateSearchConfig = `
resource "kibana_search" "china" {
	name 	        = "Chinese search - errors"
	display_columns = ["_source"]
	sort_by_columns = ["@timestamp"]
	search = {
		index   = "${data.kibana_index.main.id}"
		filters = [
			{
				match = {
					field_name = "geo.src"
					query      = "CN"
					type       = "phrase"
				},
			},
			{
				match = {
					field_name = "@tags"
					query      = "error"
					type       = "phrase"
				}
			}
		]
	}
}

data "kibana_index" "main" {
	filter = {
		name = "title"
		values = ["logstash-*"]
	}
}
`
