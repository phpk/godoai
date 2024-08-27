package convert

import (
	"bytes"
	"io"

	"godoai/convert/doc"
)

// ConvertDoc converts an MS Word .doc to text.
func ConvertDoc(r io.Reader) (string, error) {

	buf, err := doc.ParseDoc(r)

	if err != nil {
		return "", err
	}

	return buf.(*bytes.Buffer).String(), nil
}
