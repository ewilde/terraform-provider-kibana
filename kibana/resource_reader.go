package kibana

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"reflect"
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
