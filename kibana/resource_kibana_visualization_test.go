package kibana

import (
	"fmt"
	"testing"

	kibana "github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"strings"
)

var testVisualizationCreate = map[kibana.KibanaType]string{
	kibana.KibanaTypeVanilla: testCreateVisualizationConfig,
	kibana.KibanaTypeLogzio:  testCreateVisualizationLogzioConfig,
}

var testVisualizationUpdate = map[kibana.KibanaType]string{
	kibana.KibanaTypeVanilla: testUpdateVisualizationConfig,
	kibana.KibanaTypeLogzio:  testUpdateVisualizationLogzioConfig,
}

func TestAccKibanaVisualizationApi(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKibanaVisualizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testVisualizationCreate[testConfig.KibanaType],
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaVisualizationExists("kibana_visualization.china_viz"),
					resource.TestCheckResourceAttr("kibana_visualization.china_viz", "name", "Chinese visualization"),
					resource.TestCheckResourceAttr("kibana_visualization.china_viz", "description", "Chinese error visualization"),
				),
			},
			{
				Config: testVisualizationUpdate[testConfig.KibanaType],
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaVisualizationExists("kibana_visualization.china_viz"),
					resource.TestCheckResourceAttr("kibana_visualization.china_viz", "name", "Chinese visualization - updated"),
					resource.TestCheckResourceAttr("kibana_visualization.china_viz", "description", "Chinese error visualization - updated"),
				),
			},
		},
	})
}

func testAccCheckKibanaVisualizationDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*kibana.KibanaClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "kibana_visualization" {
			continue
		}

		response, err := client.Visualization().GetById(rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("error calling get visualization by id: %v", err)
		}

		if response != nil {
			return fmt.Errorf("visualization %s still exists, %+v", rs.Primary.ID, response)
		}
	}

	return nil
}

func testAccCheckKibanaVisualizationExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		api, err := testAccProvider.Meta().(*kibana.KibanaClient).Visualization().GetById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("visualization with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateVisualizationConfig = `
resource "kibana_visualization" "china_viz" {
	name 	            = "Chinese visualization"
	description         = "Chinese error visualization"
	saved_search_id     = "${kibana_search.china.id}"
	visualization_state = <<EOF
{
  "title": "Chinese search",
  "type": "gauge",
  "params": {
    "type": "gauge",
    "addTooltip": true,
    "addLegend": true,
    "gauge": {
      "verticalSplit": false,
      "extendRange": true,
      "percentageMode": false,
      "gaugeType": "Arc",
      "gaugeStyle": "Full",
      "backStyle": "Full",
      "orientation": "vertical",
      "colorSchema": "Green to Red",
      "gaugeColorMode": "Labels",
      "colorsRange": [
        {
          "from": 0,
          "to": 50
        },
        {
          "from": 50,
          "to": 75
        },
        {
          "from": 75,
          "to": 100
        }
      ],
      "invertColors": false,
      "labels": {
        "show": true,
        "color": "black"
      },
      "scale": {
        "show": true,
        "labels": false,
        "color": "#333"
      },
      "type": "meter",
      "style": {
        "bgWidth": 0.9,
        "width": 0.9,
        "mask": false,
        "bgMask": false,
        "maskBars": 50,
        "bgFill": "#eee",
        "bgColor": false,
        "subText": "",
        "fontSize": 60,
        "labelColor": true
      }
    }
  },
  "aggs": [
    {
      "id": "1",
      "enabled": true,
      "type": "count",
      "schema": "metric",
      "params": {}
    }
  ]
}
EOF
}

resource "kibana_search" "china" {
	name 	        = "Chinese search"
	description     = "Chinese search results"
	display_columns = ["_source"]
	sort_by_columns = ["@timestamp"]
	search {
		index   = "${data.kibana_index.main.id}"
		filters {
				match {
					field_name = "geo.src"
					query      = "CN"
					type       = "phrase"
				}
		}
	}
}

data "kibana_index" "main" {
	filter {
		name = "title"
		values = ["logstash-*"]
	}
}
`
const testUpdateVisualizationConfig = `
resource "kibana_visualization" "china_viz" {
	name 	            = "Chinese visualization - updated"
	description         = "Chinese error visualization - updated"
	saved_search_id     = "${kibana_search.china.id}"
	visualization_state = <<EOF
{
  "title": "Chinese search",
  "type": "gauge",
  "params": {
    "type": "gauge",
    "addTooltip": true,
    "addLegend": true,
    "gauge": {
      "verticalSplit": false,
      "extendRange": true,
      "percentageMode": false,
      "gaugeType": "Arc",
      "gaugeStyle": "Full",
      "backStyle": "Full",
      "orientation": "vertical",
      "colorSchema": "Green to Red",
      "gaugeColorMode": "Labels",
      "colorsRange": [
        {
          "from": 0,
          "to": 50
        },
        {
          "from": 50,
          "to": 75
        },
        {
          "from": 75,
          "to": 100
        }
      ],
      "invertColors": false,
      "labels": {
        "show": true,
        "color": "black"
      },
      "scale": {
        "show": true,
        "labels": false,
        "color": "#333"
      },
      "type": "meter",
      "style": {
        "bgWidth": 0.9,
        "width": 0.9,
        "mask": false,
        "bgMask": false,
        "maskBars": 50,
        "bgFill": "#eee",
        "bgColor": false,
        "subText": "",
        "fontSize": 60,
        "labelColor": true
      }
    }
  },
  "aggs": [
    {
      "id": "1",
      "enabled": true,
      "type": "count",
      "schema": "metric",
      "params": {}
    }
  ]
}
EOF
}

resource "kibana_search" "china" {
	name 	        = "Chinese search"
	description     = "Chinese search results"
	display_columns = ["_source"]
	sort_by_columns = ["@timestamp"]
	search {
		index   = "${data.kibana_index.main.id}"
		filters {
				match {
					field_name = "geo.src"
					query      = "CN"
					type       = "phrase"
				}
		}
	}
}

data "kibana_index" "main" {
	filter {
		name = "title"
		values = ["logstash-*"]
	}
}
`

const testCreateVisualizationLogzioConfig = `
resource "kibana_visualization" "china_viz" {
	name 	            = "Chinese visualization"
	description         = "Chinese error visualization"
	saved_search_id     = "${kibana_search.china.id}"
	visualization_state = <<EOF
{
  "title": "Chinese search",
  "type": "gauge",
  "params": {
    "type": "gauge",
    "addTooltip": true,
    "addLegend": true,
    "gauge": {
      "verticalSplit": false,
      "extendRange": true,
      "percentageMode": false,
      "gaugeType": "Arc",
      "gaugeStyle": "Full",
      "backStyle": "Full",
      "orientation": "vertical",
      "colorSchema": "Green to Red",
      "gaugeColorMode": "Labels",
      "colorsRange": [
        {
          "from": 0,
          "to": 50
        },
        {
          "from": 50,
          "to": 75
        },
        {
          "from": 75,
          "to": 100
        }
      ],
      "invertColors": false,
      "labels": {
        "show": true,
        "color": "black"
      },
      "scale": {
        "show": true,
        "labels": false,
        "color": "#333"
      },
      "type": "meter",
      "style": {
        "bgWidth": 0.9,
        "width": 0.9,
        "mask": false,
        "bgMask": false,
        "maskBars": 50,
        "bgFill": "#eee",
        "bgColor": false,
        "subText": "",
        "fontSize": 60,
        "labelColor": true
      }
    }
  },
  "aggs": [
    {
      "id": "1",
      "enabled": true,
      "type": "count",
      "schema": "metric",
      "params": {}
    }
  ]
}
EOF
}

resource "kibana_search" "china" {
	name 	        = "Chinese search"
	description     = "Chinese search results"
	display_columns = ["_source"]
	sort_by_columns = ["@timestamp"]
	search {
		index   = "[logzioCustomerIndex]YYMMDD"
		filters {
				match {
					field_name = "geo.src"
					query      = "CN"
					type       = "phrase"
				}
		}
	}
}
`
const testUpdateVisualizationLogzioConfig = `
resource "kibana_visualization" "china_viz" {
	name 	            = "Chinese visualization - updated"
	description         = "Chinese error visualization - updated"
	saved_search_id     = "${kibana_search.china.id}"
	visualization_state = <<EOF
{
  "title": "Chinese search",
  "type": "gauge",
  "params": {
    "type": "gauge",
    "addTooltip": true,
    "addLegend": true,
    "gauge": {
      "verticalSplit": false,
      "extendRange": true,
      "percentageMode": false,
      "gaugeType": "Arc",
      "gaugeStyle": "Full",
      "backStyle": "Full",
      "orientation": "vertical",
      "colorSchema": "Green to Red",
      "gaugeColorMode": "Labels",
      "colorsRange": [
        {
          "from": 0,
          "to": 50
        },
        {
          "from": 50,
          "to": 75
        },
        {
          "from": 75,
          "to": 100
        }
      ],
      "invertColors": false,
      "labels": {
        "show": true,
        "color": "black"
      },
      "scale": {
        "show": true,
        "labels": false,
        "color": "#333"
      },
      "type": "meter",
      "style": {
        "bgWidth": 0.9,
        "width": 0.9,
        "mask": false,
        "bgMask": false,
        "maskBars": 50,
        "bgFill": "#eee",
        "bgColor": false,
        "subText": "",
        "fontSize": 60,
        "labelColor": true
      }
    }
  },
  "aggs": [
    {
      "id": "1",
      "enabled": true,
      "type": "count",
      "schema": "metric",
      "params": {}
    }
  ]
}
EOF
}

resource "kibana_search" "china" {
	name 	        = "Chinese search"
	description     = "Chinese search results"
	display_columns = ["_source"]
	sort_by_columns = ["@timestamp"]
    search {
		index   = "[logzioCustomerIndex]YYMMDD"
		filters {
				match {
					field_name = "geo.src"
					query      = "CN"
					type       = "phrase"
				}
		}
	}
}
`
