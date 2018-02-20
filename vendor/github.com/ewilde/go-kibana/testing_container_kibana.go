package kibana

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"time"
)

var imageNameFromVersion = map[string]string{
	DefaultKibanaVersion6: "kibana-oss",
	"5.5.3":               "kibana",
}

type kibanaContainer struct {
	Name     string
	pool     *dockertest.Pool
	resource *dockertest.Resource
	Uri      string
}

func newKibanaContainer(pool *dockertest.Pool, elasticSearch *elasticSearchContainer, kibanaVersion string, client *KibanaClient) (container *kibanaContainer, indexId string, err error) {
	envVars := []string{
		fmt.Sprintf("ELASTICSEARCH_URL=http://%s:9200", elasticSearch.Name),
	}

	imageName, ok := imageNameFromVersion[kibanaVersion]
	if !ok {
		imageName = "kibana-oss"
	}

	options := &dockertest.RunOptions{
		Name:         "kibana",
		Repository:   "docker.elastic.co/kibana/" + imageName,
		Tag:          kibanaVersion,
		Env:          envVars,
		Links:        []string{elasticSearch.Name},
		ExposedPorts: []string{"5601"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5601/tcp": {{HostIP: "", HostPort: "5601"}},
		},
	}

	resource, err := pool.RunWithOptions(options)
	if err != nil {
		return nil, "", err
	}

	kibanaUri := fmt.Sprintf("http://localhost:%v", resource.GetPort("5601/tcp"))

	var indexPatternCreateResult *IndexPatternCreateResult
	pool.MaxWait = time.Minute * 5
	if err := pool.Retry(func() error {

		var err error
		if err = checkKibanaServiceIsStarted(client.client, kibanaUri); err != nil {
			return err
		}

		indexPatternClient := client.IndexPattern()
		indexPatternCreateResult, err = indexPatternClient.Create()
		if err != nil {
			log.Printf("Could not create index pattern:%s\n", err)
			return err
		}

		if error := indexPatternClient.RefreshFields(indexPatternCreateResult.Id); error != nil {
			return error
		}

		if error := indexPatternClient.SetDefault(indexPatternCreateResult.Id); error != nil {
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
	log.Printf("Kibana %s (%v): up\n", kibanaVersion, name)

	return &kibanaContainer{
		Name:     name,
		pool:     pool,
		resource: resource,
		Uri:      kibanaUri,
	}, indexPatternCreateResult.Id, nil
}

func checkKibanaServiceIsStarted(client *HttpAgent, kibanaUri string) error {
	response, body, err := client.Get(kibanaUri + "/app/kibana").End()

	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 400 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	return nil
}

func (kibana *kibanaContainer) Stop() error {
	return kibana.pool.Purge(kibana.resource)
}
