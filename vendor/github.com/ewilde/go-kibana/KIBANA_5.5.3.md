# Search for index patterns
/es_admin/.kibana/index-pattern/_search?stored_fields=
**Request body**
`{"query":{"match_all":{}},"size":100}`

**Response body**
```json
{
  "took": 3,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "failed": 0
  },
  "hits": {
    "total": 1,
    "max_score": 1,
    "hits": [
      {
        "_index": ".kibana",
        "_type": "index-pattern",
        "_id": "logstash-*",
        "_score": 1
      }
    ]
  }
}
```

# Search for saved searches
/es_admin/.kibana/search/_search?size=100

**Request body**
`{"query":{"match_all":{}}}`

**Response body**
```json
{
  "took": 0,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "failed": 0
  },
  "hits": {
    "total": 1,
    "max_score": 1,
    "hits": [
      {
        "_index": ".kibana",
        "_type": "search",
        "_id": "61d38250-e272-11e7-80e1-13cb872a4917",
        "_score": 1,
        "_source": {
          "title": "China",
          "description": "",
          "hits": 0,
          "columns": [
            "_source"
          ],
          "sort": [
            "@timestamp",
            "desc"
          ],
          "version": 1,
          "kibanaSavedObjectMeta": {
            "searchSourceJSON": "{\"index\":\"logstash-*\",\"highlightAll\":true,\"version\":true,\"query\":{\"match_all\":{}},\"filter\":[{\"meta\":{\"index\":\"logstash-*\",\"negate\":false,\"disabled\":false,\"alias\":null,\"type\":\"phrase\",\"key\":\"geo.src\",\"value\":\"CN\"},\"query\":{\"match\":{\"geo.src\":{\"query\":\"CN\",\"type\":\"phrase\"}}},\"$state\":{\"store\":\"appState\"}}]}"
          }
        }
      }
    ]
  }
}
```
