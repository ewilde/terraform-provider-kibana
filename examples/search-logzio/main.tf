provider "kibana" {
  elastic_search_path = "/kibana/elasticsearch/logzioCustomerKibanaIndex"
  kibana_type         = "KibanaTypeLogzio"
  kibana_version      = "5.5.3"
  kibana_uri          = "https://app-eu.logz.io"
  kibana_username     = "${var.kibana_username}"
  kibana_password     = "${var.kibana_password}"
  logzio_client_id    = "${var.logzio_client_id}"
  logzio_account_id   = "${var.logzio_account_id}"
}

resource "kibana_search" "china" {
  name            = "Chinese origin - errors"
  description     = "Errors occured when source was from china"
  display_columns = ["_source"]
  sort_by_columns = ["@timestamp"]

  search {
    index = "[logzioCustomerIndex]YYMMDD"

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
