package kibana

import (
	"encoding/json"
	"fmt"
)

// SpaceClient declares the required methods to implement to be a client and manage spaces
type SpaceClient interface {
	Create(request *Space) error
	Update(request *Space) error
	GetByID(id string) (*Space, error)
	Delete(id string) error
}

// Space is the api definition of a space in kibana
// can be used to create, update and get a space
type Space struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	Color            string   `json:"color,omitempty"`
	Initials         string   `json:"initials,omitempty"`
	ImageUrl         string   `json:"imageUrl,omitempty"`
	DisabledFeatures []string `json:"disabledFeatures,omitempty"`
}

// DefaultSpaceClient structure to enable operations on Spaces
// implements SpaceClient
type DefaultSpaceClient struct {
	config *Config
	client *HttpAgent
}

// Create creates a space
// based on https://www.elastic.co/guide/en/kibana/current/spaces-api-post.html
func (api *DefaultSpaceClient) Create(request *Space) error {
	response, body, err := api.client.
		Post(api.config.KibanaBaseUri+"/api/spaces/space").
		Set("kbn-version", api.config.KibanaVersion).
		Send(request).
		End()
	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 300 {
		return NewError(response, body, "Could not create space")
	}
	return nil
}

// Update updates a space
// based on https://www.elastic.co/guide/en/kibana/master/spaces-api-put.html
func (api *DefaultSpaceClient) Update(request *Space) error {
	id := request.Id
	response, body, err := api.client.
		Put(api.config.KibanaBaseUri+"/api/spaces/space/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		Send(request).
		End()
	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 300 {
		return NewError(response, body, "Could not update space")
	}
	return nil
}

// GetByID fetch an existing space
// https://www.elastic.co/guide/en/kibana/master/spaces-api-get.html
func (api *DefaultSpaceClient) GetByID(id string) (*Space, error) {
	response, body, err := api.client.
		Get(api.config.KibanaBaseUri+"/api/spaces/space/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		End()
	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not fetch space")
	}

	space := &Space{}
	error := json.Unmarshal([]byte(body), space)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from update visualization response, error: %v", error)
	}

	return space, nil

}

// Delete an existing space
// based on https://www.elastic.co/guide/en/kibana/master/spaces-api-delete.html
func (api *DefaultSpaceClient) Delete(id string) error {
	response, body, err := api.client.
		Delete(api.config.KibanaBaseUri+"/api/spaces/space/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		End()
	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 300 {
		return NewError(response, body, "Could not delete space")
	}
	return nil
}
