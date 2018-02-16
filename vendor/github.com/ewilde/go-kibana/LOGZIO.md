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

# Seach for saved visualization
POST /kibana/elasticsearch/logzioCustomerKibanaIndex/visualization/_search?size=100
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
    "total": 48,
    "max_score": 1,
    "hits": [
      {
        "_index": "logzioCustomerKibanaIndex",
        "_type": "visualization",
        "_id": "ELB-Count-number-of-4xx-and-5xx-Backend-Responses-per-URI-and-User-Agent",
        "_score": 1,
        "_source": {
          "description": "",
          "version": 1,
          "kibanaSavedObjectMeta": {
            "searchSourceJSON": "{\"index\":\"[logzioCustomerIndex]YYMMDD\",\"query\":{\"query_string\":{\"query\":\"type: elb AND backend_status_code:[400 599]\",\"analyze_wildcard\":true}},\"filter\":[]}"
          },
          "visState": "{\"type\":\"table\",\"params\":{\"perPage\":10,\"showPartialRows\":false,\"showMeticsAtAllLevels\":false},\"aggs\":[{\"id\":\"1\",\"type\":\"count\",\"schema\":\"metric\",\"params\":{}},{\"id\":\"3\",\"type\":\"terms\",\"schema\":\"bucket\",\"params\":{\"field\":\"request\",\"size\":5,\"order\":\"desc\",\"orderBy\":\"1\"}},{\"id\":\"4\",\"type\":\"terms\",\"schema\":\"bucket\",\"params\":{\"field\":\"backend_status_code\",\"size\":5,\"order\":\"desc\",\"orderBy\":\"1\"}},{\"id\":\"5\",\"type\":\"terms\",\"schema\":\"bucket\",\"params\":{\"field\":\"os\",\"size\":5,\"order\":\"desc\",\"orderBy\":\"1\"}},{\"id\":\"6\",\"type\":\"terms\",\"schema\":\"bucket\",\"params\":{\"field\":\"device\",\"size\":5,\"order\":\"desc\",\"orderBy\":\"1\"}}],\"listeners\":{}}",
          "title": "ELB Count number of 4xx and 5xx Backend Responses per URI and User Agent",
          "_logzioOriginalAppId": 36
        }
      },
      {
        "_index": "logzioCustomerKibanaIndex",
        "_type": "visualization",
        "_id": "ELB-5XX-Responses",
        "_score": 1,
        "_source": {
          "description": "",
          "version": 1,
          "kibanaSavedObjectMeta": {
            "searchSourceJSON": "{\"index\":\"[logzioCustomerIndex]YYMMDD\",\"query\":{\"query_string\":{\"query\":\"type: elb AND elb_status_code:[500 599]\",\"analyze_wildcard\":true}},\"filter\":[]}"
          },
          "visState": "{\"type\":\"table\",\"params\":{\"perPage\":10,\"showPartialRows\":false,\"showMeticsAtAllLevels\":false},\"aggs\":[{\"id\":\"1\",\"type\":\"count\",\"schema\":\"metric\",\"params\":{}},{\"id\":\"2\",\"type\":\"terms\",\"schema\":\"bucket\",\"params\":{\"field\":\"elb_status_code\",\"size\":5,\"order\":\"desc\",\"orderBy\":\"1\"}},{\"id\":\"3\",\"type\":\"terms\",\"schema\":\"bucket\",\"params\":{\"field\":\"request\",\"size\":5,\"order\":\"desc\",\"orderBy\":\"1\"}}],\"listeners\":{}}",
          "title": "ELB 5XX Responses",
          "_logzioOriginalAppId": 35
        }
      }
    ]
  }
}
```

# Create visualization
POST /kibana/elasticsearch/logzioCustomerKibanaIndex/visualization/0d41e0b0-0658-11e8-8859-6f62fb52e8a9

**Request body**

```json
{
  "title": "test kong vis",
  "visState": "{\"title\":\"test kong vis\",\"type\":\"area\",\"params\":{\"grid\":{\"categoryLines\":false,\"style\":{\"color\":\"#eee\"}},\"categoryAxes\":[{\"id\":\"CategoryAxis-1\",\"type\":\"category\",\"position\":\"bottom\",\"show\":true,\"style\":{},\"scale\":{\"type\":\"linear\"},\"labels\":{\"show\":true,\"truncate\":100},\"title\":{\"text\":\"@timestamp date ranges\"}}],\"valueAxes\":[{\"id\":\"ValueAxis-1\",\"name\":\"LeftAxis-1\",\"type\":\"value\",\"position\":\"left\",\"show\":true,\"style\":{},\"scale\":{\"type\":\"linear\",\"mode\":\"normal\"},\"labels\":{\"show\":true,\"rotate\":0,\"filter\":false,\"truncate\":100},\"title\":{\"text\":\"Count\"}}],\"seriesParams\":[{\"show\":\"true\",\"type\":\"area\",\"mode\":\"stacked\",\"data\":{\"label\":\"Count\",\"id\":\"1\"},\"drawLinesBetweenPoints\":true,\"showCircles\":true,\"interpolate\":\"linear\",\"valueAxis\":\"ValueAxis-1\"}],\"addTooltip\":true,\"addLegend\":true,\"legendPosition\":\"right\",\"times\":[],\"addTimeMarker\":false},\"aggs\":[{\"id\":\"1\",\"enabled\":true,\"type\":\"count\",\"schema\":\"metric\",\"params\":{}},{\"id\":\"2\",\"enabled\":true,\"type\":\"date_range\",\"schema\":\"segment\",\"params\":{\"field\":\"@timestamp\",\"ranges\":[{\"from\":\"now-1h\",\"to\":\"now\"}]}}],\"listeners\":{}}",
  "uiStateJSON": "{}",
  "description": "",
  "savedSearchId": "cab13e00-a8e5-11e7-8c62-75ca9e6062e7",
  "version": 1,
  "kibanaSavedObjectMeta": {
    "searchSourceJSON": "{\"filter\":[]}"
  },
  "_createdBy": {
    "userId": 19429,
    "fullName": "Edward Wilde",
    "username": "edward.wilde@form3.tech"
  },
  "_createdAt": 1517383588431,
  "_updatedBy": {
    "userId": 19429,
    "fullName": "Edward Wilde",
    "username": "edward.wilde@form3.tech"
  },
  "_updatedAt": 1517383588431
}
```

**Reponse body**
```json
{
  "_index": "logzioCustomerKibanaIndex",
  "_type": "visualization",
  "_id": "0d41e0b0-0658-11e8-8859-6f62fb52e8a9",
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
