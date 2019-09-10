provider "kibana" {
  version        = "~> 0.3"
  kibana_version = "6.2.1"
}

data "kibana_index" "main" {
  filter {
    name   = "title"
    values = ["logstash-*"]
  }
}

resource "kibana_visualization" "china_viz" {
  name            = "Chinese visualization"
  description     = "Chinese error visualization"
  saved_search_id = "${kibana_search.china.id}"

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
  name            = "Chinese origin - errors"
  description     = "Errors occured when source was from china"
  display_columns = ["_source"]
  sort_by_columns = ["@timestamp"]

  search {
    index = "${data.kibana_index.main.id}"

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
