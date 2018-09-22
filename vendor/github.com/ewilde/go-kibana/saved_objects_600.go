package kibana

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mcuadros/go-version"
)

type savedObjectsClient600 struct {
	config  *Config
	client  *HttpAgent
	version string
}

func (api *savedObjectsClient600) GetByType(request *SavedObjectRequest) (*SavedObjectResponse, error) {
	address, err := addQueryString(api.getSavedObjectsPath(), request)

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

func (api *savedObjectsClient600) getSavedObjectsPath() string {
	if version.Compare(api.config.KibanaVersion, "6.3.0", ">=") {
		return api.config.KibanaBaseUri + savedObjectsPath + "_find"

	}

	return api.config.KibanaBaseUri + savedObjectsPath
}
