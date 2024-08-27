//go:build windows

package deps

import (
	_ "embed"
)

//go:embed windows.zip
var embeddedZip []byte
