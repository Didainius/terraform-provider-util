// Package main is the entrypoint for terraform-provider-util
package main

import (
	"github.com/Didainius/terraform-provider-util/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return util.Provider()
		},
	})
}
