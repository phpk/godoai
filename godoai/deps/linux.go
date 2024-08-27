//go:build linux

package deps

import (
	_ "embed"
)

//go:embed linux.zip
var embeddedZip []byte
