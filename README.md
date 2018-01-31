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

Clone repository to: `$GOPATH/src/github.com/hashicorp/terraform-provider-kibana`

```sh
$ mkdir -p $GOPATH/src/github.com/hashicorp; cd $GOPATH/src/github.com/hashicorp
$ git clone git@github.com:hashicorp/terraform-provider-kibana
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/hashicorp/terraform-provider-kibana
$ make build
```

Using the provider
----------------------
## Example creating saved search resources
```json
provider "kibana" {
}

data "kibana_index" "main" {
  filter = {
    name = "title"
    values = ["logstash-*"]
  }
}

resource "kibana_search" "china" {
  name 	        = "Chinese origin - errors"
  description     = "Errors occured when source was from china"
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
```

More examples can be found in the [example folder](examples)

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.9+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-kibana
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
