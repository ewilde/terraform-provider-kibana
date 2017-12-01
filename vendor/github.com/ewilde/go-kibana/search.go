package kibana

import (
	"encoding/json"
	"fmt"
	"github.com/ewilde/go-kibana/containers"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

const (
	Ascending SortOrder = iota
	Descending
)

type SortOrder int

type SearchClient struct {
	config *Config
	client *gorequest.SuperAgent
}

type SearchRequest struct {
	Attributes *SearchAttributes `json:"attributes"`
}

type SearchResponse struct {
	Id         string            `json:"id"`
	Type       string            `json:"type"`
	Version    int               `json:"version"`
	Attributes *SearchAttributes `json:"attributes"`
}

type SearchAttributes struct {
	Title                 string                       `json:"title"`
	Description           string                       `json:"description"`
	Hits                  int                          `json:"hits"`
	Columns               []string                     `json:"columns"`
	Sort                  []string                     `json:"sort"`
	Version               int                          `json:"version"`
	KibanaSavedObjectMeta *SearchKibanaSavedObjectMeta `json:"kibanaSavedObjectMeta"`
}

type SearchKibanaSavedObjectMeta struct {
	SearchSourceJSON string `json:"searchSourceJSON"`
}

type SearchSource struct {
	IndexId      string          `json:"index"`
	HighlightAll bool            `json:"highlightAll"`
	Version      bool            `json:"version"`
	Query        *SearchQuery    `json:"query,omitempty"`
	Filter       []*SearchFilter `json:"filter"`
}

type SearchQuery struct {
	Query    string `json:"query"`
	Language string `json:"language"`
}

type SearchFilter struct {
	Query *SearchFilterQuery `json:"query"`
}

type SearchFilterQuery struct {
	Match map[string]*SearchFilterQueryAttributes `json:"match"`
}

type SearchFilterQueryAttributes struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}

type SearchRequestBuilder struct {
	title          string
	description    string
	displayColumns []string
	sortColumns    []string
	searchSource   *SearchSource
}

type SearchSourceBuilder struct {
	indexId      string
	highlightAll bool
	query        *SearchQuery
	filters      []*SearchFilter
}

func (api *SearchClient) Create(request *SearchRequest) (*SearchResponse, error) {
	response, body, err := api.client.
		Post(api.config.HostAddress+savedObjectsPath+"search?overwrite=true").
		Set("kbn-version", containers.KibanaVersion).
		Send(request).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	createResponse := &SearchResponse{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from create response, error: %v", error)
	}

	return createResponse, nil
}

func (api *SearchClient) GetById(id string) (*SearchResponse, error) {
	response, body, err := api.client.
		Get(api.config.HostAddress+savedObjectsPath+"search/"+id).
		Set("kbn-version", containers.KibanaVersion).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	createResponse := &SearchResponse{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from get response, error: %v", error)
	}

	return createResponse, nil
}

func NewSearchSourceBuilder() *SearchSourceBuilder {
	return &SearchSourceBuilder{filters: []*SearchFilter{}}
}

func (builder *SearchSourceBuilder) WithIndexId(indexId string) *SearchSourceBuilder {
	builder.indexId = indexId
	return builder
}

func (builder *SearchSourceBuilder) WithQuery(query *SearchQuery) *SearchSourceBuilder {
	builder.query = query
	return builder
}

func (builder *SearchSourceBuilder) WithFilter(filter *SearchFilter) *SearchSourceBuilder {
	builder.filters = append(builder.filters, filter)
	return builder
}

func (builder *SearchSourceBuilder) Build() (*SearchSource, error) {
	if builder.indexId == "" {
		return nil, errors.New("Index id is required to create a discover search source")
	}

	return &SearchSource{
		IndexId:      builder.indexId,
		HighlightAll: builder.highlightAll,
		Version:      true,
		Query:        builder.query,
		Filter:       builder.filters,
	}, nil
}

func NewRequestBuilder() *SearchRequestBuilder {
	return &SearchRequestBuilder{}
}

func (builder *SearchRequestBuilder) WithTitle(title string) *SearchRequestBuilder {
	builder.title = title
	return builder
}

func (builder *SearchRequestBuilder) WithDescription(description string) *SearchRequestBuilder {
	builder.description = description
	return builder
}

func (builder *SearchRequestBuilder) WithDisplayColumns(columns []string) *SearchRequestBuilder {
	builder.displayColumns = columns
	return builder
}

func (builder *SearchRequestBuilder) WithSortColumns(columns []string, order SortOrder) *SearchRequestBuilder {
	var sortOrder = ""
	if order == Descending {
		sortOrder = "desc"
	} else {
		sortOrder = "asc"
	}

	builder.sortColumns = append(columns, sortOrder)
	return builder
}

func (builder *SearchRequestBuilder) WithSearchSource(searchSource *SearchSource) *SearchRequestBuilder {
	builder.searchSource = searchSource
	return builder
}

func (builder *SearchRequestBuilder) Build() (*SearchRequest, error) {
	searchSourceBytes, err := json.Marshal(builder.searchSource)
	if err != nil {
		return nil, err
	}

	request := &SearchRequest{
		Attributes: &SearchAttributes{
			Title:       builder.title,
			Description: builder.description,
			Hits:        0,
			Columns:     builder.displayColumns,
			Sort:        builder.sortColumns,
			Version:     1,
			KibanaSavedObjectMeta: &SearchKibanaSavedObjectMeta{
				SearchSourceJSON: string(searchSourceBytes),
			},
		},
	}
	return request, nil
}
