// Package util holds the main code for terraform-provider-util
package util

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"util_file": resourceFile(),
		},
	}
}
