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

type SearchClient interface {
	Create(request *CreateSearchRequest) (*Search, error)
	Update(id string, request *UpdateSearchRequest) (*Search, error)
	GetById(id string) (*Search, error)
	Delete(id string) error
}
type SearchClient600 struct {
	config *Config
	client *HttpAgent
}

type SearchClient553 struct {
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

type SearchCreateResult553 struct {
	Id      string `json:"_id"`
	Type    string `json:"_type"`
	Version int    `json:"_version"`
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

func (api *SearchClient600) Create(request *CreateSearchRequest) (*Search, error) {
	response, body, err := api.client.
		Post(api.config.KibanaBaseUri+savedObjectsPath+"search?overwrite=true").
		Set("kbn-version", api.config.KibanaVersion).
		Send(request).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	createResponse := &Search{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from create response, error: %v", error)
	}

	return createResponse, nil
}

func (api *SearchClient600) Update(id string, request *UpdateSearchRequest) (*Search, error) {
	response, body, err := api.client.
		Post(api.config.KibanaBaseUri+savedObjectsPath+"search/"+id+"?overwrite=true").
		Set("kbn-version", api.config.KibanaVersion).
		Send(request).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	createResponse := &Search{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from create response, error: %v", error)
	}

	return createResponse, nil
}

func (api *SearchClient600) GetById(id string) (*Search, error) {
	response, body, err := api.client.
		Get(api.config.KibanaBaseUri+savedObjectsPath+"search/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	createResponse := &Search{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from get response, error: %v", error)
	}

	return createResponse, nil
}

func (api *SearchClient600) Delete(id string) error {
	response, body, err := api.client.
		Delete(api.config.KibanaBaseUri+savedObjectsPath+"search/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return fmt.Errorf("could not delete search with id:%s response: %+v %s errors: %+v", id, response, body, err)
	}

	return nil
}

func (api *SearchClient553) Create(request *CreateSearchRequest) (*Search, error) {
	id := uuid.NewV4().String()
	response, body, errs := api.client.
		Post(fmt.Sprintf("%s/%s/%s", api.config.KibanaBaseUri, api.getUriBase("search"), id)).
		Set("kbn-version", api.config.KibanaVersion).
		Send(request.Attributes).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	if response.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	createResponse := &SearchCreateResult553{}
	err := json.Unmarshal([]byte(body), createResponse)
	if err != nil {
		return nil, fmt.Errorf("could not parse fields from create response, error: %v", err)
	}

	return &Search{
		Id:         createResponse.Id,
		Version:    createResponse.Version,
		Attributes: request.Attributes,
		Type:       createResponse.Type,
	}, nil
}

func (api *SearchClient553) Update(id string, request *UpdateSearchRequest) (*Search, error) {
	response, body, err := api.client.
		Post(fmt.Sprintf("%s/%s/%s", api.config.KibanaBaseUri, api.getUriBase("search"), id)).
		Set("kbn-version", api.config.KibanaVersion).
		Send(request.Attributes).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	createResponse := &SearchCreateResult553{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from create response, error: %v", error)
	}

	return &Search{
		Id:         id,
		Version:    createResponse.Version,
		Attributes: request.Attributes,
		Type:       createResponse.Type,
	}, nil
}

func (api *SearchClient553) GetById(id string) (*Search, error) {
	response, body, err := api.client.
		Get(fmt.Sprintf("%s/%s/%s", api.config.KibanaBaseUri, api.getUriBase("search"), id)).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	createResponse := &Search{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from get response, error: %v", error)
	}

	return createResponse, nil
}

func (api *SearchClient553) Delete(id string) error {
	response, body, err := api.client.
		Delete(fmt.Sprintf("%s/%s/%s", api.config.KibanaBaseUri, api.getUriBase("search"), id)).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return fmt.Errorf("could not delete search with id:%s response: %+v %s errors: %+v", id, response, body, err)
	}

	return nil
}

func (api *SearchClient553) getUriBase(action string) string {
	if api.config.KibanaType == KibanaTypeLogzio {
		return action
	}

	return fmt.Sprintf("es_admin/.kibana/%s", action)
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
