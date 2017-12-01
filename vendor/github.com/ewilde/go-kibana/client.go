package kibana

import (
	"github.com/google/go-querystring/query"
	"github.com/parnurzeal/gorequest"
	"net/url"
	"os"
	"reflect"
	"strings"
)

const EnvKibanaUri = "KIBANA_URI"
const EnvKibanaIndexId = "KIBANA_INDEX_ID"
const DefaultKibanaUri = "http://localhost:5601"

type Config struct {
	HostAddress    string
	DefaultIndexId string
}

type KibanaClient struct {
	Config *Config
	client *gorequest.SuperAgent
}

func NewDefaultConfig() *Config {
	config := &Config{
		HostAddress: DefaultKibanaUri,
	}

	if os.Getenv(EnvKibanaUri) != "" {
		config.HostAddress = strings.TrimRight(os.Getenv(EnvKibanaUri), "/")
	}

	if os.Getenv(EnvKibanaIndexId) != "" {
		config.DefaultIndexId = os.Getenv(EnvKibanaIndexId)
	}

	return config
}

func NewClient(config *Config) *KibanaClient {
	return &KibanaClient{
		Config: config,
		client: gorequest.New(),
	}
}

func (kibanaClient *KibanaClient) Search() *SearchClient {
	return &SearchClient{
		config: kibanaClient.Config,
		client: kibanaClient.client,
	}
}

func (kibanaClient *KibanaClient) SavedObjects() *SavedObjectsClient {
	return &SavedObjectsClient{
		config: kibanaClient.Config,
		client: kibanaClient.client,
	}
}

func addQueryString(currentUrl string, filter interface{}) (string, error) {
	v := reflect.ValueOf(filter)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return currentUrl, nil
	}

	url, err := url.Parse(currentUrl)
	if err != nil {
		return currentUrl, err
	}

	queryStringValues, err := query.Values(filter)
	if err != nil {
		return currentUrl, err
	}

	url.RawQuery = queryStringValues.Encode()
	return url.String(), nil
}
