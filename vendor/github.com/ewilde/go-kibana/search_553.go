package kibana

import (
	"encoding/json"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type searchClient553 struct {
	config  *Config
	client  *HttpAgent
	version string
}

type searchSourceBuilder553 struct {
	indexId      string
	highlightAll bool
	query        *SearchQuery553
	filters      []*SearchFilter
}

func (api *searchClient553) Version() string {
	return api.version
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
		if api.config.KibanaType == KibanaTypeLogzio && response.StatusCode >= 400 { // bug in their api reports missing search as bad request or server error
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

func (api *searchClient553) List() ([]*Search, error) {
	return nil, errors.New("not implemnted")
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
