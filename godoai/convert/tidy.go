package convert

import (
	"bytes"
	"fmt"
	"io"

	"github.com/beevik/etree"
)

// TidyWithEtree 使用beevik/etree库进行简单的XML清理
func Tidy(r io.Reader) ([]byte, error) {
	// 读取并解析XML
	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(r); err != nil {
		return nil, fmt.Errorf("error reading and parsing XML: %w", err)
	}

	// 清理操作：例如，移除空节点
	removeEmptyNodes(doc.Root())

	// 格式化XML
	var buf bytes.Buffer
	if _, err := doc.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("error writing formatted XML: %w", err)
	}

	return buf.Bytes(), nil
}

// removeEmptyNodes 遍历XML树并移除空节点
func removeEmptyNodes(node *etree.Element) {
	for i := len(node.Child) - 1; i >= 0; i-- { // 逆序遍历以安全删除
		token := node.Child[i]
		element, ok := token.(*etree.Element) // 检查是否为etree.Element类型
		if ok {
			text := element.Text() // 获取元素的文本
			if text == "" && len(element.Attr) == 0 && len(element.Child) == 0 {
				node.RemoveChildAt(i)
			} else {
				removeEmptyNodes(element) // 递归处理子节点，传入指针
			}
		}
	}
}
