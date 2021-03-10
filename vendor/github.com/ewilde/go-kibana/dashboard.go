package kibana

import (
	"encoding/json"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// Enums for DashboardReferencesType
const (
	DashboardReferencesTypeSearch        dashboardReferencesType = "search"
	DashboardReferencesTypeVisualization dashboardReferencesType = "visualization"
)

type dashboardReferencesType string

func (r dashboardReferencesType) String() string {
	return string(r)
}

type DashboardClient interface {
	Create(request *CreateDashboardRequest) (*Dashboard, error)
	GetById(id string) (*Dashboard, error)
	List() ([]*Dashboard, error)
	Update(id string, request *UpdateDashboardRequest) (*Dashboard, error)
	Delete(id string) error
}

type CreateDashboardRequest struct {
	Attributes *DashboardAttributes   `json:"attributes"`
	References []*DashboardReferences `json:"references,omitempty"`
}

type UpdateDashboardRequest struct {
	Attributes *DashboardAttributes   `json:"attributes"`
	References []*DashboardReferences `json:"references,omitempty"`
}

type Dashboard struct {
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	Version    version                `json:"version"`
	Attributes *DashboardAttributes   `json:"attributes"`
	References []*DashboardReferences `json:"references,omitempty"`
}

type DashboardReferences struct {
	Name string                  `json:"name"`
	Type dashboardReferencesType `json:"type"`
	Id   string                  `json:"id"`
}

type DashboardAttributes struct {
	Title                 string                       `json:"title"`
	Description           string                       `json:"description"`
	Version               int                          `json:"version"`
	PanelsJson            string                       `json:"panelsJSON"`
	OptionsJson           string                       `json:"optionsJSON"`
	UiStateJSON           string                       `json:"uiStateJSON,omitempty"`
	TimeRestore           bool                         `json:"timeRestore"`
	KibanaSavedObjectMeta *SearchKibanaSavedObjectMeta `json:"kibanaSavedObjectMeta"`
}

type DashboardRequestBuilder struct {
	title                 string
	description           string
	panelsJson            string
	optionsJson           string
	uiStateJson           string
	timeRestore           bool
	kibanaSavedObjectMeta *SearchKibanaSavedObjectMeta
	references            []*DashboardReferences
}

type dashboardClient600 struct {
	config *Config
	client *HttpAgent
}

type dashboardClient553 struct {
	config *Config
	client *HttpAgent
}

type dashboardReadResult553 struct {
	Id      string               `json:"_id"`
	Type    string               `json:"_type"`
	Version version              `json:"_version"`
	Source  *DashboardAttributes `json:"_source"`
}

func NewDashboardRequestBuilder() *DashboardRequestBuilder {
	return &DashboardRequestBuilder{}
}

func (builder *DashboardRequestBuilder) WithTitle(title string) *DashboardRequestBuilder {
	builder.title = title
	return builder
}

func (builder *DashboardRequestBuilder) WithDescription(description string) *DashboardRequestBuilder {
	builder.description = description
	return builder
}

func (builder *DashboardRequestBuilder) WithPanelsJson(panelsJson string) *DashboardRequestBuilder {
	builder.panelsJson = panelsJson
	return builder
}

func (builder *DashboardRequestBuilder) WithOptionsJson(optionsJson string) *DashboardRequestBuilder {
	builder.optionsJson = optionsJson
	return builder
}

func (builder *DashboardRequestBuilder) WithUiStateJson(uiStateJson string) *DashboardRequestBuilder {
	builder.uiStateJson = uiStateJson
	return builder
}

func (builder *DashboardRequestBuilder) WithTimeRestore(timeRestore bool) *DashboardRequestBuilder {
	builder.timeRestore = timeRestore
	return builder
}

func (builder *DashboardRequestBuilder) WithKibanaSavedObjectMeta(meta *SearchKibanaSavedObjectMeta) *DashboardRequestBuilder {
	builder.kibanaSavedObjectMeta = meta
	return builder
}

func (builder *DashboardRequestBuilder) WithReferences(refs []*DashboardReferences) *DashboardRequestBuilder {
	builder.references = refs
	return builder
}

func (builder *DashboardRequestBuilder) Build() (*CreateDashboardRequest, error) {

	return &CreateDashboardRequest{
		Attributes: &DashboardAttributes{
			Title:                 builder.title,
			Description:           builder.description,
			PanelsJson:            builder.panelsJson,
			OptionsJson:           builder.optionsJson,
			UiStateJSON:           builder.uiStateJson,
			TimeRestore:           builder.timeRestore,
			KibanaSavedObjectMeta: builder.kibanaSavedObjectMeta,
		},
		References: builder.references,
	}, nil
}

func (api *dashboardClient600) Create(request *CreateDashboardRequest) (*Dashboard, error) {
	response, body, err := api.client.
		Post(api.config.KibanaBaseUri+savedObjectsPath+"dashboard?overwrite=true").
		Set("kbn-version", api.config.KibanaVersion).
		Send(request).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not create dashboard")
	}

	createResponse := &Dashboard{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from create dashboard response, error: %v", error)
	}

	return createResponse, nil
}

func (api *dashboardClient600) GetById(id string) (*Dashboard, error) {
	response, body, err := api.client.
		Get(api.config.KibanaBaseUri+savedObjectsPath+"dashboard/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		if api.config.KibanaType == KibanaTypeLogzio && response.StatusCode >= 400 { // bug in their api reports missing dashboard as bad request / server error
			response.StatusCode = 404
		}
		return nil, NewError(response, body, "Could not fetch dashboard")
	}

	createResponse := &Dashboard{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from get dashboard response, error: %v", error)
	}

	return createResponse, nil
}

func (api *dashboardClient600) List() ([]*Dashboard, error) {
	response, body, err := api.client.
		Get(api.config.KibanaBaseUri+savedObjectsPath+"_find?type=dashboard&per_page=9999").
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		if api.config.KibanaType == KibanaTypeLogzio && response.StatusCode >= 400 { // bug in their api reports missing dashboard as bad request / server error
			response.StatusCode = 404
		}
		return nil, NewError(response, body, "Could not list dashboards")
	}

	var listResp = struct {
		SavedObjects []*Dashboard `json:"saved_objects"`
	}{}
	var listErr error
	listErr = json.Unmarshal([]byte(body), &listResp)
	if listErr != nil {
		return nil, fmt.Errorf("could not parse fields from list dashboard response, error: %v", listErr)
	}

	return listResp.SavedObjects, nil
}

func (api *dashboardClient600) Update(id string, request *UpdateDashboardRequest) (*Dashboard, error) {
	response, body, err := api.client.
		Post(api.config.KibanaBaseUri+savedObjectsPath+"dashboard/"+id+"?overwrite=true").
		Set("kbn-version", api.config.KibanaVersion).
		Send(request).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not update dashboard")
	}

	createResponse := &Dashboard{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from update dashboard response, error: %v", error)
	}

	return createResponse, nil
}

func (api *dashboardClient600) Delete(id string) error {
	response, body, err := api.client.
		Delete(api.config.KibanaBaseUri+savedObjectsPath+"dashboard/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return NewError(response, body, "Could not delete dashboard")
	}

	return nil
}

func (api *dashboardClient553) Create(request *CreateDashboardRequest) (*Dashboard, error) {
	id := uuid.NewV4().String()
	response, body, errs := api.client.
		Post(api.config.BuildFullPath("/%s/%s", "dashboard", id)).
		Set("kbn-version", api.config.KibanaVersion).
		Send(request.Attributes).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not create dashboard")
	}

	createResponse := &createResourceResult553{}
	err := json.Unmarshal([]byte(body), createResponse)
	if err != nil {
		return nil, fmt.Errorf("could not parse fields from create dashboard response, error: %v", err)
	}

	return api.GetById(createResponse.Id)
}

func (api *dashboardClient553) GetById(id string) (*Dashboard, error) {
	response, body, err := api.client.
		Get(api.config.BuildFullPath("/%s/%s", "dashboard", id)).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		if api.config.KibanaType == KibanaTypeLogzio && response.StatusCode >= 400 { // bug in their api reports missing dashboard as bad request / server error
			response.StatusCode = 404
		}

		return nil, NewError(response, body, "Could not fetch dashboard")
	}

	createResponse := &dashboardReadResult553{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from dashboard get response, error: %v", error)
	}

	return &Dashboard{
		Id:         createResponse.Id,
		Version:    createResponse.Version,
		Type:       createResponse.Type,
		Attributes: createResponse.Source,
	}, nil
}

func (api *dashboardClient553) List() ([]*Dashboard, error) {
	return nil, errors.New("not implemnted")
}

func (api *dashboardClient553) Update(id string, request *UpdateDashboardRequest) (*Dashboard, error) {
	response, body, err := api.client.
		Post(api.config.BuildFullPath("/%s/%s", "dashboard", id)).
		Set("kbn-version", api.config.KibanaVersion).
		Send(request.Attributes).
		End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not update dashboard")
	}

	createResponse := &createResourceResult553{}
	error := json.Unmarshal([]byte(body), createResponse)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from update dashboard response, error: %v", error)
	}

	return api.GetById(createResponse.Id)
}

func (api *dashboardClient553) Delete(id string) error {
	response, body, err := api.client.
		Delete(api.config.BuildFullPath("/%s/%s", "dashboard", id)).
		Set("kbn-version", api.config.KibanaVersion).
		End()

	if err != nil {
		return NewError(response, body, "Could not delete dashboard")
	}

	return nil
}
