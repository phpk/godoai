package convert

import (
	"fmt"
	"io"
	"net/http"

	"jaytaylor.com/html2text"
)

func resErr(err error) Res {
	return Res{
		Status: 201,
		Data:   fmt.Sprintf("error opening file: %v", err),
	}
}
func ConvertHttp(url string) Res {
	resp, err := http.Get(url)
	if err != nil {
		return resErr(err)
	}
	defer resp.Body.Close()

	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return resErr(errRead)
	}
	text, err := html2text.FromString(string(body), html2text.Options{PrettyTables: false})
	if err != nil {
		return resErr(err)
	}
	return Res{
		Status: 0,
		Data:   text,
	}
}
