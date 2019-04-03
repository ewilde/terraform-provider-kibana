package kibana

import (
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"os"
	"strings"
	"testing"
)

type testContext struct {
	containers    []container
	KibanaUri     string
	KibanaIndexId string
}

type container interface {
	Stop() error
}

var authForContainerVersion = map[string]map[KibanaType]AuthenticationHandler{
	"5.5.3": {
		KibanaTypeVanilla: &BasicAuthenticationHandler{"elastic", "changeme"},
		KibanaTypeLogzio:  createLogzAuthenticationHandler(),
	},
	DefaultKibanaVersion6: {KibanaTypeVanilla: &NoAuthenticationHandler{}},
}

func getAuthForContainerVersion(version string, kibanaType KibanaType) AuthenticationHandler {
	handler, ok := authForContainerVersion[version]
	if !ok {
		handler = authForContainerVersion[DefaultKibanaVersion6]
	}

	return handler[kibanaType]
}

func RunTestsWithoutContainers(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func RunTestsWithContainers(m *testing.M, client *KibanaClient) {
	testContext, err := startKibana(GetEnvVarOrDefault("ELK_VERSION", DefaultKibanaVersion6), client)
	if err != nil {
		log.Fatalf("Could not start kibana: %v", err)
	}

	err = os.Setenv(EnvKibanaUri, testContext.KibanaUri)
	if err != nil {
		log.Fatalf("Could not set kibana uri env variable: %v", err)
	}

	err = os.Setenv(EnvKibanaIndexId, testContext.KibanaIndexId)
	if err != nil {
		log.Fatalf("Could not set kibana index id env variable: %v", err)
	}

	code := m.Run()

	if client.Config.KibanaType == KibanaTypeVanilla {
		stopKibana(testContext)
	}

	os.Exit(code)
}

func DefaultTestKibanaClient() *KibanaClient {
	kibanaClient := NewClient(NewDefaultConfig())
	kibanaClient.SetAuth(getAuthForContainerVersion(kibanaClient.Config.KibanaVersion, kibanaClient.Config.KibanaType))
	return kibanaClient
}

func createLogzAuthenticationHandler() *LogzAuthenticationHandler {
	return &LogzAuthenticationHandler{
		Auth0Uri:  "https://logzio.auth0.com",
		LogzUri:   "https://app-eu.logz.io",
		ClientId:  os.Getenv(EnvLogzClientId),
		UserName:  os.Getenv(EnvKibanaUserName),
		Password:  os.Getenv(EnvKibanaPassword),
		MfaSecret: os.Getenv(EnvLogzMfaSecret),
	}
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
