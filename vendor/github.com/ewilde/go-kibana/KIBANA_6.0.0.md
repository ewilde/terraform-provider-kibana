# Search for saved searches
GET /api/saved_objects/?type=search&per_page=100&page=1&search_fields=title%5E3&search_fields=description

**Request body**
`{"query":{"match_all":{}}}`

**Response body**
```json
{
  "page": 1,
  "per_page": 100,
  "total": 1,
  "saved_objects": [
    {
      "id": "ccb36ef0-e272-11e7-9252-193185beeed2",
      "type": "search",
      "version": 1,
      "attributes": {
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
          "searchSourceJSON": "{\"index\":\"bc4272f0-e272-11e7-9252-193185beeed2\",\"highlightAll\":true,\"version\":true,\"query\":{\"query\":\"\",\"language\":\"lucene\"},\"filter\":[{\"meta\":{\"index\":\"bc4272f0-e272-11e7-9252-193185beeed2\",\"negate\":false,\"disabled\":false,\"alias\":null,\"type\":\"phrase\",\"key\":\"geo.src\",\"value\":\"CN\",\"params\":{\"query\":\"CN\",\"type\":\"phrase\"}},\"query\":{\"match\":{\"geo.src\":{\"query\":\"CN\",\"type\":\"phrase\"}}},\"$state\":{\"store\":\"appState\"}}]}"
        }
      }
    }
  ]
}
```

## Saved search
```json
{
  "index": "55590620-134c-11e8-8d71-2547028c7370",
  "highlightAll": true,
  "version": true,
  "query": {
    "query": "",
    "language": "lucene"
  },
  "filter": [
    {
      "meta": {
        "index": "55590620-134c-11e8-8d71-2547028c7370",
        "negate": false,
        "disabled": false,
        "alias": null,
        "type": "phrase",
        "key": "geo.src",
        "value": "CN",
        "params": {
          "query": "CN",
          "type": "phrase"
        }
      },
      "query": {
        "match": {
          "geo.src": {
            "query": "CN",
            "type": "phrase"
          }
        }
      },
      "$state": {
        "store": "appState"
      }
    },
    {
      "meta": {
        "index": "55590620-134c-11e8-8d71-2547028c7370",
        "negate": false,
        "disabled": false,
        "alias": null,
        "type": "phrase",
        "key": "@tags",
        "value": "error",
        "params": {
          "query": "error",
          "type": "phrase"
        }
      },
      "query": {
        "match": {
          "@tags": {
            "query": "error",
            "type": "phrase"
          }
        }
      },
      "$state": {
        "store": "appState"
      }
    }
  ]
}
```
