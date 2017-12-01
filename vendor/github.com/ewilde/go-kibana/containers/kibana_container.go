package containers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/ory-am/dockertest.v3"
	"log"
)

const KibanaVersion = "6.0.0"

type kibanaContainer struct {
	Name     string
	pool     *dockertest.Pool
	resource *dockertest.Resource
	Uri      string
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

type indexPattern struct {
	Attributes *indexPatternAttributes `json:"attributes"`
}

type indexPatternCreateResult struct {
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Version    int                     `json:"version"`
	Attributes *indexPatternAttributes `json:"attributes"`
}

type indexPatternBulkGet struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type indexPatternAttributes struct {
	Title         string `json:"title"`
	TimeFieldName string `json:"timeFieldName"`
	Fields        string `json:"fields"`
}

type valuePair struct {
	Value string `json:"value"`
}

func NewKibanaContainer(pool *dockertest.Pool, elasticSearch *elasticSearchContainer) (container *kibanaContainer, indexId string) {
	envVars := []string{
		fmt.Sprintf("ELASTICSEARCH_URL=http://%s:9200", elasticSearch.Name),
	}

	options := &dockertest.RunOptions{
		Name:         "kibana",
		Repository:   "docker.elastic.co/kibana/kibana-oss",
		Tag:          KibanaVersion,
		Env:          envVars,
		Links:        []string{elasticSearch.Name},
		ExposedPorts: []string{"5601"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5601/tcp": {{HostIP: "", HostPort: "5601"}},
		},
	}

	resource, err := pool.RunWithOptions(options)
	kibanaUri := fmt.Sprintf("http://localhost:%v", resource.GetPort("5601/tcp"))

	var indexPatternCreateResult *indexPatternCreateResult
	if err := pool.Retry(func() error {
		client := gorequest.New()

		var error error
		if error = checkKibanaServiceIsStarted(client, kibanaUri); error != nil {
			return error
		}

		indexPatternCreateResult, error = createIndexPattern(client, kibanaUri)
		if error != nil {
			return error
		}

		if error := updateIndexPatternFields(client, kibanaUri, indexPatternCreateResult.Id); error != nil {
			return error
		}

		if error := setDefaultIndexPattern(client, kibanaUri, indexPatternCreateResult.Id); error != nil {
			return error
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to kibana: %s", err)
	}

	if err != nil {
		log.Fatalf("Could not connect to kibana: %s", err)
	}

	name := getContainerName(resource)
	log.Printf("Kibana (%v): up", name)

	return &kibanaContainer{
		Name:     name,
		pool:     pool,
		resource: resource,
		Uri:      kibanaUri,
	}, indexPatternCreateResult.Id
}

func setDefaultIndexPattern(client *gorequest.SuperAgent, kibanaUri string, indexPatternId string) error {
	response, body, err := client.Post(fmt.Sprintf("%s/api/kibana/settings/defaultIndex", kibanaUri)).
		Set("kbn-version", KibanaVersion).
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

func createIndexPattern(client *gorequest.SuperAgent, kibanaUri string) (*indexPatternCreateResult, error) {
	response, body, err := client.Post(kibanaUri+"/api/saved_objects/index-pattern").
		Set("kbn-version", KibanaVersion).
		Send("{\"attributes\":{\"title\":\"logstash-*\",\"timeFieldName\":\"@timestamp\"}}").End()

	if err != nil {
		return nil, err[0]
	}

	if response.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	indexPatternCreateResult := &indexPatternCreateResult{}
	error := json.Unmarshal([]byte(body), indexPatternCreateResult)
	if error != nil {
		return nil, fmt.Errorf("could not parse fields from index pattern create response, error: %v", error)
	}

	return indexPatternCreateResult, nil
}

func updateIndexPatternFields(client *gorequest.SuperAgent, kibanaUri string, indexPatternId string) error {
	response, body, err := client.Get(kibanaUri + "/api/index_patterns/_fields_for_wildcard").
		Query(`{ "pattern" : "logstash-*", "meta_fields" : "[\"_source\",\"_id\",\"_type\",\"_index\",\"_score\"]" }`).End()

	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	fields := &metaFieldsResult{}
	error := json.Unmarshal([]byte(body), fields)
	if error != nil {
		return fmt.Errorf("could not parse fields for wildcard response, error: %v", error)
	}

	fieldJson, error := json.Marshal(fields.Fields)
	if error != nil {
		return fmt.Errorf("could not marshal fields from wildcard response, error: %v", error)
	}

	indexPattern := &indexPattern{
		Attributes: &indexPatternAttributes{
			Title:         "logstash-*",
			TimeFieldName: "@timestamp",
			Fields:        string(fieldJson),
		},
	}

	response, body, err = client.Put(fmt.Sprintf("%s/api/saved_objects/index-pattern/%s", kibanaUri, indexPatternId)).
		Set("kbn-version", KibanaVersion).
		Send(indexPattern).
		End()

	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	return nil
}

func bulkGetPost(client *gorequest.SuperAgent, kibanaUri string, indexPatternCreateResult *indexPatternCreateResult) error {
	response, body, err := client.Post(kibanaUri+"/api/saved_objects/bulk_get").
		Set("kbn-version", KibanaVersion).
		Send(&[]*indexPatternBulkGet{{Id: indexPatternCreateResult.Id, Type: indexPatternCreateResult.Type}}).End()

	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	return nil
}

func checkKibanaServiceIsStarted(client *gorequest.SuperAgent, kibanaUri string) error {
	response, body, err := client.Get(kibanaUri + "/app/kibana").End()

	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	return nil
}

func (kibana *kibanaContainer) Stop() error {
	return kibana.pool.Purge(kibana.resource)
}
