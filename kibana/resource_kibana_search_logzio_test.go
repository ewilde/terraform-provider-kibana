package kibana

import (
	"fmt"
	"github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"

	"strings"
)

func TestAccKibanaSearchApiLogzio(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKibanaSearchLogzioDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateSearchLogzioConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaSearchLogzioExists("kibana_search.china"),
					resource.TestCheckResourceAttr("kibana_search.china", "name", "Chinese search"),
					resource.TestCheckResourceAttr("kibana_search.china", "description", "Chinese search results"),
					resource.TestCheckResourceAttr("kibana_search.china", "display_columns.0", "_source"),
					resource.TestCheckResourceAttr("kibana_search.china", "sort_by_columns.0", "@timestamp"),
					resource.TestCheckResourceAttr("kibana_search.china", "sort_ascending", "false"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.0.match.#.field_name", "geo.src"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.0.match.#.query", "CN"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.0.match.#.type", "phrase"),
				),
			},
			{
				Config: testUpdateSearchLogzioConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaSearchLogzioExists("kibana_search.china"),
					resource.TestCheckResourceAttr("kibana_search.china", "name", "Chinese search - errors"),
					resource.TestCheckResourceAttr("kibana_search.china", "description", "Chinese errors"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.1.match.#.field_name", "@tags"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.1.match.#.query", "error"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.1.match.#.type", "phrase"),
				),
			},
		},
	})
}

func testAccCheckKibanaSearchLogzioDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*kibana.KibanaClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "kibana_search" {
			continue
		}

		response, err := client.Search().GetById(rs.Primary.ID)

		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "400")) {
			return fmt.Errorf("error calling get search by id: %v", err)
		}

		if response != nil {
			return fmt.Errorf("search %s still exists, %+v", rs.Primary.ID, response)
		}
	}

	return nil
}

func testAccCheckKibanaSearchLogzioExists(resourceKey string) resource.TestCheckFunc {

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

const testCreateSearchLogzioConfig = `
resource "kibana_search" "china" {
	name 	        = "Chinese search"
	description     = "Chinese search results"
	display_columns = ["_source"]
	sort_by_columns = ["@timestamp"]
	search = {
		index   = "[logzioCustomerIndex]YYMMDD"
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
`
const testUpdateSearchLogzioConfig = `
resource "kibana_search" "china" {
	name 	        = "Chinese search - errors"
	description     = "Chinese errors"
	display_columns = ["_source"]
	sort_by_columns = ["@timestamp"]
	search = {
		index   = "[logzioCustomerIndex]YYMMDD"
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
`
