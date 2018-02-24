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

# Dashboard
POST http://localhost:5601/api/saved_objects/dashboard?overwrite=true

**Request body**
```json
{
  "attributes": {
    "title": "China dashboard",
    "hits": 0,
    "description": "China dashboard description",
    "panelsJSON": "[{\"size_x\":6,\"size_y\":3,\"panelIndex\":1,\"type\":\"visualization\",\"id\":\"bc8a1970-175b-11e8-accb-65182aaf9591\",\"col\":1,\"row\":1},{\"size_x\":6,\"size_y\":3,\"panelIndex\":2,\"type\":\"search\",\"id\":\"aca8b340-175b-11e8-accb-65182aaf9591\",\"col\":7,\"row\":1,\"columns\":[\"_source\"],\"sort\":[\"@timestamp\",\"desc\"]}]",
    "optionsJSON": "{\"darkTheme\":false}",
    "uiStateJSON": "{\"P-1\":{\"vis\":{\"defaultColors\":{\"0 - 50\":\"rgb(0,104,55)\",\"50 - 75\":\"rgb(255,255,190)\",\"75 - 100\":\"rgb(165,0,38)\"}}}}",
    "version": 1,
    "timeRestore": false,
    "kibanaSavedObjectMeta": {
      "searchSourceJSON": "{\"query\":{\"query\":\"\",\"language\":\"lucene\"},\"filter\":[],\"highlightAll\":true,\"version\":true}"
    }
  }
}
```


**Response body**
```json
{
  "id": "e41e1680-175b-11e8-accb-65182aaf9591",
  "type": "dashboard",
  "version": 1,
  "attributes": {
    "title": "China dashboard",
    "hits": 0,
    "description": "China dashboard description",
    "panelsJSON": "[{\"size_x\":6,\"size_y\":3,\"panelIndex\":1,\"type\":\"visualization\",\"id\":\"bc8a1970-175b-11e8-accb-65182aaf9591\",\"col\":1,\"row\":1},{\"size_x\":6,\"size_y\":3,\"panelIndex\":2,\"type\":\"search\",\"id\":\"aca8b340-175b-11e8-accb-65182aaf9591\",\"col\":7,\"row\":1,\"columns\":[\"_source\"],\"sort\":[\"@timestamp\",\"desc\"]}]",
    "optionsJSON": "{\"darkTheme\":false}",
    "uiStateJSON": "{\"P-1\":{\"vis\":{\"defaultColors\":{\"0 - 50\":\"rgb(0,104,55)\",\"50 - 75\":\"rgb(255,255,190)\",\"75 - 100\":\"rgb(165,0,38)\"}}}}",
    "version": 1,
    "timeRestore": false,
    "kibanaSavedObjectMeta": {
      "searchSourceJSON": "{\"query\":{\"query\":\"\",\"language\":\"lucene\"},\"filter\":[],\"highlightAll\":true,\"version\":true}"
    }
  }
}
```
