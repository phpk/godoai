package embeddings

import (
	"context"
	"log"
	"sync"

	"godoai/knowledge/dbtype"
)

type LimitGroup struct {
	limit  int
	group  sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewLimitGroup(maxConcurrency int) *LimitGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return &LimitGroup{
		limit:  maxConcurrency,
		group:  sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
	}
}

func (lg *LimitGroup) AddTask(f func()) {
	lg.group.Add(1)
	go func() {
		defer lg.group.Done()
		select {
		case <-lg.ctx.Done():
			return
		default:
			f()
		}
	}()
}

func (lg *LimitGroup) Wait() {
	lg.group.Wait()
}

func (lg *LimitGroup) Cancel() {
	lg.cancel()
}

func GetEmbeddings(config dbtype.EmbedConifig, model string, docs []string) ([]interface{}, error) {
	embeddings := make([]interface{}, len(docs))

	lg := NewLimitGroup(10) // 设置最大并发数为10
	defer lg.Wait()

	for i, doc := range docs {
		lg.AddTask(func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Error getting embedding for doc %d: %v", i, r)
				}
			}()

			embedding, err := getEmbeddingForDoc(config, model, doc)
			if err != nil {
				return
			}
			embeddings[i] = embedding
		})
	}

	return embeddings, nil
}

func getEmbeddingForDoc(config dbtype.EmbedConifig, model string, text string) ([]float32, error) {
	//var res []float32

	//if libs.SupportedModels[model].Engine == "ollama" {
	// if config.ApiType == "ollama" || config.ApiType == "aiok" {
	// 	//log.Printf("Using llama")
	// 	res, err := NewEmbeddingFuncOllama(config.ApiUrl, model, text)
	// 	if err != nil {
	// 		return res, err
	// 	}
	// 	return res, nil
	// } else {
	// 	return res, fmt.Errorf("model %s not supported", model)
	// }
	res, err := NewEmbeddingFuncOllama(config, model, text)
	if err != nil {
		return res, err
	}
	return res, nil
}
