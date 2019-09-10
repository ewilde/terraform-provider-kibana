package kibana

import (
	"fmt"
	"strings"
	"testing"

	kibana "github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
)

var testDataSource = map[kibana.KibanaType]string{
	kibana.KibanaTypeVanilla: testAccDataSourceKibanaConfig,
	kibana.KibanaTypeLogzio:  testAccDataSourceKibanaConfigLogzio,
}
var testKibanaIndex = map[kibana.KibanaType]func(string) resource.TestCheckFunc{
	kibana.KibanaTypeVanilla: testAccDataSourceKibanaIndex,
	kibana.KibanaTypeLogzio:  testAccDataSourceKibanaIndexLogz,
}

func TestAccDataSourceKibanaIndex_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSource[testConfig.KibanaType],
				Check: resource.ComposeTestCheckFunc(
					testKibanaIndex[testConfig.KibanaType]("data.kibana_index.basic"),
				),
			},
		},
	})
}

func testAccDataSourceKibanaIndex(dataSource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[dataSource]
		if r == nil {
			return errors.New("expected kibana index data source, but was not present")
		}

		a := r.Primary.Attributes

		expectedTitle := "logstash-*"
		if a["title"] != expectedTitle {
			return fmt.Errorf("expected kibana index title %s actual %s", expectedTitle, a["title"])
		}

		expectedTimeFieldName := "@timestamp"
		if a["time_field_name"] != expectedTimeFieldName {
			return fmt.Errorf("expected kibana index time field name %s actual %s", expectedTimeFieldName, a["time_field_name"])
		}

		if len(r.Primary.ID) <= 0 {
			return fmt.Errorf("expected id to be greater than 0 characters actual length: %d value: %s", len(r.Primary.ID), r.Primary.ID)
		}
		return nil
	}
}

func testAccDataSourceKibanaIndexLogz(dataSource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[dataSource]
		if r == nil {
			return errors.New("expected kibana index data source, but was not present")
		}

		a := r.Primary.Attributes

		expectedTitle := "[logzio"
		if !strings.HasPrefix(a["title"], expectedTitle) {
			return fmt.Errorf("expected kibana index title start with %s actual %s", expectedTitle, a["title"])
		}

		expectedTimeFieldName := "@timestamp"
		if a["time_field_name"] != expectedTimeFieldName {
			return fmt.Errorf("expected kibana index time field name %s actual %s", expectedTimeFieldName, a["time_field_name"])
		}

		if r.Primary.ID != "[logzioCustomerIndex]YYMMDD" {
			return fmt.Errorf("expected id to be [logzioCustomerIndex]YYMMDD characters actual %s", r.Primary.ID)
		}
		return nil
	}
}

const testAccDataSourceKibanaConfig = `
data "kibana_index" "basic" {
	filter {
		name = "title"
		values = ["logstash-*"]
	}
}
`

const testAccDataSourceKibanaConfigLogzio = `
data "kibana_index" "basic" {
	filter {
		name = "id"
		values = ["[logzioCustomerIndex]YYMMDD"]
	}
}
`
