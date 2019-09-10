package kibana

import (
	"fmt"
	"testing"

	kibana "github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"strings"
)

var testDashboardCreate = map[kibana.KibanaType]string{
	kibana.KibanaTypeVanilla: fmt.Sprintf(testCreateDashboardConfig, "${data.kibana_index.main.id}", dataKibanaIndex),
	kibana.KibanaTypeLogzio:  fmt.Sprintf(testCreateDashboardConfig, "[logzioCustomerIndex]YYMMDD", ""),
}

var testDashboardUpdate = map[kibana.KibanaType]string{
	kibana.KibanaTypeVanilla: fmt.Sprintf(testUpdateDashboardConfig, "${data.kibana_index.main.id}", dataKibanaIndex),
	kibana.KibanaTypeLogzio:  fmt.Sprintf(testUpdateDashboardConfig, "[logzioCustomerIndex]YYMMDD", ""),
}

func TestAccKibanaDashboardApi(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKibanaDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDashboardCreate[testConfig.KibanaType],
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaDashboardExists("kibana_dashboard.china_dash"),
					resource.TestCheckResourceAttr("kibana_dashboard.china_dash", "name", "Chinese dashboard"),
					resource.TestCheckResourceAttr("kibana_dashboard.china_dash", "description", "Chinese dashboard description"),
				),
			},
			{
				Config: testDashboardUpdate[testConfig.KibanaType],
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaDashboardExists("kibana_dashboard.china_dash"),
					resource.TestCheckResourceAttr("kibana_dashboard.china_dash", "name", "Chinese dashboard - updated"),
					resource.TestCheckResourceAttr("kibana_dashboard.china_dash", "description", "Chinese dashboard description - updated"),
				),
			},
		},
	})
}

func testAccCheckKibanaDashboardDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*kibana.KibanaClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "kibana_dashboard" {
			continue
		}

		response, err := client.Dashboard().GetById(rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("error calling get dashboard by id: %v", err)
		}

		if response != nil {
			return fmt.Errorf("dashboard %s still exists, %+v", rs.Primary.ID, response)
		}
	}

	return nil
}

func testAccCheckKibanaDashboardExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		api, err := testAccProvider.Meta().(*kibana.KibanaClient).Dashboard().GetById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("dashboard with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateDashboardConfig = `
resource "kibana_dashboard" "china_dash" {
	name = "Chinese dashboard"
	description = "Chinese dashboard description"
	panels_json = <<EOF
[
  {
    "gridData": {
      "w": 6,
      "h": 3,
      "x": 0,
      "y": 0,
      "i": "1"
    },
	"version": "6.2.1",
    "panelIndex": "1",
    "type": "visualization",
    "id": "${kibana_visualization.china_viz.id}"
  },
  {
    "gridData": {
      "w": 6,
      "h": 3,
      "x": 6,
      "y": 0,
      "i": "2"
    },
	"version": "6.2.1",
    "panelIndex": "2",
    "type": "search",
    "id": "${kibana_search.china.id}"
  }
]
EOF
}


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
		index   = "%s"
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
const testUpdateDashboardConfig = `
resource "kibana_dashboard" "china_dash" {
	name = "Chinese dashboard - updated"
	description = "Chinese dashboard description - updated"
	panels_json = <<EOF
[
  {
    "gridData": {
      "w": 6,
      "h": 3,
      "x": 0,
      "y": 0,
      "i": "1"
    },
	"version": "6.2.1",
    "panelIndex": "1",
    "type": "visualization",
    "id": "${kibana_visualization.china_viz.id}"
  },
  {
    "gridData": {
      "w": 6,
      "h": 3,
      "x": 6,
      "y": 0,
      "i": "2"
    },
	"version": "6.2.1",
    "panelIndex": "2",
    "type": "search",
    "id": "${kibana_search.china.id}"
  }
]
EOF
}


resource "kibana_visualization" "china_viz" {
	name 	            = "Chinese visualization updated"
	description         = "Chinese error visualization  updated"
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
const dataKibanaIndex = `
data "kibana_index" "main" {
  filter {
    name   = "title"
    values = ["logstash-*"]
  }
}
`
