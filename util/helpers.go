package util

import (
	"flag"
	"os"
)

// getTerraformStdout returns std out to write message in terraform output
func getTerraformStdout() *os.File {
	// Needed to avoid errors when uintptr(4) is used
	if v := flag.Lookup("test.v"); v == nil || v.Value.String() != "true" {
		return os.NewFile(uintptr(4), "stdout")
	} else {
		return os.Stdout
	}
}
