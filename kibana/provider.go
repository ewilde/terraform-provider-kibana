package kibana

import (
	"fmt"
	"log"
	"os"
	"sync"

	kibana "github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

type KibanaProvider struct {
	*schema.Provider
	client *kibana.KibanaClient
	once   sync.Once
}

func Provider() terraform.ResourceProvider {
	kp := &KibanaProvider{
		Provider: &schema.Provider{
			Schema: map[string]*schema.Schema{
				"elastic_search_path": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: envDefaultFuncWithDefault(kibana.EnvElasticSearchPath, kibana.DefaultElasticSearchPath),
					Description: "The elastic search path, defaults to: " + kibana.DefaultElasticSearchPath,
				},
				"kibana_uri": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: envDefaultFuncWithDefault(kibana.EnvKibanaUri, kibana.DefaultKibanaUri),
					Description: "The address of the kibana admin url, defaults to: " + kibana.DefaultKibanaUri,
				},
				"kibana_type": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: envDefaultFuncWithDefault(kibana.EnvKibanaType, kibana.KibanaTypeVanilla.String()),
					Description: "The type of the kibana either vanilla or logz.io, defaults to: " + kibana.KibanaTypeVanilla.String(),
				},
				"kibana_version": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: envDefaultFuncWithDefault(kibana.EnvKibanaVersion, kibana.DefaultKibanaVersion),
					Description: "The version of kibana being terraformed either 6.0.0 or 5.5.3, defaults to: " + kibana.DefaultKibanaVersion,
				},
				"kibana_username": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: envDefaultFuncWithDefault(kibana.EnvKibanaUserName, ""),
					Description: "The username used to connect to kibana",
				},
				"kibana_password": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: envDefaultFuncWithDefault(kibana.EnvKibanaPassword, ""),
					Description: "The password used to connect to kibana",
				},
				"logzio_client_id": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: envDefaultFuncWithDefault(kibana.EnvLogzClientId, ""),
					Description: "The logz.io client id used when connecting to logz.io",
				},
				"logzio_account_id": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: envDefaultFuncWithDefault("LOGZIO_ACCOUNT_ID", ""),
					Description: "The logz.io account id used when connecting to logz.io",
				},
				"logzio_mfa_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: envDefaultFuncWithDefault(kibana.EnvLogzMfaSecret, ""),
					Description: "The logz.io MFA secret if the account has it enabled.",
				},
				"kibana_insecure": {
					Type:        schema.TypeBool,
					Default:     false,
					Optional:    true,
					Description: "Disable SSL verification",
				},
			},

			DataSourcesMap: map[string]*schema.Resource{
				"kibana_index": dataSourceKibanaIndex(),
			},

			ResourcesMap: map[string]*schema.Resource{
				"kibana_search":        resourceKibanaSearch(),
				"kibana_visualization": resourceKibanaVisualization(),
				"kibana_dashboard":     resourceKibanaDashboard(),
				"kibana_role":          resourceKibanaRole(),
			},
		},
	}

	kp.ConfigureFunc = kp.providerConfigure

	return kp
}

func envDefaultFuncWithDefault(key string, defaultValue string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(key); v != "" {
			if v == "true" {
				return true, nil
			} else if v == "false" {
				return false, nil
			}
			return v, nil
		}
		return defaultValue, nil
	}
}

func GetEnvVarOrDefaultBool(key string, defaultValue bool) bool {
	result := os.Getenv(key)

	if result == "" {
		return defaultValue
	}

	if result == "true" || result == "1" {
		return true
	} else {
		return false
	}
}

func (p *KibanaProvider) providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var err error

	p.once.Do(func() {
		config := &kibana.Config{
			ElasticSearchPath: d.Get("elastic_search_path").(string),
			KibanaBaseUri:     d.Get("kibana_uri").(string),
			KibanaType:        kibana.ParseKibanaType(d.Get("kibana_type").(string)),
			KibanaVersion:     d.Get("kibana_version").(string),
			Insecure:          d.Get("kibana_insecure").(bool),
		}

		client := kibana.NewClient(config)
		client.SetAuth(authForContainerVersion[config.KibanaType](config, d))
		client.Config.Debug = GetEnvVarOrDefaultBool("KIBANA_DEBUG", false)

		if accountId, ok := d.GetOk("logzio_account_id"); ok && len(accountId.(string)) > 0 {
			err = client.ChangeAccount(accountId.(string))
			if err != nil {
				return
			}
		}

		p.client = client
	})

	if err != nil {
		return nil, err
	}

	return p.client, nil
}

var authForContainerVersion = map[kibana.KibanaType]func(config *kibana.Config, d *schema.ResourceData) kibana.AuthenticationHandler{
	kibana.KibanaTypeLogzio:  getLogzioAuthHandler,
	kibana.KibanaTypeVanilla: getAuthHandler,
}

func getAuthHandler(config *kibana.Config, d *schema.ResourceData) kibana.AuthenticationHandler {
	userName := ""
	password := ""

	if v, ok := d.GetOk("kibana_username"); ok {
		userName = v.(string)
	}

	if v, ok := d.GetOk("kibana_password"); ok {
		password = v.(string)
	}

	if userName != "" && password != "" {
		return kibana.NewBasicAuthentication(userName, password)
	}

	return &kibana.NoAuthenticationHandler{}
}

func getLogzioAuthHandler(config *kibana.Config, d *schema.ResourceData) kibana.AuthenticationHandler {
	return &kibana.LogzAuthenticationHandler{
		Auth0Uri:  "https://logzio.auth0.com",
		LogzUri:   config.KibanaBaseUri,
		ClientId:  d.Get("logzio_client_id").(string),
		UserName:  d.Get("kibana_username").(string),
		Password:  d.Get("kibana_password").(string),
		MfaSecret: d.Get("logzio_mfa_secret").(string),
	}
}

func handleNotFoundError(err error, d *schema.ResourceData) error {
	if httpError, ok := err.(*kibana.HttpError); ok && httpError.Code == 404 {
		log.Printf("[WARN] Removing %s because it's gone", d.Id())
		d.SetId("")
		return nil
	}

	return fmt.Errorf("error reading: %s: %s", d.Id(), err)
}
