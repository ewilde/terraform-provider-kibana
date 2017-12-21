package containers

import (
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"github.com/fsouza/go-dockerclient"
	"bytes"
)

type elasticSearchContainer struct {
	Name     string
	pool     *dockertest.Pool
	resource *dockertest.Resource
	Uri      string
}

func NewElasticSearchContainer(pool *dockertest.Pool) (*elasticSearchContainer, error) {

	envVars := []string{
		"discovery.type=single-node",
	}

	options := &dockertest.RunOptions{
		Name:       "elasticsearch",
		Hostname:   "elasticsearch",
		Repository: "elastic-local",
		Tag:        "6.0.0",
		Env:        envVars,
	}

	resource, err := pool.RunWithOptions(options)
	if err != nil {
		return nil, err
	}

	elasticSearchAddress := fmt.Sprintf("http://localhost:%v", resource.GetPort("9200/tcp"))

	if err := pool.Retry(func() error {
		container, err := pool.Client.InspectContainer(resource.Container.ID)
		if err != nil {
			return err
		}

		if !container.State.Running {
			logOptions := &docker.LogsOptions{Container: container.ID}
			if err := pool.Client.Logs(*logOptions); err != nil {
				errors.New("Container is not running: " + container.State.StateString())
			}

			buf := new(bytes.Buffer)
			logOptions.OutputStream.Write(buf.Bytes())


			return errors.New("Container is not running: " + buf.String())
		}

		client := gorequest.New()
		response, body, errs := client.Get(elasticSearchAddress).End()
		if errs != nil {
			return errs[0]
		}

		if response.StatusCode >= 300 {
			return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to elastic search: %s", err)
	}

	if err != nil {
		log.Fatalf("Could not connect to elastic search: %s", err)
	}

	name := getContainerName(resource)
	log.Printf("Elastic search (%v): up", name)

	return &elasticSearchContainer{
		Name:     name,
		pool:     pool,
		resource: resource,
		Uri:      elasticSearchAddress,
	}, nil
}

func (elasticSearch *elasticSearchContainer) Stop() error {
	return elasticSearch.pool.Purge(elasticSearch.resource)
}
