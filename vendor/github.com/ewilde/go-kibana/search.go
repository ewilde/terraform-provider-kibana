package kibana

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

const (
	Ascending SortOrder = iota
	Descending
)

type SortOrder int

type SearchSourceBuilderFactory interface {
	NewSearchSource() SearchSourceBuilder
}

type SearchClient interface {
	Create(request *CreateSearchRequest) (*Search, error)
	Update(id string, request *UpdateSearchRequest) (*Search, error)
	GetById(id string) (*Search, error)
	Delete(id string) error
	NewSearchSource() SearchSourceBuilder
}

type searchClient600 struct {
	config *Config
	client *HttpAgent
}

type searchClient553 struct {
	config *Config
	client *HttpAgent
}

type CreateSearchRequest struct {
	Attributes *SearchAttributes `json:"attributes"`
}

type UpdateSearchRequest struct {
	Attributes *SearchAttributes `json:"attributes"`
}

type Search struct {
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

type searchReadResult553 struct {
	Id      string            `json:"_id"`
	Type    string            `json:"_type"`
	Version int               `json:"_version"`
	Source  *SearchAttributes `json:"_source"`
}

type SearchKibanaSavedObjectMeta struct {
	SearchSourceJSON string `json:"searchSourceJSON"`
}

type SearchSource struct {
	IndexId      string          `json:"index"`
	HighlightAll bool            `json:"highlightAll"`
	Version      bool            `json:"version"`
	Query        interface{}     `json:"query,omitempty"`
	Filter       []*SearchFilter `json:"filter"`
}

type SearchQuery600 struct {
	Query    string `json:"query"`
	Language string `json:"language"`
}

type SearchQuery553 struct {
	QueryString *searchQueryString `json:"query_string"`
}

type searchQueryString struct {
	Query   string `json:"query"`
	Analyze bool   `json:"analyze_wildcard"`
}

type SearchFilter struct {
	Query  *SearchFilterQuery    `json:"query"`
	Exists *SearchFilterExists   `json:"exists"`
	Meta   *SearchFilterMetaData `json:"meta,omitempty"`
}

type SearchFilterQuery struct {
	Match map[string]*SearchFilterQueryAttributes `json:"match"`
}

type SearchFilterExists struct {
	Field string `json:"field"`
}

type SearchFilterMetaData struct {
	Index    string                       `json:"index"`
	Negate   bool                         `json:"negate"`
	Disabled bool                         `json:"disabled"`
	Alias    string                       `json:"alias"`
	Type     string                       `json:"type"`
	Key      string                       `json:"key"`
	Value    string                       `json:"value"`
	Params   *SearchFilterQueryAttributes `json:"params"`
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

type SearchSourceBuilder interface {
	WithIndexId(indexId string) SearchSourceBuilder
	WithQuery(query string) SearchSourceBuilder
	WithFilter(filter *SearchFilter) SearchSourceBuilder
	Build() (*SearchSource, error)
}

type searchSourceBuilder600 struct {
	indexId      string
	highlightAll bool
	query        *SearchQuery600
	filters      []*SearchFilter
}

type searchSourceBuilder553 struct {
	indexId      string
	highlightAll bool
	query        *SearchQuery553
	filters      []*SearchFilter
}

func (api *searchClient600) NewSearchSource() SearchSourceBuilder {
	return &searchSourceBuilder600{filters: []*SearchFilter{}}
}

func (api *searchClient600) Create(request *CreateSearchRequest) (*Search, error) {
	response, body, err := api.client.
		Post(api.config.KibanaBaseUri+savedObjectsPath+"search?overwrite=true").
		Set("kbn-version", api.config.KibanaVersion).
		Send(request).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not create search")
	}

	createResponse := &Search{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from create search response, error: %v", error)
	}

	return createResponse, nil
}

func (api *searchClient600) Update(id string, request *UpdateSearchRequest) (*Search, error) {
	response, body, err := api.client.
		Post(api.config.KibanaBaseUri+savedObjectsPath+"search/"+id+"?overwrite=true").
		Set("kbn-version", api.config.KibanaVersion).
		Send(request).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not update search")
	}

	createResponse := &Search{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from update search response, error: %v", error)
	}

	return createResponse, nil
}

func (api *searchClient600) GetById(id string) (*Search, error) {
	response, body, err := api.client.
		Get(api.config.KibanaBaseUri+savedObjectsPath+"search/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not fetch search")
	}

	createResponse := &Search{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from get search response, error: %v", error)
	}

	return createResponse, nil
}

func (api *searchClient600) Delete(id string) error {
	response, body, err := api.client.
		Delete(api.config.KibanaBaseUri+savedObjectsPath+"search/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return NewError(response, body, "Could not delete search")
	}

	return nil
}

func (api *searchClient553) NewSearchSource() SearchSourceBuilder {
	return &searchSourceBuilder553{filters: []*SearchFilter{}}
}

func (api *searchClient553) Create(request *CreateSearchRequest) (*Search, error) {
	id := uuid.NewV4().String()
	response, body, errs := api.client.
		Post(api.config.BuildFullPath("/%s/%s", "search", id)).
		Set("kbn-version", api.config.KibanaVersion).
		Send(request.Attributes).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not create search")
	}

	createResponse := &createResourceResult553{}
	err := json.Unmarshal([]byte(body), createResponse)
	if err != nil {
		return nil, fmt.Errorf("could not parse fields from create search response, error: %v", err)
	}

	return &Search{
		Id:         createResponse.Id,
		Version:    createResponse.Version,
		Attributes: request.Attributes,
		Type:       createResponse.Type,
	}, nil
}

func (api *searchClient553) Update(id string, request *UpdateSearchRequest) (*Search, error) {
	response, body, err := api.client.
		Post(api.config.BuildFullPath("/%s/%s", "search", id)).
		Set("kbn-version", api.config.KibanaVersion).
		Send(request.Attributes).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not update search")
	}

	createResponse := &createResourceResult553{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from update search response, error: %v", error)
	}

	return &Search{
		Id:         id,
		Version:    createResponse.Version,
		Attributes: request.Attributes,
		Type:       createResponse.Type,
	}, nil
}

func (api *searchClient553) GetById(id string) (*Search, error) {
	response, body, err := api.client.
		Get(api.config.BuildFullPath("/%s/%s", "search", id)).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		if api.config.KibanaType == KibanaTypeLogzio && response.StatusCode == 400 { // bug in their api reports missing search as bad request
			response.StatusCode = 404
		}

		return nil, NewError(response, body, "Could not fetch search")
	}

	createResponse := &searchReadResult553{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from get response, error: %v", error)
	}

	return &Search{
		Id:         createResponse.Id,
		Version:    createResponse.Version,
		Type:       createResponse.Type,
		Attributes: createResponse.Source,
	}, nil
}

func (api *searchClient553) Delete(id string) error {
	response, body, err := api.client.
		Delete(api.config.BuildFullPath("/%s/%s", "search", id)).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return NewError(response, body, "Could not delete search")
	}

	return nil
}

func (builder *searchSourceBuilder600) WithIndexId(indexId string) SearchSourceBuilder {
	builder.indexId = indexId
	return builder
}

func (builder *searchSourceBuilder600) WithQuery(query string) SearchSourceBuilder {
	builder.query = &SearchQuery600{Query: query, Language: "lucene"}
	return builder
}

func (builder *searchSourceBuilder600) WithFilter(filter *SearchFilter) SearchSourceBuilder {
	builder.filters = append(builder.filters, filter)
	return builder
}

func (builder *searchSourceBuilder600) Build() (*SearchSource, error) {
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

func (builder *searchSourceBuilder553) WithIndexId(indexId string) SearchSourceBuilder {
	builder.indexId = indexId
	return builder
}

func (builder *searchSourceBuilder553) WithQuery(query string) SearchSourceBuilder {
	builder.query = &SearchQuery553{
		QueryString: &searchQueryString{
			Query:   query,
			Analyze: true,
		},
	}
	return builder
}

func (builder *searchSourceBuilder553) WithFilter(filter *SearchFilter) SearchSourceBuilder {
	builder.filters = append(builder.filters, filter)
	return builder
}

func (builder *searchSourceBuilder553) Build() (*SearchSource, error) {
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

func NewSearchRequestBuilder() *SearchRequestBuilder {
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

func (builder *SearchRequestBuilder) Build() (*CreateSearchRequest, error) {
	searchSourceBytes, err := json.Marshal(builder.searchSource)
	if err != nil {
		return nil, err
	}

	request := &CreateSearchRequest{
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
