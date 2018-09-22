package kibana

import (
	"encoding/json"
	"errors"
	"fmt"
)

type searchClient600 struct {
	config  *Config
	client  *HttpAgent
	version string
}

type searchSourceBuilder600 struct {
	indexId      string
	highlightAll bool
	query        *SearchQuery600
	filters      []*SearchFilter
}

func (api *searchClient600) Version() string {
	return api.version
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
