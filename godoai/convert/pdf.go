package convert

import (
	"io"

	"godoai/convert/libs"
)

func ConvertPDF(r io.Reader) (string, error) {
	// 获取临时文件的绝对路径
	absFilePath, tmpfile, err := libs.GetTempFile(r, "prefix-pdf")
	if err != nil {
		return "", err
	}
	output, err := libs.RunXpdf(absFilePath)
	if err != nil {
		return "", err
	}
	libs.CloseTempFile(tmpfile)
	return output, nil

}
