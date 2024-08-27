package convert

import (
	"io"

	"godoai/convert/libs"
)

func ConvertImage(r io.Reader) (string, error) {
	// 获取临时文件的绝对路径
	absFilePath, tmpfile, err := libs.GetTempFile(r, "prefix-image")
	if err != nil {
		return "", err
	}
	paths := []string{absFilePath}
	// 识别文本
	output, err := libs.RunRapid(paths)
	if err != nil {
		return "", err
	}
	libs.CloseTempFile(tmpfile)
	// resultString, err := libs.ExtractText(output)
	// if err != nil {
	// 	return "", err
	// }
	// fmt.Println(resultString)
	return output, nil
}
