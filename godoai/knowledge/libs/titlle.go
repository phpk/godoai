package libs

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"godoai/knowledge/dbtype"
)

// underNonAlphaRatio 检查文本中非字母字符的比例是否超过给定阈值
// 不计算空格
func underNonAlphaRatio(text string, threshold float64) bool {
	if len(text) == 0 {
		return false
	}

	alphaCount, totalCount := countAlphaAndTotalChars(text)
	if totalCount == 0 {
		return false
	}

	ratio := float64(alphaCount) / float64(totalCount)
	return ratio < threshold
}

// isPossibleTitle 检查文本是否符合标题的条件
func isPossibleTitle(text string, titleMaxWordLength int, nonAlphaThreshold float64) bool {
	if len(text) == 0 {
		fmt.Println("Not a title. Text is empty.")
		return false
	}

	if !endsWithAcceptablePunctuation(text) {
		return false
	}

	if len(text) > titleMaxWordLength {
		return false
	}

	if underNonAlphaRatio(text, nonAlphaThreshold) {
		return false
	}

	// 添加日志语句，输出 text 的长度
	fmt.Printf("Checking text: '%s' with length %d\n", text, len(text))

	if len(text) < 5 || isNumericOnly(text) {
		fmt.Printf("Not a title. Text is too short or all numeric:\n\n%s\n", text)
		return false
	}

	if hasNoLettersInFirstFive(text) {
		return false
	}

	return true
}

// isNumericOnly 检查字符串是否只包含数字
func isNumericOnly(s string) bool {
	for _, char := range s {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

// Document 类表示一个文档
type Document struct {
	PageContent string
	Metadata    map[string]string
}

func ZhTitleEnhance(docs []dbtype.DocumentParams) []dbtype.DocumentParams {
	// 假设我们希望标题至少有8个字符，且非字母字符比例小于0.2
	nonAlphaThreshold := 0.2
	title := ""
	for _, doc := range docs {
		// 判断当前文档内容是否可能为标题，并据此进行分类标记。
		if isPossibleTitle(doc.Content, 20, nonAlphaThreshold) {
			doc.Metadata["category"] = "cn_Title"
			title = doc.Content
		} else if title != "" {
			// 如果之前已识别到标题，则将当前文档内容与之前识别的标题相关联。
			doc.Content = fmt.Sprintf("下文与(%s)有关。%s", title, doc.Content)
		}
	}
	return docs
}

// countAlphaAndTotalChars 统计文本中字母字符和总字符的数量
func countAlphaAndTotalChars(text string) (int, int) {
	alphaCount := 0
	totalCount := 0
	for _, char := range text {
		if strings.TrimSpace(string(char)) != "" {
			if unicode.IsLetter(char) {
				alphaCount++
			}
			totalCount++
		}
	}
	return alphaCount, totalCount
}

// endsWithAcceptablePunctuation 检查输入文本是否以合适的标点符号结尾。
// 合适的标点符号是指除了字母、数字和空格之外的字符。
// 函数返回一个布尔值，如果文本以标点符号结尾，则返回 true，否则返回 false。
func endsWithAcceptablePunctuation(text string) bool {
	// 正则表达式模式，匹配以非字母、数字和空格字符结尾的字符串
	endsInPunctPattern := "[^\\w\\s]\\z"

	// 编译正则表达式
	endsInPunctRE := regexp.MustCompile(endsInPunctPattern)

	// 检查文本是否以非字母、数字和空格字符结尾
	return !endsInPunctRE.MatchString(text)
}

// hasNoLettersInFirstFive 检查文本前五个字符中是否有字母
func hasNoLettersInFirstFive(text string) bool {
	text5 := text[:5]
	if len(text5) < 5 {
		return false
	}
	alphaInText5 := 0
	for _, char := range text5 {
		if unicode.IsLetter(char) {
			alphaInText5++
		}
	}
	return alphaInText5 == 0
}
