package util

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// var testAccProviders map[string]terraform.ResourceProvider
// var testAccProvider *schema.Provider

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {

	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"util": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("HASHICUPS_USERNAME"); err == "" {
		t.Fatal("HASHICUPS_USERNAME must be set for acceptance tests")
	}
	if err := os.Getenv("HASHICUPS_PASSWORD"); err == "" {
		t.Fatal("HASHICUPS_PASSWORD must be set for acceptance tests")
	}
}
