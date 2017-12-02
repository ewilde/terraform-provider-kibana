package kibana

import (
	"fmt"
	"github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"log"
)

func dataSourceKibanaIndex() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKibanaIndexRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"time_field_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fields": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKibanaIndexRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*kibana.KibanaClient)

	log.Printf("[INFO] Reading kibana indexes")

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return errors.New("No filter provided")
	}

	result, err := client.SavedObjects().GetByType(
		kibana.NewSavedObjectRequestBuilder().
			WithFields([]string{"title", "timeFieldName", "fields"}).
			WithType("index-pattern").
			WithPerPage(100).
			Build())

	if err != nil {
		return err
	}

	matchingObject := &kibana.SavedObject{}
	for _, savedObject := range result.SavedObjects {
		if !matchesFilter(savedObject, filters.(*schema.Set)) {
			continue
		}

		matchingObject = savedObject
		break
	}

	if matchingObject == nil {
		return fmt.Errorf("unable to locate a saved index matching the provided filter: %s", filters)
	}

	d.SetId(matchingObject.Id)
	d.Set("id", matchingObject.Id)
	d.Set("time_field_name", matchingObject.Attributes["timeFieldName"])
	d.Set("title", matchingObject.Attributes["title"])
	d.Set("fields", matchingObject.Attributes["fields"])

	return nil
}

func matchesFilter(savedObject *kibana.SavedObject, filters *schema.Set) bool {
	for _, filterList := range filters.List() {
		filterMap := filterList.(map[string]interface{})
		passed := false

		for _, matchOnValue := range filterMap["values"].([]interface{}) {
			switch filterMap["name"].(string) {
			default:
				if savedObject.Attributes["title"] == matchOnValue {
					passed = true
				}
			}
		}

		if passed {
			continue
		} else {
			return false
		}

	}
	return true
}
