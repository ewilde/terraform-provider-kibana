[![Build Status](https://travis-ci.org/ewilde/go-kibana.svg?branch=master)](https://travis-ci.org/ewilde/go-kibana) [![GoDoc](https://godoc.org/github.com/ewilde/go-kibana?status.svg)](https://godoc.org/github.com/ewilde/go-kibana)

# go-kibana
go-kibana is a [go](https://golang.org/) client library for [kibana](https://github.com/elastic/kibana)

## Installation

```
go get github.com/ewilde/go-kibana
```

## Usage
```go
package examples

import (
	"github.com/ewilde/go-kibana"
	"github.com/stretchr/testify/assert"
)

func createSearch() (*kibana.SearchResponse, error) {
	client := kibana.NewClient(kibana.NewDefaultConfig())

	requestSearch, _ := kibana.NewSearchSourceBuilder().
		WithIndexId(client.Config.DefaultIndexId).
		WithFilter(&kibana.SearchFilter{
			Query: &kibana.SearchFilterQuery{
				Match: map[string]*kibana.SearchFilterQueryAttributes{
					"geo.src": {
						Query: "CN",
						Type:  "phrase",
					},
				},
			},
		}).
		Build()

	request, _ := kibana.NewSearchRequestBuilder().
		WithTitle("Geography filter on china with errors").
		WithDisplayColumns([]string{"_source"}).
		WithSortColumns([]string{"@timestamp"}, kibana.Descending).
		WithSearchSource(requestSearch).
		Build()

	return client.Search().Create(request)
}

func createVisualization(search *kibana.Search) (*kibana.Visualization, error) {
	client := kibana.NewClient(kibana.NewDefaultConfig())
	client.Config.KibanaVersion = kibana.DefaultKibanaVersion6

	request, _ := kibana.NewVisualizationRequestBuilder().
		WithTitle("Geography filter on china with errors").
		WithDescription("Gauge visualization based on a saved search").
		WithVisualizationState(`{"title":"Chinese search","type":"gauge","params":{"type":"gauge","addTooltip":true,"addLegend":true,"gauge":{"verticalSplit":false,"extendRange":true,"percentageMode":false,"gaugeType":"Arc","gaugeStyle":"Full","backStyle":"Full","orientation":"vertical","colorSchema":"Green to Red","gaugeColorMode":"Labels","colorsRange":[{"from":0,"to":50},{"from":50,"to":75},{"from":75,"to":100}],"invertColors":false,"labels":{"show":true,"color":"black"},"scale":{"show":true,"labels":false,"color":"#333"},"type":"meter","style":{"bgWidth":0.9,"width":0.9,"mask":false,"bgMask":false,"maskBars":50,"bgFill":"#eee","bgColor":false,"subText":"","fontSize":60,"labelColor":true}}},"aggs":[{"id":"1","enabled":true,"type":"count","schema":"metric","params":{}}]}`).
		WithSavedSearchId(search.Id).
		Build()

	return client.Visualization().Create(request)
}
```

### Displaying the filter with a saved search
By default a saved search won't display the `filter` on the search UI.
Use the `meta` structure to enable this, as shown below:

![image](https://user-images.githubusercontent.com/329397/36351467-714aeac2-14a2-11e8-9f83-225844da579e.png)

```go
client := DefaultTestKibanaClient()

	requestSearch, err := NewSearchSourceBuilder().
		WithIndexId(client.Config.DefaultIndexId).
		WithFilter(&SearchFilter{
			Query: &SearchFilterQuery{
				Match: map[string]*SearchFilterQueryAttributes{
					"geo.src": {
						Query: "CN",
						Type:  "phrase",
					},
				},
			},
			Meta: &SearchFilterMetaData{
				Index: client.Config.DefaultIndexId,
				Negate: false,
				Disabled: false,
				Alias: "China",
				Type: "phrase",
				Key: "geo.src",
				Value: "CN",
				Params: &SearchFilterQueryAttributes {
					Query: "CN",
					Type: "phrase",
				},
			},
		}).
		Build()

	request, err := NewSearchRequestBuilder().
		WithTitle("Geography filter on china").
		WithDisplayColumns([]string{"_source"}).
		WithSortColumns([]string{"@timestamp"}, Descending).
		WithSearchSource(requestSearch).
		Build()

	searchApi := client.Search()
	response, err := searchApi.Create(request)
```

### All Resources and Actions
Complete examples can be found in the [examples folder](examples) or
in the unit tests


## Developing
### Running test
**Logzio - running tests**

example:
```
env ELK_VERSION=5.5.3 KIBANA_TYPE=KibanaTypeLogzio \
    KIBANA_URI="https://app-eu.logz.io" \
    ELASTIC_SEARCH_PATH="/kibana/elasticsearch/logzioCustomerKibanaIndex" \
    LOGZ_CLIENT_ID=zzedfwe3424fsdf KIBANA_USERNAME=foo@acme.com \
    LOGZ_IO_ACCOUNT_ID_1=123233 \
    LOGZ_IO_ACCOUNT_ID_2=232333
    KIBANA_PASSWORD=mypwd make fmt test
```

| Environment variables           | Description                             |
|:----------------|:----------------------------------------|
| ELK_VERSION| Version of ELK to run while test against logzio |
| KIBANA_TYPE| Always  KibanaTypeLogzio|
| KIBANA_URI| Your logz.io base uri i.e. https://app-eu.logz.io |
| ELASTIC_SEARCH_PATH| Always /kibana/elasticsearch/logzioCustomerKibanaIndex for logz.io|
| LOGZ_CLIENT_ID| Obtained for the POST data call to https://logzio.auth0.com/oauth/ro. Use chrome developer tools when you login to logz.io to obtain this. |
| KIBANA_USERNAME| Your logz.io username|
| KIBANA_PASSWORD| Your logz.io password|
| LOGZ_IO_ACCOUNT_ID_1| *Optional* Your primary logz.io account id, you can obtain this from the result or GET https://app-eu.logz.io/session. If not given will not run some tests to do with switching between multiple logz.io accounts|
| LOGZ_IO_ACCOUNT_ID_2| *Optional* A secondary primary logz.io account id, you can obtain this from the result or GET https://app-eu.logz.io/session after you switch accounts in the logz.io UI. If not given will not run some tests to do with switching between multiple logz.io accounts|
| KIBANA_DEBUG| *Optional* If set to any value i.e. 1 will print http request and response debug information|
