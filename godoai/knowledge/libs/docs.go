package libs

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"godoai/knowledge/dbtype"
)

type ResContentInfo struct {
	Content string       `json:"content"`
	Images  []ImagesInfo `json:"image"`
}
type ImagesInfo struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func GetWords(content string, contextLength int) ([]string, error) {
	res, err := SplitTokenText(content, contextLength)
	if err != nil {
		return res, err
	}
	return res, nil
}
func GetDocs(files []string, index int, contextLength int) ([]dbtype.DocumentParams, error) {
	docs := []dbtype.DocumentParams{}
	for _, v := range files {
		// 获取文件的基本名称，不包含路径
		baseName := filepath.Base(v)
		// 去除文件扩展名
		fileName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
		savePath := v + "_result.json"
		var reqBodies ResContentInfo
		content, err := os.ReadFile(savePath)
		if err != nil {
			log.Printf("Failed to open file %s", savePath)
			continue
		}
		err = json.Unmarshal(content, &reqBodies)
		if err != nil {
			log.Printf("Failed to read file %s", savePath)
			continue
		}
		if reqBodies.Content != "" {
			splitArr, err := SplitTokenText(reqBodies.Content, contextLength)
			//log.Printf("content: %s", splitArr)
			if err != nil {
				log.Printf("Failed to split file: %s", savePath)
				continue
			}
			if len(splitArr) > 0 {
				for _, article := range splitArr {
					index++
					doc := dbtype.DocumentParams{
						ID:       strconv.Itoa(index),
						Metadata: map[string]string{"category": fileName, "file": v, "type": "text"},
						Content:  article,
					}

					docs = append(docs, doc)
				}
			}
		}

		if len(reqBodies.Images) > 0 {
			for _, image := range reqBodies.Images {
				if image.Content != "" {
					splitArr, err := SplitTokenText(image.Content, contextLength)
					if err != nil {
						log.Printf("Failed to split file image: %s", savePath)
						continue
					}
					if len(splitArr) > 0 {
						for _, article := range splitArr {
							index++
							doc := dbtype.DocumentParams{
								ID:       strconv.Itoa(index),
								Metadata: map[string]string{"category": fileName, "file": image.Path, "type": "image"},
								Content:  article,
							}
							docs = append(docs, doc)
						}
					}
				}

			}
		}

	}
	if len(docs) > 0 {
		docs = ZhTitleEnhance(docs)
	}
	return docs, nil
}
