package kibana

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

const EnvElasticSearchPath = "ELASTIC_SEARCH_PATH"
const EnvKibanaUri = "KIBANA_URI"
const EnvKibanaUserName = "KIBANA_USERNAME"
const EnvKibanaPassword = "KIBANA_PASSWORD"
const EnvKibanaVersion = "ELK_VERSION"
const EnvKibanaIndexId = "KIBANA_INDEX_ID"
const EnvKibanaType = "KIBANA_TYPE"
const EnvKibanaDebug = "KIBANA_DEBUG"
const EnvLogzClientId = "LOGZ_CLIENT_ID"
const EnvLogzMfaSecret = "LOGZ_MFA_SECRET"
const DefaultKibanaUri = "http://localhost:5601"
const DefaultElasticSearchPath = "/es_admin/.kibana"
const DefaultKibanaVersion6 = "6.0.0"
const DefaultKibanaVersion7 = "7.3.1"
const DefaultLogzioVersion = "6.3.2"
const DefaultKibanaVersion553 = "5.5.3"
const DefaultKibanaVersion = DefaultKibanaVersion6
const DefaultKibanaIndexId = "logstash-*"
const DefaultKibanaIndexIdLogzio = "[logzioCustomerIndex]YYMMDD"

type KibanaType int

var kibanaTypeNames = map[string]KibanaType{
	KibanaTypeVanilla.String(): KibanaTypeVanilla,
	KibanaTypeLogzio.String():  KibanaTypeLogzio,
}

const (
	KibanaTypeUnknown KibanaType = iota
	KibanaTypeVanilla
	KibanaTypeLogzio
)

func ParseKibanaType(value string) KibanaType {
	kibanaType, ok := kibanaTypeNames[value]

	if !ok {
		return KibanaTypeUnknown
	}

	return kibanaType
}

type version string

func (v *version) UnmarshalJSON(data []byte) error {
	var tmp int
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*v = version(strconv.Itoa(tmp))
		return nil
	} else {
		var tmp string
		err = json.Unmarshal(data, &tmp)
		if err != nil {
			return err
		} else {
			*v = version(tmp)
		}
	}
	return nil
}

type Config struct {
	Debug             bool
	DefaultIndexId    string
	ElasticSearchPath string
	KibanaBaseUri     string
	KibanaVersion     string
	KibanaType        KibanaType
	Insecure          bool
}

type KibanaClient struct {
	Config *Config
	client *HttpAgent
}

type createResourceResult553 struct {
	Id      string  `json:"_id"`
	Type    string  `json:"_type"`
	Version version `json:"_version"`
}

var indexClientFromVersion = map[string]func(kibanaClient *KibanaClient) IndexPatternClient{
	DefaultKibanaVersion6: func(kibanaClient *KibanaClient) IndexPatternClient {
		return &IndexPatternClient600{config: kibanaClient.Config, client: kibanaClient.client}
	},
	"5.5.3": func(kibanaClient *KibanaClient) IndexPatternClient {
		return &IndexPatternClient553{config: kibanaClient.Config, client: kibanaClient.client}
	},
}

func getIndexClientFromVersion(version string, kibanaClient *KibanaClient) IndexPatternClient {
	indexClient, ok := indexClientFromVersion[version]
	if !ok {
		indexClient = indexClientFromVersion[DefaultKibanaVersion6]
	}

	return indexClient(kibanaClient)
}

var searchClientFromVersion = map[string]func(kibanaClient *KibanaClient) SearchClient{
	DefaultKibanaVersion6: func(kibanaClient *KibanaClient) SearchClient {
		return &searchClient600{config: kibanaClient.Config, client: kibanaClient.client}
	},
	"5.5.3": func(kibanaClient *KibanaClient) SearchClient {
		return &searchClient553{config: kibanaClient.Config, client: kibanaClient.client}
	},
}

func getSearchClientFromVersion(version string, kibanaClient *KibanaClient) SearchClient {
	searchClient, ok := searchClientFromVersion[version]
	if !ok {
		searchClient = searchClientFromVersion[DefaultKibanaVersion6]
	}

	return searchClient(kibanaClient)
}

var visualizationClientFromVersion = map[string]func(kibanaClient *KibanaClient) VisualizationClient{
	DefaultKibanaVersion6: func(kibanaClient *KibanaClient) VisualizationClient {
		return &visualizationClient600{config: kibanaClient.Config, client: kibanaClient.client}
	},
	"5.5.3": func(kibanaClient *KibanaClient) VisualizationClient {
		return &visualizationClient553{config: kibanaClient.Config, client: kibanaClient.client}
	},
}

func getVisualizationClientFromVersion(version string, kibanaClient *KibanaClient) VisualizationClient {
	visualizationClient, ok := visualizationClientFromVersion[version]
	if !ok {
		visualizationClient = visualizationClientFromVersion[DefaultKibanaVersion6]
	}

	return visualizationClient(kibanaClient)
}

var dashboardClientFromVersion = map[string]func(kibanaClient *KibanaClient) DashboardClient{
	DefaultKibanaVersion6: func(kibanaClient *KibanaClient) DashboardClient {
		return &dashboardClient600{config: kibanaClient.Config, client: kibanaClient.client}
	},
	"5.5.3": func(kibanaClient *KibanaClient) DashboardClient {
		return &dashboardClient553{config: kibanaClient.Config, client: kibanaClient.client}
	},
}

func getDashboardClientFromVersion(version string, kibanaClient *KibanaClient) DashboardClient {
	dashboardClient, ok := dashboardClientFromVersion[version]
	if !ok {
		dashboardClient = dashboardClientFromVersion[DefaultKibanaVersion6]
	}

	return dashboardClient(kibanaClient)
}

var savedObjectsClientFromVersion = map[string]func(kibanaClient *KibanaClient) SavedObjectsClient{
	DefaultKibanaVersion6: func(kibanaClient *KibanaClient) SavedObjectsClient {
		return &savedObjectsClient600{config: kibanaClient.Config, client: kibanaClient.client}
	},
	"5.5.3": func(kibanaClient *KibanaClient) SavedObjectsClient {
		return &savedObjectsClient553{config: kibanaClient.Config, client: kibanaClient.client}
	},
}

func getSavedObjectsClientFromVersion(version string, kibanaClient *KibanaClient) SavedObjectsClient {
	savedObjectsClient, ok := savedObjectsClientFromVersion[version]
	if !ok {
		savedObjectsClient = savedObjectsClientFromVersion[DefaultKibanaVersion6]
	}

	return savedObjectsClient(kibanaClient)
}

func NewDefaultConfig() *Config {
	config := &Config{
		ElasticSearchPath: DefaultElasticSearchPath,
		KibanaBaseUri:     DefaultKibanaUri,
		KibanaVersion:     DefaultKibanaVersion,
		KibanaType:        KibanaTypeVanilla,
		Insecure:          false,
	}

	if value := os.Getenv(EnvElasticSearchPath); value != "" {
		config.ElasticSearchPath = value
	}

	if value := os.Getenv(EnvKibanaUri); value != "" {
		config.KibanaBaseUri = strings.TrimRight(value, "/")
	}

	if value := os.Getenv(EnvKibanaVersion); value != "" {
		config.KibanaVersion = value
	}

	if value := os.Getenv(EnvKibanaType); value != "" {
		config.KibanaType = ParseKibanaType(value)
	}

	if value := os.Getenv(EnvKibanaIndexId); value != "" {
		config.DefaultIndexId = value
	} else {
		if config.KibanaType == KibanaTypeVanilla {
			config.DefaultIndexId = DefaultKibanaIndexId
		} else {
			config.DefaultIndexId = DefaultKibanaIndexIdLogzio
		}
	}

	if value := os.Getenv(EnvKibanaDebug); value != "" {
		config.Debug = true
	}

	return config
}

func NewClient(config *Config) *KibanaClient {
	agent := NewHttpAgent(config, &NoAuthenticationHandler{})
	return &KibanaClient{
		Config: config,
		client: agent,
	}
}

func (kibanaClient *KibanaClient) SetAuth(handler AuthenticationHandler) *KibanaClient {
	kibanaClient.client.authHandler = handler
	return kibanaClient
}

func (kibanaClient *KibanaClient) ChangeAccount(accountId string) error {
	return kibanaClient.client.authHandler.ChangeAccount(accountId, kibanaClient.client)
}

func (kibanaClient *KibanaClient) Search() SearchClient {
	return getSearchClientFromVersion(kibanaClient.Config.KibanaVersion, kibanaClient)
}

func (kibanaClient *KibanaClient) Visualization() VisualizationClient {
	return getVisualizationClientFromVersion(kibanaClient.Config.KibanaVersion, kibanaClient)
}

func (kibanaClient *KibanaClient) Dashboard() DashboardClient {
	return getDashboardClientFromVersion(kibanaClient.Config.KibanaVersion, kibanaClient)
}

func (kibanaClient *KibanaClient) IndexPattern() IndexPatternClient {
	return getIndexClientFromVersion(kibanaClient.Config.KibanaVersion, kibanaClient)
}

func (kibanaClient *KibanaClient) SavedObjects() SavedObjectsClient {
	return getSavedObjectsClientFromVersion(kibanaClient.Config.KibanaVersion, kibanaClient)
}

func (kibanaClient *KibanaClient) SetLogger(logger *log.Logger) *KibanaClient {
	kibanaClient.client.SetLogger(logger)
	return kibanaClient
}

func (config *Config) BuildFullPath(format string, a ...interface{}) string {
	return config.KibanaBaseUri + config.ElasticSearchPath + fmt.Sprintf(format, a...)
}

func addQueryString(currentUrl string, filter interface{}) (string, error) {
	v := reflect.ValueOf(filter)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return currentUrl, nil
	}

	uri, err := url.Parse(currentUrl)
	if err != nil {
		return currentUrl, err
	}

	queryStringValues, err := query.Values(filter)
	if err != nil {
		return currentUrl, err
	}

	uri.RawQuery = queryStringValues.Encode()
	return uri.String(), nil
}
