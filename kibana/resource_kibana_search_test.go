package kibana

import (
	"fmt"
	"github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"

	"strings"
)

var testSearchCreate = map[kibana.KibanaType]string{
	kibana.KibanaTypeVanilla: testCreateSearchConfig,
	kibana.KibanaTypeLogzio:  testCreateSearchLogzioConfig,
}

var testSearchCreateMeta = map[kibana.KibanaType]string{
	kibana.KibanaTypeVanilla: testCreateSearchConfigMeta,
	kibana.KibanaTypeLogzio:  testCreateSearchLogzioConfigMeta,
}

var testSearchUpdate = map[kibana.KibanaType]string{
	kibana.KibanaTypeVanilla: testUpdateSearchConfig,
	kibana.KibanaTypeLogzio:  testUpdateSearchLogzioConfig,
}

func TestAccKibanaSearchApi(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKibanaSearchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSearchCreate[testConfig.KibanaType],
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaSearchExists("kibana_search.china"),
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
				Config: testSearchUpdate[testAccProvider.Meta().(*kibana.KibanaClient).Config.KibanaType],
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaSearchExists("kibana_search.china"),
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

func TestAccKibanaSearchApi_WithMetaFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKibanaSearchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSearchCreateMeta[testConfig.KibanaType],
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaSearchExists("kibana_search.china"),
					resource.TestCheckResourceAttr("kibana_search.china", "name", "Chinese search with filter meta"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.0.meta.#.negate", "false"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.0.meta.#.disabled", "false"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.0.meta.#.alias", "China"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.0.meta.#.type", "phrase"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.0.meta.#.key", "geo.src"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.0.meta.#.params.#.query", "CN"),
					CheckResourceAttrSet("kibana_search.china", "search.#.filters.0.meta.#.params.#.type", "phrase"),
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

		if err != nil && !strings.Contains(err.Error(), "404") {
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
	description     = "Chinese search results"
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

const testCreateSearchConfigMeta = `
resource "kibana_search" "china" {
	name 	        = "Chinese search with filter meta"
	description     = "Chinese search results with filter meta"
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

				meta = {
					index = "${data.kibana_index.main.id}"
					alias = "China"
					type  = "phrase"
                    key   = "geo.src"
					value = "CN"
 					params = {
						query = "CN"
						type  = "phrase"
					}
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

const testCreateSearchLogzioConfigMeta = `
resource "kibana_search" "china" {
	name 	        = "Chinese search with filter meta"
	description     = "Chinese search results with filter meta"
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

				meta = {
					index = "[logzioCustomerIndex]YYMMDD"
					alias = "China"
					type  = "phrase"
                    key   = "geo.src"
					value = "CN"
 					params = {
						query = "CN"
						type  = "phrase"
					}
				}
			}
		]
	}
}
`
const testUpdateSearchConfig = `
resource "kibana_search" "china" {
	name 	        = "Chinese search - errors"
	description     = "Chinese errors"
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
