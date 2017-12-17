package kibana

import (
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"os"
	"strings"
)

type testContext struct {
	containers    []container
	KibanaUri     string
	KibanaIndexId string
}

type container interface {
	Stop() error
}

func startKibana(elkVersion string, client *KibanaClient) (*testContext, error) {
	log.SetOutput(os.Stdout)

	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	elasticSearch, err := newElasticSearchContainer(pool, elkVersion)
	if err != nil {
		return nil, err
	}

	kibana, index, err := newKibanaContainer(pool, elasticSearch, elkVersion, client)
	if err != nil {
		return nil, err
	}

	return &testContext{
		containers:    []container{elasticSearch, kibana},
		KibanaUri:     kibana.Uri,
		KibanaIndexId: index}, nil
}

func stopKibana(testContext *testContext) {

	for _, container := range testContext.containers {
		err := container.Stop()
		if err != nil {
			log.Printf("Could not stop container: %v \n", err)
		}
	}

}

func getContainerName(container *dockertest.Resource) string {
	return strings.TrimPrefix(container.Container.Name, "/")
}
