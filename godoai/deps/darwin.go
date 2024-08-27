//go:build darwin

package deps

import (
	_ "embed"
)

//go:embed darwin.zip
var embeddedZip []byte
