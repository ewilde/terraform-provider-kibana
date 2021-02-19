package kibana

import (
	"encoding/json"
	"errors"
	"fmt"

	goversion "github.com/mcuadros/go-version"
	uuid "github.com/satori/go.uuid"
)

const (
	VisualizationReferencesTypeSearch visualizationReferencesType = "search"
	VisualizationReferencesTypeIndex  visualizationReferencesType = "index"
)

type visualizationReferencesType string

func (r visualizationReferencesType) String() string {
	return string(r)
}

type VisualizationClient interface {
	Create(request *CreateVisualizationRequest) (*Visualization, error)
	GetById(id string) (*Visualization, error)
	List() ([]*Visualization, error)
	Update(id string, request *UpdateVisualizationRequest) (*Visualization, error)
	Delete(id string) error
}

type CreateVisualizationRequest struct {
	Attributes *VisualizationAttributes   `json:"attributes"`
	References []*VisualizationReferences `json:"references,omitempty"`
}

type UpdateVisualizationRequest struct {
	Attributes *VisualizationAttributes   `json:"attributes"`
	References []*VisualizationReferences `json:"references,omitempty"`
}

type Visualization struct {
	Id         string                     `json:"id"`
	Type       string                     `json:"type"`
	Version    version                    `json:"version"`
	Attributes *VisualizationAttributes   `json:"attributes"`
	References []*VisualizationReferences `json:"references,omitempty"`
}

type VisualizationReferences struct {
	Name string                      `json:"name"`
	Type visualizationReferencesType `json:"type"`
	Id   string                      `json:"id"`
}

type VisualizationAttributes struct {
	Title                 string                       `json:"title"`
	Description           string                       `json:"description"`
	Version               int                          `json:"version"`
	VisualizationState    string                       `json:"visState"`
	SavedSearchId         string                       `json:"savedSearchId,omitempty"`
	SavedSearchRefName    string                       `json:"savedSearchRefName,omitempty"`
	KibanaSavedObjectMeta *SearchKibanaSavedObjectMeta `json:"kibanaSavedObjectMeta"`
}

type VisualizationRequestBuilder struct {
	title                 string
	description           string
	visualizationState    string
	savedSearchId         string
	savedSearchRefName    string
	kibanaSavedObjectMeta *SearchKibanaSavedObjectMeta
	references            []*VisualizationReferences
}

type visualizationClient600 struct {
	config *Config
	client *HttpAgent
}

type visualizationClient553 struct {
	config *Config
	client *HttpAgent
}

type visualizationReadResult553 struct {
	Id      string                   `json:"_id"`
	Type    string                   `json:"_type"`
	Version version                  `json:"_version"`
	Source  *VisualizationAttributes `json:"_source"`
}

func NewVisualizationRequestBuilder() *VisualizationRequestBuilder {
	return &VisualizationRequestBuilder{}
}

func (builder *VisualizationRequestBuilder) WithTitle(title string) *VisualizationRequestBuilder {
	builder.title = title
	return builder
}

func (builder *VisualizationRequestBuilder) WithDescription(description string) *VisualizationRequestBuilder {
	builder.description = description
	return builder
}

func (builder *VisualizationRequestBuilder) WithVisualizationState(visualizationState string) *VisualizationRequestBuilder {
	builder.visualizationState = visualizationState
	return builder
}

func (builder *VisualizationRequestBuilder) WithSavedSearchId(savedSearchId string) *VisualizationRequestBuilder {
	builder.savedSearchId = savedSearchId
	return builder
}

func (builder *VisualizationRequestBuilder) WithSavedSearchRefName(savedSearchRefName string) *VisualizationRequestBuilder {
	builder.savedSearchRefName = savedSearchRefName
	return builder
}

func (builder *VisualizationRequestBuilder) WithKibanaSavedObjectMeta(meta *SearchKibanaSavedObjectMeta) *VisualizationRequestBuilder {
	builder.kibanaSavedObjectMeta = meta
	return builder
}

func (builder *VisualizationRequestBuilder) WithReferences(refs []*VisualizationReferences) *VisualizationRequestBuilder {
	builder.references = refs
	return builder
}

func (builder *VisualizationRequestBuilder) Build(version string) (*CreateVisualizationRequest, error) {
	if goversion.Compare(version, "7.0.0", "<") {
		return &CreateVisualizationRequest{
			Attributes: &VisualizationAttributes{
				Title:                 builder.title,
				Description:           builder.description,
				SavedSearchId:         builder.savedSearchId,
				Version:               1,
				VisualizationState:    builder.visualizationState,
				KibanaSavedObjectMeta: builder.kibanaSavedObjectMeta,
			},
		}, nil
	} else {
		req := &CreateVisualizationRequest{
			Attributes: &VisualizationAttributes{
				Title:                 builder.title,
				Description:           builder.description,
				Version:               1,
				VisualizationState:    builder.visualizationState,
				KibanaSavedObjectMeta: builder.kibanaSavedObjectMeta,
				SavedSearchRefName:    "search_1",
			},
			References: []*VisualizationReferences{
				{
					Name: "search_1",
					Type: "search",
					Id:   builder.savedSearchId,
				},
			},
		}

		if len(builder.references) > 0 {
			req.References = builder.references
			req.Attributes.SavedSearchRefName = ""
		}

		return req, nil
	}
}

func (api *visualizationClient600) Create(request *CreateVisualizationRequest) (*Visualization, error) {
	response, body, err := api.client.
		Post(api.config.KibanaBaseUri+savedObjectsPath+"visualization?overwrite=true").
		Set("kbn-version", api.config.KibanaVersion).
		Send(request).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not create visualization")
	}

	createResponse := &Visualization{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from create visualization response, error: %v", error)
	}

	return createResponse, nil
}

func (api *visualizationClient600) GetById(id string) (*Visualization, error) {
	response, body, err := api.client.
		Get(api.config.KibanaBaseUri+savedObjectsPath+"visualization/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		if api.config.KibanaType == KibanaTypeLogzio && response.StatusCode >= 400 { // bug in their api reports missing visualization as bad request / server error
			response.StatusCode = 404
		}
		return nil, NewError(response, body, "Could not fetch visualization")
	}

	createResponse := &Visualization{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from get visualization response, error: %v", error)
	}

	return createResponse, nil
}

func (api *visualizationClient600) List() ([]*Visualization, error) {
	response, body, err := api.client.
		Get(api.config.KibanaBaseUri+savedObjectsPath+"_find?type=visualization&per_page=9999").
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		if api.config.KibanaType == KibanaTypeLogzio && response.StatusCode >= 400 { // bug in their api reports missing visualization as bad request / server error
			response.StatusCode = 404
		}
		return nil, NewError(response, body, "Could not fetch visualization")
	}

	var listResp = struct {
		SavedObjects []*Visualization `json:"saved_objects"`
	}{}
	var listErr error
	listErr = json.Unmarshal([]byte(body), &listResp)
	if listErr != nil {
		return nil, fmt.Errorf("could not parse fields from list visualization response, error: %v", listErr)
	}

	return listResp.SavedObjects, nil
}

func (api *visualizationClient600) Update(id string, request *UpdateVisualizationRequest) (*Visualization, error) {
	response, body, err := api.client.
		Post(api.config.KibanaBaseUri+savedObjectsPath+"visualization/"+id+"?overwrite=true").
		Set("kbn-version", api.config.KibanaVersion).
		Send(request).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not update visualization")
	}

	createResponse := &Visualization{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from update visualization response, error: %v", error)
	}

	return createResponse, nil
}

func (api *visualizationClient600) Delete(id string) error {
	response, body, err := api.client.
		Delete(api.config.KibanaBaseUri+savedObjectsPath+"visualization/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return NewError(response, body, "Could not delete visualization")
	}

	return nil
}

func (api *visualizationClient553) Create(request *CreateVisualizationRequest) (*Visualization, error) {
	id := uuid.NewV4().String()
	response, body, errs := api.client.
		Post(api.config.BuildFullPath("/%s/%s", "visualization", id)).
		Set("kbn-version", api.config.KibanaVersion).
		Send(request.Attributes).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not create visualization")
	}

	createResponse := &createResourceResult553{}
	err := json.Unmarshal([]byte(body), createResponse)
	if err != nil {
		return nil, fmt.Errorf("could not parse fields from create visualization response, error: %v", err)
	}

	return api.GetById(createResponse.Id)
}

func (api *visualizationClient553) GetById(id string) (*Visualization, error) {
	response, body, err := api.client.
		Get(api.config.BuildFullPath("/%s/%s", "visualization", id)).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		if api.config.KibanaType == KibanaTypeLogzio && response.StatusCode >= 400 { // bug in their api reports missing visualization as bad request / server error
			response.StatusCode = 404
		}

		return nil, NewError(response, body, "Could not fetch visualization")
	}

	createResponse := &visualizationReadResult553{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from visualization get response, error: %v", error)
	}

	return &Visualization{
		Id:         createResponse.Id,
		Version:    createResponse.Version,
		Type:       createResponse.Type,
		Attributes: createResponse.Source,
	}, nil
}

func (api *visualizationClient553) List() ([]*Visualization, error) {
	return nil, errors.New("not implemented")
}

func (api *visualizationClient553) Update(id string, request *UpdateVisualizationRequest) (*Visualization, error) {
	response, body, err := api.client.
		Post(api.config.BuildFullPath("/%s/%s", "visualization", id)).
		Set("kbn-version", api.config.KibanaVersion).
		Send(request.Attributes).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not update visualization")
	}

	createResponse := &createResourceResult553{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from update visualization response, error: %v", error)
	}

	return api.GetById(createResponse.Id)
}

func (api *visualizationClient553) Delete(id string) error {
	response, body, err := api.client.
		Delete(api.config.BuildFullPath("/%s/%s", "visualization", id)).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return NewError(response, body, "Could not delete visualization")
	}

	return nil
}
