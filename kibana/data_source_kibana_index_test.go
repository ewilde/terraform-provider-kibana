package kibana

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
)

func TestAccDataSourceKibanaIndex_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKibanaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceKibanaIndex("data.kibana_index.by_title"),
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

		if len(r.Primary.ID) != 36 {
			return fmt.Errorf("expected id to be 36 characters actual length: %d value: %s", len(r.Primary.ID), r.Primary.ID)
		}
		return nil
	}
}

const testAccDataSourceKibanaConfig = `
data "kibana_index" "by_title" {
	filter = {
		name = "title"
		values = ["logstash-*"]
	}
}
`