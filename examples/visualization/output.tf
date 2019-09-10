output "index_id" {
  value = "${data.kibana_index.main.id}"
}

output "search_china_id" {
  value = "${kibana_search.china.id}"
}


output "visualization_china_id" {
  value = "${kibana_search.china.id}"
}
