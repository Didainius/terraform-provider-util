package util

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUtilFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("util_file.ubuntu-image", "id"),
				),
			},
		},
	})
}

const configText = `
resource "util_file" "ubuntu-image" {
	remote_address = "http://dl-cdn.alpinelinux.org/alpine/v3.12/releases/x86_64/alpine-netboot-3.12.0-x86_64.tar.gz"
	file_name      = "fetched-item"
}
`
