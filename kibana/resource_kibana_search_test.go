package kibana

import (
	"fmt"
	"testing"

	kibana "github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"strings"
)

const kibanaIndexVanilla = "${data.kibana_index.main.id}"
const kibanaIndexLogzio = "[logzioCustomerIndex]YYMMDD"

var testSearchCreate = map[kibana.KibanaType]string{
	kibana.KibanaTypeVanilla: fmt.Sprintf(testCreateSearchConfig, kibanaIndexVanilla, dataKibanaIndex),
	kibana.KibanaTypeLogzio:  fmt.Sprintf(testCreateSearchConfig, kibanaIndexLogzio, ""),
}

var testSearchCreateMeta = map[kibana.KibanaType]string{
	kibana.KibanaTypeVanilla: fmt.Sprintf(testCreateSearchConfigMeta, kibanaIndexVanilla, kibanaIndexVanilla, kibanaIndexVanilla, dataKibanaIndex),
	kibana.KibanaTypeLogzio:  fmt.Sprintf(testCreateSearchConfigMeta, kibanaIndexLogzio, kibanaIndexLogzio, "", ""),
}

var testSearchCreateQuery = map[kibana.KibanaType]string{
	kibana.KibanaTypeVanilla: fmt.Sprintf(testCreateSearchConfigQuery, kibanaIndexVanilla, dataKibanaIndex),
	kibana.KibanaTypeLogzio:  fmt.Sprintf(testCreateSearchConfigQuery, kibanaIndexLogzio, ""),
}

var testSearchUpdate = map[kibana.KibanaType]string{
	kibana.KibanaTypeVanilla: fmt.Sprintf(testUpdateSearchConfig, kibanaIndexVanilla, dataKibanaIndex),
	kibana.KibanaTypeLogzio:  fmt.Sprintf(testUpdateSearchConfig, kibanaIndexLogzio, ""),
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
				Config: testSearchUpdate[testConfig.KibanaType],
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

func TestAccKibanaSearchApi_WithQuery(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKibanaSearchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSearchCreateQuery[testConfig.KibanaType],
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaSearchExists("kibana_search.china"),
					resource.TestCheckResourceAttr("kibana_search.china", "name", "Chinese search with query"),
					resource.TestCheckResourceAttr("kibana_search.china", "description", "Chinese search results with query"),
					CheckResourceAttrSet("kibana_search.china", "search.#.query", "geo.src:china"),
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
  name            = "Chinese search"
  description     = "Chinese search results"
  display_columns = ["_source"]
  sort_by_columns = ["@timestamp"]

  search {
    index = "%s"

    filters {
      match {
        field_name = "geo.src"
        query      = "CN"
        type       = "phrase"
      }
    }
  }
}

%s
`

const testCreateSearchConfigMeta = `
resource "kibana_search" "china" {
  name            = "Chinese search with filter meta"
  description     = "Chinese search results with filter meta"
  display_columns = ["_source"]
  sort_by_columns = ["@timestamp"]

  search {
    index = "%s"

    filters {
      match {
        field_name = "geo.src"
        query      = "CN"
        type       = "phrase"
      }

      meta {
        index = "%s"
        alias = "China"
        type  = "phrase"
        key   = "geo.src"
        value = "CN"

        params {
          query = "CN"
          type  = "phrase"
        }
      }
    }

    filters {
      exists = "geoip.region_name"

      meta {
        index = "%s"
        type  = "exists"
        key   = "geoip.region_name"
        value = "exists"
      }
    }
  }
}

%s
`
const testCreateSearchConfigQuery = `
resource "kibana_search" "china" {
	name 	        = "Chinese search with query"
	description     = "Chinese search results with query"
	display_columns = ["_source"]
	sort_by_columns = ["@timestamp"]
	search {
		index   = "%s"
		query   = "geo.src:china"
	}
}

%s
`

const testUpdateSearchConfig = `
resource "kibana_search" "china" {
  name            = "Chinese search - errors"
  description     = "Chinese errors"
  display_columns = ["_source"]
  sort_by_columns = ["@timestamp"]

  search {
    index = "%s"

    filters {
      match {
        field_name = "geo.src"
        query      = "CN"
        type       = "phrase"
      }
    }

    filters {
      match {
        field_name = "@tags"
        query      = "error"
        type       = "phrase"
      }
    }
  }
}

%s
`
