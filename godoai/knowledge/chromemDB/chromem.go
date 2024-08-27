package chromemDB

import (
	"context"
	"godoai/knowledge/chromemDB/chromem"
	"godoai/knowledge/dbtype"
	"godoai/knowledge/libs"
	"runtime"
)

type DB struct {
	// 存储 ChromemDB 实例
	ChromemDB   *chromem.DB
	BaseURL     string
	EmbedConfig dbtype.EmbedConifig
}

// 初始化 DB，创建 ChromemDB 实例
func NewDB(config dbtype.DbConfig) (*DB, error) {
	db, err := chromem.NewPersistentDB("", false)
	if err != nil {
		return nil, err
	}
	return &DB{ChromemDB: db, BaseURL: config.ApiUrl, EmbedConfig: config.Embedding}, nil
}

func (db *DB) Create(name string, model string) (dbtype.CreateResponse, error) {
	//db := getChromemDb()
	res := dbtype.CreateResponse{
		Name: name,
		Id:   name,
	}
	embed := chromem.NewEmbeddingFuncOllama(model, db.EmbedConfig)
	_, err := db.ChromemDB.GetOrCreateCollection(name, nil, embed)
	if err != nil {
		return res, err
	}
	return res, nil
}
func (db *DB) Add(name string, model string, files []string) error {
	embed := chromem.NewEmbeddingFuncOllama(model, db.EmbedConfig)
	collection, err := db.ChromemDB.GetOrCreateCollection(name, nil, embed)
	if err != nil {
		return err
	}
	index := collection.Count()
	docs, err := libs.GetDocs(files, index, db.EmbedConfig.ContextLength)
	if err != nil {
		return err
	}
	if len(docs) > 0 {
		saveDocs := make([]chromem.Document, len(docs))
		for i, doc := range docs {
			saveDocs[i] = chromem.Document{
				Content:   doc.Content,
				ID:        doc.ID,
				Metadata:  doc.Metadata,
				Embedding: doc.Embedding,
			}
		}
		ctx := context.Background()
		err = collection.AddDocuments(ctx, saveDocs, runtime.NumCPU())
		if err != nil {
			return err
		}
	}
	return nil // 示例，实际应根据情况返回
}

func (db *DB) List() ([]string, error) {
	// 实现列出数据的方法，这取决于 chromem 库的 API
	// 这只是一个示例，可能需要根据实际情况调整
	var res = []string{}
	list := db.ChromemDB.ListCollections()
	for _, collection := range list {
		res = append(res, collection.Name)
	}
	return res, nil
}

func (db *DB) Delete(name string) error {
	return db.ChromemDB.DeleteCollection(name)
}

func (db *DB) Ask(name string, model string, message string) ([]dbtype.AskResponse, error) {
	var res []dbtype.AskResponse
	embed := chromem.NewEmbeddingFuncOllama(model, db.EmbedConfig)
	collection, err := db.ChromemDB.GetOrCreateCollection(name, nil, embed)
	if err != nil {
		return res, err
	}
	ctx := context.Background()
	docRes, err := collection.Query(ctx, message, 10, nil, nil)
	if err != nil {
		return res, err
	}
	if len(docRes) > 0 {
		for _, doc := range docRes {
			res = append(res, dbtype.AskResponse{
				Similarity: doc.Similarity,
				Content:    doc.Content,
				Metadata:   doc.Metadata,
				ID:         doc.ID,
				Embedding:  doc.Embedding,
			})
		}
	}

	return res, nil
}

func (db *DB) DeleteFile(name string, model string, file string) error {
	embed := chromem.NewEmbeddingFuncOllama(model, db.EmbedConfig)
	collection, err := db.ChromemDB.GetOrCreateCollection(name, nil, embed)
	if err != nil {
		return err
	}
	ctx := context.Background()
	if err := collection.Delete(ctx, map[string]string{
		"file": file,
	}, nil); err != nil {
		return err
	}
	return nil
}
