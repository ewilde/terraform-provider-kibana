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
