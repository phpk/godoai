package embeddings

import (
	"log"

	"godoai/knowledge/dbtype"
)

// SelectDocuments selects up to 3 records of each type ('text' and 'image') from the given documents.
func SelectDocuments(docs []dbtype.AskResponse) []dbtype.AskResponse {
	textDocs := make([]dbtype.AskResponse, 0)
	imageDocs := make([]dbtype.AskResponse, 0)
	if len(docs) < 1 {
		return docs
	}
	// 分类记录
	for _, doc := range docs {
		docType, ok := doc.Metadata["type"]
		if !ok {
			continue // 忽略没有'type'元数据的文档
		}
		if docType == "text" {
			textDocs = append(textDocs, doc)
		} else if docType == "image" {
			imageDocs = append(imageDocs, doc)
		}
	}

	// 选择文本类型的前5条记录
	selectedTextDocs := textDocs[:min(5, len(textDocs))]

	// 如果文本类型的记录不足5条，用图像类型的记录补充，同时标记这些记录
	borrowedForText := 5 - len(selectedTextDocs)
	if borrowedForText > 0 && len(imageDocs) >= borrowedForText {
		selectedTextDocs = append(selectedTextDocs, imageDocs[:borrowedForText]...)
		// 移除已经被借用的记录
		imageDocs = imageDocs[borrowedForText:]
	}
	selectedImageDocs := imageDocs
	// 确保图像类型的记录不超过3条
	if len(imageDocs) > 3 {
		selectedImageDocs = imageDocs[:3]
	}

	// 合并结果
	result := append(selectedTextDocs, selectedImageDocs...)
	log.Printf("Selected %d documents: %v", len(result), result)
	return result
}

// min returns the smaller of x or y.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
