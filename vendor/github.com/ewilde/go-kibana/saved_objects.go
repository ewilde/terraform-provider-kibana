package kibana

import (
	"encoding/json"
	"errors"
	"fmt"
)

const savedObjectsPath = "/api/saved_objects/"

type SavedObjectRequest struct {
	Type    string   `json:"type" url:"type"`
	Fields  []string `json:"fields" url:"fields"`
	PerPage int      `json:"per_page" url:"per_page"`
}

type SavedObjectRequestBuilder struct {
	objectType string
	fields     []string
	perPage    int
}

type SavedObjectsClient interface {
	GetByType(request *SavedObjectRequest) (*SavedObjectResponse, error)
}

type SavedObjectResponse struct {
	Page         int            `json:"page"`
	PerPage      int            `json:"per_page"`
	Total        int            `json:"total"`
	SavedObjects []*SavedObject `json:"saved_objects"`
}

type SavedObject struct {
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	Version    int                    `json:"version"`
	Attributes map[string]interface{} `json:"attributes"`
}

type savedObjectsClient600 struct {
	config *Config
	client *HttpAgent
}

type savedObjectsClient553 struct {
	config *Config
	client *HttpAgent
}

type savedObjectSearchResponse553 struct {
	Hits savedObjectSearchResponseHits553 `json:"hits"`
}

type savedObjectSearchResponseHits553 struct {
	Total int                               `json:"total"`
	Hits  []savedObjectSearchResponseHit553 `json:"hits"`
}

type savedObjectSearchResponseHit553 struct {
	Id     string                 `json:"_id"`
	Type   string                 `json:"_type"`
	Source map[string]interface{} `json:"_source"`
}

func (api *savedObjectsClient600) GetByType(request *SavedObjectRequest) (*SavedObjectResponse, error) {
	address, err := addQueryString(api.config.KibanaBaseUri+savedObjectsPath, request)

	if err != nil {
		return nil, fmt.Errorf("could not build query string for get saved objects by type, error: %v", err)
	}

	apiResponse, body, errs := api.client.Get(address).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get saved objects, error: %v", errs)
	}

	if apiResponse.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", apiResponse.StatusCode, body))
	}

	response := &SavedObjectResponse{}
	err = json.Unmarshal([]byte(body), response)
	if err != nil {
		return nil, fmt.Errorf("could not parse saved objects response, error: %v, response body: %s", err, body)
	}

	return response, nil
}

func (api *savedObjectsClient553) GetByType(request *SavedObjectRequest) (*SavedObjectResponse, error) {
	address := fmt.Sprintf("%s/%s/_search?size=%d", api.config.KibanaBaseUri, api.getUriBase(request.Type), request.PerPage)
	apiResponse, body, errs := api.client.
		Post(address).
		Set("kbn-version", api.config.KibanaVersion).
		Send(`{"query":{"match_all":{}}}`).
		End()
	if errs != nil {
		return nil, fmt.Errorf("could not get saved objects, error: %v", errs)
	}

	if apiResponse.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", apiResponse.StatusCode, body))
	}

	response := &savedObjectSearchResponse553{}
	err := json.Unmarshal([]byte(body), response)
	if err != nil {
		return nil, fmt.Errorf("could not parse saved objects response, error: %v, response body: %s", err, body)
	}

	var savedObjects []*SavedObject
	for _, item := range response.Hits.Hits {
		version := 1
		if val, ok := item.Source["version"]; ok {
			version = val.(int)
		}

		savedObjects = append(savedObjects, &SavedObject{
			Type:       item.Type,
			Id:         item.Id,
			Version:    version,
			Attributes: item.Source,
		})
	}

	savedObjectResponse := &SavedObjectResponse{
		PerPage:      request.PerPage,
		Page:         1,
		Total:        response.Hits.Total,
		SavedObjects: savedObjects,
	}

	return savedObjectResponse, nil
}

func (api *savedObjectsClient553) getUriBase(action string) string {
	if api.config.KibanaType == KibanaTypeLogzio {
		return action
	}

	return fmt.Sprintf("es_admin/.kibana/%s", action)
}

func NewSavedObjectRequestBuilder() *SavedObjectRequestBuilder {
	return &SavedObjectRequestBuilder{perPage: 20}
}

func (builder *SavedObjectRequestBuilder) WithType(objectType string) *SavedObjectRequestBuilder {
	builder.objectType = objectType
	return builder
}

func (builder *SavedObjectRequestBuilder) WithFields(fields []string) *SavedObjectRequestBuilder {
	builder.fields = fields
	return builder
}

func (builder *SavedObjectRequestBuilder) WithPerPage(perPage int) *SavedObjectRequestBuilder {
	builder.perPage = perPage
	return builder
}

func (builder *SavedObjectRequestBuilder) Build() *SavedObjectRequest {
	return &SavedObjectRequest{
		Fields:  builder.fields,
		Type:    builder.objectType,
		PerPage: builder.perPage,
	}
}
