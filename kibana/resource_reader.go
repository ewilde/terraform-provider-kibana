package kibana

import (
	"fmt"
	"reflect"

	"github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func readArrayFromResource(d *schema.ResourceData, key string) []string {

	if attr, ok := d.GetOk(key); ok {
		var array []string
		items := attr.([]interface{})
		for _, x := range items {
			item := x.(string)
			array = append(array, item)
		}

		return array
	}

	return nil
}

func readStringFromResource(d *schema.ResourceData, key string) string {
	if attr, ok := d.GetOk(key); ok {
		return attr.(string)
	}
	return ""
}

func readBoolFromResource(d *schema.ResourceData, key string) bool {
	if attr, ok := d.GetOk(key); ok {
		return attr.(bool)
	}
	return false
}

func readIntFromResource(d *schema.ResourceData, key string) int {
	if attr, ok := d.GetOk(key); ok {
		return attr.(int)
	}
	return 0
}

func readMapFromResource(d *schema.ResourceData, key string) map[string]interface{} {

	if attr, ok := d.GetOk(key); ok {
		result := attr.(map[string]interface{})

		for _, value := range result {
			t := reflect.TypeOf(value)
			fmt.Printf("type is %s", t)
		}

		return result
	}

	return nil
}

func readSearchReferencesFromResource(d *schema.ResourceData) []*kibana.SearchReferences {
	return readSearchReferencesFromInterface(d.Get("references"))
}

func readSearchReferencesFromInterface(val interface{}) []*kibana.SearchReferences {
	var searchRefs []*kibana.SearchReferences
	references := val.(*schema.Set).List()

	for _, reference := range references {
		ref := reference.(map[string]interface{})

		searchRef := &kibana.SearchReferences{
			Id:   ref["id"].(string),
			Name: ref["name"].(string),
		}

		switch ref["type"].(string) {
		case kibana.SearchReferencesTypeIndexPattern.String():
			searchRef.Type = kibana.SearchReferencesTypeIndexPattern
		}

		searchRefs = append(searchRefs, searchRef)
	}

	return searchRefs
}

func readDashboardReferencesFromResource(d *schema.ResourceData) []*kibana.DashboardReferences {
	return readDashboardReferencesFromInterface(d.Get("references"))
}

func readDashboardReferencesFromInterface(val interface{}) []*kibana.DashboardReferences {
	var dashboardRefs []*kibana.DashboardReferences
	references := val.(*schema.Set).List()

	for _, reference := range references {
		ref := reference.(map[string]interface{})

		dashboardRef := &kibana.DashboardReferences{
			Id:   ref["id"].(string),
			Name: ref["name"].(string),
		}

		switch ref["type"].(string) {
		case kibana.DashboardReferencesTypeSearch.String():
			dashboardRef.Type = kibana.DashboardReferencesTypeSearch
		case kibana.DashboardReferencesTypeVisualization.String():
			dashboardRef.Type = kibana.DashboardReferencesTypeVisualization
		case kibana.DashboardReferencesTypeIndexPattern.String():
			dashboardRef.Type = kibana.DashboardReferencesTypeIndexPattern
		case kibana.DashboardReferencesTypeTag.String():
			dashboardRef.Type = kibana.DashboardReferencesTypeTag
		case kibana.DashboardReferencesTypeLens.String():
			dashboardRef.Type = kibana.DashboardReferencesTypeLens
		}

		dashboardRefs = append(dashboardRefs, dashboardRef)
	}

	return dashboardRefs
}

func readVisualizationReferencesFromResource(d *schema.ResourceData) []*kibana.VisualizationReferences {
	return readVisualizationReferencesFromInterface(d.Get("references"))
}

func readVisualizationReferencesFromInterface(val interface{}) []*kibana.VisualizationReferences {
	var visRefs []*kibana.VisualizationReferences
	references := val.(*schema.Set).List()

	for _, reference := range references {
		ref := reference.(map[string]interface{})

		visRef := &kibana.VisualizationReferences{
			Id:   ref["id"].(string),
			Name: ref["name"].(string),
		}

		switch ref["type"].(string) {
		case kibana.VisualizationReferencesTypeSearch.String():
			visRef.Type = kibana.VisualizationReferencesTypeSearch
		case kibana.VisualizationReferencesTypeIndexPattern.String():
			visRef.Type = kibana.VisualizationReferencesTypeIndexPattern
		}

		visRefs = append(visRefs, visRef)
	}

	return visRefs
}
