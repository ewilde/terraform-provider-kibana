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

	request, _ := kibana.NewRequestBuilder().
		WithTitle("Geography filter on china with errors").
		WithDisplayColumns([]string{"_source"}).
		WithSortColumns([]string{"@timestamp"}, kibana.Descending).
		WithSearchSource(requestSearch).
		Build()

	return client.Search().Create(request)
}
```

### All Resources and Actions
Complete examples can be found in the [examples folder](examples) or
in the unit tests
