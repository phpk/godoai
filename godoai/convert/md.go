package convert

import (
	"io"
	"regexp"
	"strings"
)

func ConvertMd(r io.Reader) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`<[^>]*>`)
	content := re.ReplaceAllString(string(b), "")
	reMarkdown := regexp.MustCompile(`(\*{1,4}|_{1,4}|\#{1,6})`)
	content = reMarkdown.ReplaceAllString(content, "")
	// 移除换行符
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, "\n", "")

	// 移除多余的空格
	content = strings.TrimSpace(content)
	return content, nil
}
