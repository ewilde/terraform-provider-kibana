package kibana

import (
	"testing"

	"fmt"
	"strconv"
	"strings"

	kibana "github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testConfig *kibana.Config

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"kibana": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMain(m *testing.M) {
	client := kibana.DefaultTestKibanaClient()
	testConfig = client.Config
	if client.Config.KibanaType == kibana.KibanaTypeVanilla {
		kibana.RunTestsWithContainers(m, client)
	} else {
		kibana.RunTestsWithoutContainers(m)
	}
}

func CheckResourceAttrSet(name, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		is, err := primaryInstanceState(s, name)
		if err != nil {
			return err
		}

		expectedSplit := strings.Split(key, ".")
		matchedValue := ""
		keysAreNotEqual := false

		for k, v := range is.Attributes {
			attrSplit := strings.Split(k, ".")
			if len(attrSplit) == len(expectedSplit) {
				keysAreNotEqual = false
				for index, item := range attrSplit {
					if item == expectedSplit[index] {
						continue
					}

					if _, err := strconv.ParseInt(item, 10, 64); err == nil && len(item) > 5 {
						continue
					}

					keysAreNotEqual = true
					break
				}

				if !keysAreNotEqual {
					matchedValue = v
					break
				}
			}
		}

		if keysAreNotEqual {
			return fmt.Errorf("%s: Attribute '%s' not found", name, key)
		}

		if matchedValue != value {
			return fmt.Errorf(
				"%s: Attribute '%s' expected %#v, got %#v",
				name,
				key,
				value,
				matchedValue)
		}

		return nil
	}
}

func primaryInstanceState(s *terraform.State, name string) (*terraform.InstanceState, error) {
	ms := s.RootModule()
	rs, ok := ms.Resources[name]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", name)
	}

	is := rs.Primary
	if is == nil {
		return nil, fmt.Errorf("No primary instance: %s", name)
	}

	return is, nil
}
