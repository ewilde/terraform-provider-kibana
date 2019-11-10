package containers

import (
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"os"
)

type TestContext struct {
	containers    []container
	KibanaUri     string
	KibanaIndexId string
}

func StartKibana() (*TestContext, error) {
	log.SetOutput(os.Stdout)

	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	elasticSearch, err := NewElasticSearchContainer(pool)
	if err != nil {
		return nil, err
	}

	kibana, index := NewKibanaContainer(pool, elasticSearch)
	return &TestContext{
		containers:    []container{elasticSearch, kibana},
		KibanaUri:     kibana.Uri,
		KibanaIndexId: index}, nil
}

func StopKibana(testContext *TestContext) {

	for _, container := range testContext.containers {
		err := container.Stop()
		if err != nil {
			log.Printf("Could not stop container: %v \n", err)
		}
	}

}
