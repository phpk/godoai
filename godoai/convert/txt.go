package convert

import (
	"io"
)

func ConvertTxt(r io.Reader) (string, error) {

	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
