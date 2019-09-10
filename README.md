[![Build Status](https://travis-ci.org/ewilde/terraform-provider-kibana.svg?branch=master)](https://travis-ci.org/ewilde/terraform-provider-kibana)

Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.9 (to build the provider plugin)

Usage
---------------------

```
# For example, restrict kibana version in 1.0.x
provider "kibana" {
  version = "~> 1.0"
}
```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/ewilde/terraform-provider-kibana`

```sh
$ mkdir -p $GOPATH/src/github.com/ewilde; cd $GOPATH/src/github.com/ewilde
$ git clone git@github.com:ewilde/terraform-provider-kibana
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/ewilde/terraform-provider-kibana
$ ELK_VERSION=6.2.3 KIBANA_TYPE=KibanaTypeVanilla make build
```

Using the provider
----------------------
## Example creating saved search, visualization and dashboard resources
```hcl
provider "kibana" {}

data "kibana_index" "main" {
  filter {
    name   = "title"
    values = ["logstash-*"]
  }
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

resource "kibana_visualization" "china_viz" {
  name            = "Chinese visualization - updated"
  description     = "Chinese error visualization - updated"
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

resource "kibana_dashboard" "china_dash" {
  name        = "Chinese dashboard"
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
```

More examples can be found in the [example folder](examples)

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.9+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ ELK_VERSION=6.2.3 KIBANA_TYPE=KibanaTypeVanilla make build
...
$ $GOPATH/bin/terraform-provider-kibana
...
```

In order to test the provider, you can simply run `make test`. Note that `ELK_VERSION` and `KIBANA_TYPE` are used to control the test targets. The full list of test targets is visible in `.travis.yml`.

```sh
$ ELK_VERSION=6.2.3 KIBANA_TYPE=KibanaTypeVanilla make test
```

In order to run the full suite of Acceptance tests, run `ELK_VERSION=6.2.3 KIBANA_TYPE=KibanaTypeVanilla make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

## Debuging
Set `KIBANA_DEBUG=1` to see http debug output

```sh
$ make testacc
```

## Adding dependencies
This project uses [govendor](https://github.com/kardianos/govendor) to manage dependencies

### Add /Update a package
`govendor fetch github.com/owner/repo`

*Recursive*
`govendor fetch github.com/owner/repo/...`
