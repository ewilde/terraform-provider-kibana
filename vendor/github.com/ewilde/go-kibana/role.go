package kibana

import (
	"encoding/json"
	"fmt"
)

// RoleClient declares the required methods to implement to be a client and manage roles
type RoleClient interface {
	CreateOrUpdate(request *Role) error
	GetByID(id string) (*Role, error)
	Delete(id string) error
}

// Role is the api definition of a role in kibana
// can be used to create and get a role
type Role struct {
	Name              string                 `json:"name,omitempty"`
	Metadata          map[string]interface{} `json:"metadata"`
	TransientMetadata map[string]interface{} `json:"transient_metadata,omitempty"`
	ElasticSearch     *RoleElasticSearch     `json:"elasticsearch"`
	Kibana            []*RoleKibana          `json:"kibana"`
}

type RoleElasticSearch struct {
	Cluster []string      `json:"cluster"`
	Indices []interface{} `json:"indices"`
	RunAs   []string      `json:"run_as"`
}

type RoleKibana struct {
	Base    []string            `json:"base"`
	Feature map[string][]string `json:"feature"`
	Spaces  []string            `json:"spaces"`
}

// DefaultRoleClient structure to enable operations on Roles
// implements RoleClient
type DefaultRoleClient struct {
	config *Config
	client *HttpAgent
}

// CreateOrUpdate creates or updates a role
// based on https://www.elastic.co/guide/en/kibana/current/role-management-api-put.html
func (api *DefaultRoleClient) CreateOrUpdate(request *Role) error {
	id := request.Name
	request.Name = ""
	response, body, err := api.client.
		Put(api.config.KibanaBaseUri+"/api/security/role/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		Send(request).
		End()
	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 300 {
		return NewError(response, body, "Could not create role")
	}
	return nil
}

// GetByID fetch an existing role
// https://www.elastic.co/guide/en/kibana/current/role-management-api-get.html
func (api *DefaultRoleClient) GetByID(id string) (*Role, error) {
	response, body, err := api.client.
		Get(api.config.KibanaBaseUri+"/api/security/role/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		End()
	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, NewError(response, body, "Could not fetch role")
	}

	role := &Role{}
	error := json.Unmarshal([]byte(body), role)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from update visualization response, error: %v", error)
	}

	return role, nil

}

// Delete an existing Role
// based on https://www.elastic.co/guide/en/kibana/current/role-management-api-delete.html
func (api *DefaultRoleClient) Delete(id string) error {
	response, body, err := api.client.
		Delete(api.config.KibanaBaseUri+"/api/security/role/"+id).
		Set("kbn-version", api.config.KibanaVersion).
		End()
	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 300 {
		return NewError(response, body, "Could not fetch role")
	}
	return nil
}
