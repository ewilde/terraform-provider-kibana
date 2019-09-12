package kibana

import (
	"encoding/json"
	"errors"
	"fmt"
)

var urlFromVersion = map[string]map[string]func(config *Config, name string, id string) string{
	DefaultKibanaVersion6: {
		"create_index": func(config *Config, name string, id string) string {
			return fmt.Sprintf("%s/api/saved_objects/index-pattern", config.KibanaBaseUri)
		},
		"refresh_index": func(config *Config, name string, id string) string {
			return fmt.Sprintf("%s/api/saved_objects/index-pattern/%s", config.KibanaBaseUri, id)
		},
	},
	"5.5.3": {
		"create_index": func(config *Config, name string, id string) string {
			return fmt.Sprintf("%s/es_admin/.kibana/index-pattern/%s/_create", config.KibanaBaseUri, name)
		},
		"refresh_index": func(config *Config, name string, id string) string {
			return fmt.Sprintf("%s/es_admin/.kibana/index-pattern/%s", config.KibanaBaseUri, name)
		},
	},
}

func getUrlFromVersion(version string, key string, config *Config, name string, id string) string {
	urlMap, ok := urlFromVersion[version]
	if !ok {
		urlMap = urlFromVersion[DefaultKibanaVersion6]
	}

	return urlMap[key](config, name, id)
}

type IndexPatternClient interface {
	SetDefault(indexPatternId string) error
	Create() (*IndexPatternCreateResult, error)
	RefreshFields(indexPatternId string) error
}

type IndexPatternClient553 struct {
	config *Config
	client *HttpAgent
}

type IndexPatternClient600 struct {
	config *Config
	client *HttpAgent
}

type IndexPattern struct {
	Attributes *IndexPatternAttributes `json:"attributes"`
}

type IndexPatternCreateResult struct {
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Version    version                 `json:"version"`
	Attributes *IndexPatternAttributes `json:"attributes"`
}

type IndexPatternCreateResult553 struct {
	Id      string  `json:"_id"`
	Type    string  `json:"_type"`
	Version version `json:"_version"`
}

type IndexPatternAttributes struct {
	Title         string `json:"title"`
	TimeFieldName string `json:"timeFieldName"`
	Fields        string `json:"fields"`
}

type metaFieldsResult struct {
	Fields []metaFields `json:"fields"`
}

type metaFields struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	Count             int    `json:"count"`
	Scripted          bool   `json:"scripted"`
	Searchable        bool   `json:"searchable"`
	Aggregatable      bool   `json:"aggregatable"`
	ReadFromDocValues bool   `json:"readFromDocValues"`
}

type valuePair struct {
	Value string `json:"value"`
}

func (api *IndexPatternClient600) SetDefault(indexPatternId string) error {
	response, body, err := api.client.Post(fmt.Sprintf("%s/api/kibana/settings/defaultIndex", api.config.KibanaBaseUri)).
		Set("kbn-version", api.config.KibanaVersion).
		Send(&valuePair{Value: indexPatternId}).
		End()

	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	return nil
}

func (api *IndexPatternClient600) Create() (*IndexPatternCreateResult, error) {
	uri := getUrlFromVersion(api.config.KibanaVersion, "create_index", api.config, "logstash-*", "")
	response, body, errs := api.client.Post(uri).
		Set("kbn-version", api.config.KibanaVersion).
		Send("{\"attributes\":{\"title\":\"logstash-*\",\"timeFieldName\":\"@timestamp\"}}").End()

	if errs != nil {
		return nil, errs[0]
	}

	if response.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	indexPatternCreateResult := &IndexPatternCreateResult{}
	err := json.Unmarshal([]byte(body), indexPatternCreateResult)
	if err != nil {
		return nil, fmt.Errorf("could not parse fields from index pattern create response, error: %v", err)
	}

	return indexPatternCreateResult, nil
}

func (api *IndexPatternClient600) RefreshFields(indexPatternId string) error {
	response, body, errs := api.client.Get(api.config.KibanaBaseUri + "/api/index_patterns/_fields_for_wildcard").
		Query(`{ "pattern" : "logstash-*", "meta_fields" : "[\"_source\",\"_id\",\"_type\",\"_index\",\"_score\"]" }`).End()

	if errs != nil {
		return errs[0]
	}

	if response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	fields := &metaFieldsResult{}
	err := json.Unmarshal([]byte(body), fields)
	if err != nil {
		return fmt.Errorf("could not parse fields for wildcard response, error: %v", err)
	}

	fieldJson, err := json.Marshal(fields.Fields)
	if err != nil {
		return fmt.Errorf("could not marshal fields from wildcard response, error: %v", err)
	}

	indexPattern := &IndexPattern{
		Attributes: &IndexPatternAttributes{
			Title:         "logstash-*",
			TimeFieldName: "@timestamp",
			Fields:        string(fieldJson),
		},
	}

	uri := getUrlFromVersion(api.config.KibanaVersion, "refresh_index", api.config, "logstash-*", indexPatternId)
	response, body, errs = api.client.Put(uri).
		Set("kbn-version", api.config.KibanaVersion).
		Send(indexPattern).
		End()

	if errs != nil {
		return errs[0]
	}

	if response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	return nil
}

func (api *IndexPatternClient553) SetDefault(indexPatternId string) error {
	response, body, err := api.client.Post(fmt.Sprintf("%s/api/kibana/settings/defaultIndex", api.config.KibanaBaseUri)).
		Set("kbn-version", api.config.KibanaVersion).
		Send(&valuePair{Value: indexPatternId}).
		End()

	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	return nil
}

func (api *IndexPatternClient553) Create() (*IndexPatternCreateResult, error) {
	uri := getUrlFromVersion(api.config.KibanaVersion, "create_index", api.config, "logstash-*", "")
	response, body, errs := api.client.Post(uri).
		Set("kbn-version", api.config.KibanaVersion).
		Send("{\"title\":\"logstash-*\",\"timeFieldName\":\"@timestamp\"}").End()

	if errs != nil {
		return nil, errs[0]
	}

	if response.StatusCode == 409 {
		return &IndexPatternCreateResult{
			Id:      "logstash-*",
			Type:    "index-pattern",
			Version: "1",
		}, nil
	} else if response.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	result := &IndexPatternCreateResult553{}

	err := json.Unmarshal([]byte(body), result)
	if err != nil {
		return nil, fmt.Errorf("could not parse fields from index pattern create response, error: %v", err)
	}

	return &IndexPatternCreateResult{
		Id:      result.Id,
		Type:    result.Type,
		Version: result.Version,
	}, nil

}

func (api *IndexPatternClient553) RefreshFields(indexPatternId string) error {
	response, body, errs := api.client.Get(api.config.KibanaBaseUri + "/api/index_patterns/_fields_for_wildcard").
		Query(`{ "pattern" : "logstash-*", "meta_fields" : "[\"_source\",\"_id\",\"_type\",\"_index\",\"_score\"]" }`).End()

	if errs != nil {
		return errs[0]
	}

	if response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	fields := &metaFieldsResult{}
	err := json.Unmarshal([]byte(body), fields)
	if err != nil {
		return fmt.Errorf("could not parse fields for wildcard response, error: %v", err)
	}

	fieldJson, err := json.Marshal(fields.Fields)
	if err != nil {
		return fmt.Errorf("could not marshal fields from wildcard response, error: %v", err)
	}

	indexPattern := &IndexPatternAttributes{
		Title:         "logstash-*",
		TimeFieldName: "@timestamp",
		Fields:        string(fieldJson),
	}

	uri := getUrlFromVersion(api.config.KibanaVersion, "refresh_index", api.config, "logstash-*", indexPatternId)
	response, body, errs = api.client.Put(uri).
		Set("kbn-version", api.config.KibanaVersion).
		Send(indexPattern).
		End()

	if errs != nil {
		return errs[0]
	}

	if response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	return nil
}
