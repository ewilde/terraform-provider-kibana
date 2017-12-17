# Search for index patterns
POST /kibana/elasticsearch/logzioCustomerKibanaIndex/index-pattern/_search?stored_fields=

**Request body**
`{"query":{"match_all":{}},"size":100}`

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
        "_index": "logzioCustomerKibanaIndex",
        "_type": "index-pattern",
        "_id": "[logzioCustomerIndex]YYMMDD",
        "_score": 1
      }
    ]
  }
}
```

# Search for saved searches
POST /kibana/elasticsearch/logzioCustomerKibanaIndex/search/_search?size=100

**Request body**
`{"query":{"match_all":{}}}`

**Response body**
```json
{
  "took": 1,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "failed": 0
  },
  "hits": {
    "total": 69,
    "max_score": 1,
    "hits": [
      {
        "_index": "logzioCustomerKibanaIndex",
        "_type": "search",
        "_id": "3d29c470-62f0-11e7-b189-0f8cdb432680",
        "_score": 1,
        "_source": {
          "title": "application logs (noseflute)",
          "description": "",
          "hits": 0,
          "columns": [
            "message",
            "type",
            "stack"
          ],
          "sort": [
            "@timestamp",
            "desc"
          ],
          "version": 1,
          "kibanaSavedObjectMeta": {
            "searchSourceJSON": "{\"index\":\"[logzioCustomerIndex]YYMMDD\",\"highlightAll\":true,\"filter\":[{\"meta\":{\"negate\":false,\"index\":\"[logzioCustomerIndex]YYMMDD\",\"key\":\"stack\",\"value\":\"noseflute\",\"disabled\":false,\"alias\":null},\"query\":{\"match\":{\"stack\":{\"query\":\"noseflute\",\"type\":\"phrase\"}}},\"$state\":{\"store\":\"appState\"}},{\"meta\":{\"negate\":true,\"index\":\"[logzioCustomerIndex]YYMMDD\",\"key\":\"type\",\"value\":\"docker-stats\",\"disabled\":false,\"alias\":null},\"query\":{\"match\":{\"type\":{\"query\":\"docker-stats\",\"type\":\"phrase\"}}},\"$state\":{\"store\":\"appState\"}},{\"meta\":{\"negate\":true,\"index\":\"[logzioCustomerIndex]YYMMDD\",\"key\":\"type\",\"value\":\"metricsets\",\"disabled\":false,\"alias\":null},\"query\":{\"match\":{\"type\":{\"query\":\"metricsets\",\"type\":\"phrase\"}}},\"$state\":{\"store\":\"appState\"}},{\"meta\":{\"negate\":true,\"index\":\"[logzioCustomerIndex]YYMMDD\",\"key\":\"type\",\"value\":\"amazon-ecs-agent\",\"disabled\":false,\"alias\":null},\"query\":{\"match\":{\"type\":{\"query\":\"amazon-ecs-agent\",\"type\":\"phrase\"}}},\"$state\":{\"store\":\"appState\"}},{\"meta\":{\"negate\":true,\"index\":\"[logzioCustomerIndex]YYMMDD\",\"key\":\"type\",\"value\":\"tech.form3/userapi\",\"disabled\":false,\"alias\":null},\"query\":{\"match\":{\"type\":{\"query\":\"tech.form3/userapi\",\"type\":\"phrase\"}}},\"$state\":{\"store\":\"appState\"}},{\"meta\":{\"negate\":true,\"index\":\"[logzioCustomerIndex]YYMMDD\",\"key\":\"type\",\"value\":\"tech.form3/paymentapi\",\"disabled\":false,\"alias\":null},\"query\":{\"match\":{\"type\":{\"query\":\"tech.form3/paymentapi\",\"type\":\"phrase\"}}},\"$state\":{\"store\":\"appState\"}},{\"meta\":{\"negate\":true,\"index\":\"[logzioCustomerIndex]YYMMDD\",\"key\":\"logger_name\",\"value\":\"tech.form3.corelib.aws.queues.PollingQueueListener\",\"disabled\":false,\"alias\":null},\"query\":{\"match\":{\"logger_name\":{\"query\":\"tech.form3.corelib.aws.queues.PollingQueueListener\",\"type\":\"phrase\"}}},\"$state\":{\"store\":\"appState\"}},{\"meta\":{\"negate\":true,\"index\":\"[logzioCustomerIndex]YYMMDD\",\"key\":\"type\",\"value\":\"tech.form3/consul-agent\",\"disabled\":false,\"alias\":null},\"query\":{\"match\":{\"type\":{\"query\":\"tech.form3/consul-agent\",\"type\":\"phrase\"}}},\"$state\":{\"store\":\"appState\"}},{\"meta\":{\"negate\":true,\"index\":\"[logzioCustomerIndex]YYMMDD\",\"key\":\"logger_name\",\"value\":\"tech.form3.corelib.aws.queues.ScheduledQueueListener\",\"disabled\":false,\"alias\":null},\"query\":{\"match\":{\"logger_name\":{\"query\":\"tech.form3.corelib.aws.queues.ScheduledQueueListener\",\"type\":\"phrase\"}}},\"$state\":{\"store\":\"appState\"}}],\"query\":{\"query_string\":{\"analyze_wildcard\":true,\"query\":\"*\"}}}"
          },
          "_createdBy": {
            "userId": 19430,
            "fullName": "Steve Cook",
            "username": "steve.cook@form3.tech"
          },
          "_createdAt": 1499416961600,
          "_updatedBy": {
            "userId": 19430,
            "fullName": "Steve Cook",
            "username": "steve.cook@form3.tech"
          },
          "_updatedAt": 1499416961600
        }
      }
    ]
  }
}
```

# Create search
POST /kibana/elasticsearch/logzioCustomerKibanaIndex/search/9c2f2320-e252-11e7-96f8-397bd34fab6c

**Request body**

```json
{
  "title": "test s",
  "description": "",
  "hits": 0,
  "columns": [
    "message"
  ],
  "sort": [
    "@timestamp",
    "desc"
  ],
  "version": 1,
  "kibanaSavedObjectMeta": {
    "searchSourceJSON": "{\"index\":\"[logzioCustomerIndex]YYMMDD\",\"highlightAll\":true,\"version\":true,\"query\":{\"query_string\":{\"query\":\"message:\\\"GET\\\"\",\"analyze_wildcard\":true}},\"filter\":[]}"
  },
  "_createdBy": {
    "userId": 19429,
    "fullName": "Edward Wilde",
    "username": "edward.wilde@form3.tech"
  },
  "_createdAt": 1513423009383,
  "_updatedBy": {
    "userId": 19429,
    "fullName": "Edward Wilde",
    "username": "edward.wilde@form3.tech"
  },
  "_updatedAt": 1513423009383
}
```

**Response**
```json
{
  "_index": "logzioCustomerKibanaIndex",
  "_type": "search",
  "_id": "9c2f2320-e252-11e7-96f8-397bd34fab6c",
  "_version": 1,
  "result": "created",
  "_shards": {
    "total": 2,
    "successful": 2,
    "failed": 0
  },
  "created": true
}
```
