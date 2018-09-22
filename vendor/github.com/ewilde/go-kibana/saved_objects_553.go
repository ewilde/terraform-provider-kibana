package kibana

import (
	"encoding/json"
	"errors"
	"fmt"
)

type savedObjectsClient553 struct {
	config  *Config
	client  *HttpAgent
	version string
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

func (api *savedObjectsClient553) GetByType(request *SavedObjectRequest) (*SavedObjectResponse, error) {
	address := api.config.BuildFullPath("/%s/_search?size=%d", request.Type, request.PerPage)
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
