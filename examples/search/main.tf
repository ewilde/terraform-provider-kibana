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
