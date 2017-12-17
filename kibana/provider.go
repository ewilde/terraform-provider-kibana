package kibana

import (
	"github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
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
		},

		DataSourcesMap: map[string]*schema.Resource{
			"kibana_index": dataSourceKibanaIndex(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"kibana_search": resourceKibanaSearch(),
		},

		ConfigureFunc: providerConfigure,
	}
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

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &kibana.Config{
		KibanaBaseUri: d.Get("kibana_uri").(string),
		KibanaType:    kibana.ParseKibanaType(d.Get("kibana_type").(string)),
		KibanaVersion: d.Get("kibana_version").(string),
	}

	return kibana.NewClient(config), nil
}
