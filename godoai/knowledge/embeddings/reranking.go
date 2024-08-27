package embeddings

/*
// DocumentEmbedding represents a document with its embedding vector and similarity score.
type DocumentEmbedding struct {
	Document   dbtype.AskResponse
	Similarity float32
}

// BySimilarity implements sort.Interface for sorting by similarity score.
type BySimilarity []DocumentEmbedding

func (ds BySimilarity) Len() int           { return len(ds) }
func (ds BySimilarity) Swap(i, j int)      { ds[i], ds[j] = ds[j], ds[i] }
func (ds BySimilarity) Less(i, j int) bool { return ds[i].Similarity > ds[j].Similarity } // Descending order for highest similarity first

// Reranking reorders the documents based on their cosine similarity to the query.
func RerankingByEmbedding(model string, docs []dbtype.AskResponse, query string) ([]dbtype.AskResponse, error) {
	// Assuming queryEmbedding is the precomputed embedding for the query.
	queryEmbedding, err := getEmbeddingForDoc(model, query)
	if err != nil {
		return docs, fmt.Errorf("failed to get query embedding")
	}

	documentEmbeddings := make([]DocumentEmbedding, 0, len(docs))
	for _, doc := range docs {
		if doc.Embedding == nil {
			continue // Skip documents without embeddings
		}

		cosineSim := cosineSimilarity(queryEmbedding, doc.Embedding)
		documentEmbeddings = append(documentEmbeddings, DocumentEmbedding{Document: doc, Similarity: cosineSim})
	}

	sort.Sort(BySimilarity(documentEmbeddings))

	// Reconstruct the slice of AskResponse in the new order.
	sortedDocs := make([]dbtype.AskResponse, len(documentEmbeddings))
	for i, de := range documentEmbeddings {
		sortedDocs[i] = de.Document
	}

	return sortedDocs, nil
}

// cosineSimilarity calculates the cosine similarity between two vectors.
func cosineSimilarity(a, b []float32) float32 {
	aDotB := 0.0
	aMagnitude := 0.0
	bMagnitude := 0.0

	for i := range a {
		aDotB += float64(a[i]) * float64(b[i])
		aMagnitude += float64(a[i]) * float64(a[i])
		bMagnitude += float64(b[i]) * float64(b[i])
	}

	if aMagnitude == 0 || bMagnitude == 0 {
		return 0 // 防止除以零错误
	}

	return float32(aDotB / (math.Sqrt(float64(aMagnitude)) * math.Sqrt(float64(bMagnitude))))
}

// getQueryEmbedding returns the embedding for the given query.
// Replace this with your actual logic to retrieve the embedding from the model.
// func getQueryEmbedding(model string, query string) []float32 {
// 	// Your code here...
// 	return nil
// }
*/
