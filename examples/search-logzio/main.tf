provider "kibana" {
    kibana_type    = "KibanaTypeLogzio"
    kibana_version = "5.5.3"
    kibana_uri     = "https://app-eu.logz.io/kibana/elasticsearch/logzioCustomerKibanaIndex"
}

resource "kibana_search" "china" {
  name 	        = "Chinese origin - errors"
  description     = "Errors occured when source was from china"
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
