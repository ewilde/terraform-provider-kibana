package kibana

import (
	"testing"

	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/helper/schema"
	"os"
	"log"
	"github.com/ewilde/go-kibana/containers"
	"github.com/ewilde/go-kibana"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

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

	testContext, err := containers.StartKibana()
	if err != nil {
		log.Fatalf("Could not start kibana: %v", err)
	}

	err = os.Setenv(kibana.EnvKibanaUri, testContext.KibanaUri)
	if err != nil {
		log.Fatalf("Could not set kibana host address env variable: %v", err)
	}

	err = os.Setenv(kibana.EnvKibanaIndexId, testContext.KibanaIndexId)
	if err != nil {
		log.Fatalf("Could not set kibana index id env variable: %v", err)
	}

	code := m.Run()

	containers.StopKibana(testContext)

	os.Exit(code)

}