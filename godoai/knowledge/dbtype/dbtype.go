package dbtype

type Database interface {
	Create(name string, model string) (CreateResponse, error)
	List() ([]string, error)
	Delete(name string) error
	Add(name string, model string, files []string) error
	Ask(name string, model string, message string) ([]AskResponse, error)
	DeleteFile(name string, model string, file string) error
}
type EmbedConifig struct {
	ApiUrl        string `json:"apiUrl"`
	ApiKey        string `json:"apiKey"`
	ApiType       string `json:"apiType"`
	ContextLength int    `json:"contextLength"`
}
type DbConfig struct {
	Type      string       `json:"type"`
	ApiUrl    string       `json:"apiUrl"`
	ApiKey    string       `json:"apiKey"`
	Embedding EmbedConifig `json:"embedding"`
}
type DbFactory struct {
	DB Database
}
type ConfigParams struct {
	Config DbConfig `json:"config"`
}
type CreateParams struct {
	Name   string   `json:"name"`
	Model  string   `json:"model"`
	Config DbConfig `json:"config"`
}
type DeleteParams struct {
	Name   string   `json:"name"`
	Config DbConfig `json:"config"`
}
type Files struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
type AddParams struct {
	Name   string   `json:"name"`
	Model  string   `json:"model"`
	Files  []string `json:"files"`
	Config DbConfig `json:"config"`
}
type DeleteFilearams struct {
	Name   string   `json:"name"`
	Model  string   `json:"model"`
	File   string   `json:"file"`
	Config DbConfig `json:"config"`
}
type AskParams struct {
	Name    string   `json:"name"`
	Model   string   `json:"model"`
	Message string   `json:"message"`
	Config  DbConfig `json:"config"`
}
type Article struct {
	Text     string `json:"text"`
	File     string `json:"file"`
	Category string `json:"category"`
}
type DocumentParams struct {
	ID        string
	Metadata  map[string]string
	Embedding []float32
	Content   string
}
type AskResponse struct {
	ID         string            `json:"id"`
	Metadata   map[string]string `json:"metadata"`
	Embedding  []float32         `json:"embedding"`
	Content    string            `json:"content"`
	Similarity float32           `json:"similarity"`
}
type CreateResponse struct {
	Name string `json:"name"`
	Id   string `josn:"id"`
}
