package voice

import (
	"encoding/base64"
	"encoding/json"
	"godoai/libs"
	"net/http"
	"os"
	"path/filepath"
)

type TtsRequest struct {
	Model  string    `json:"model"`
	Text   string    `json:"text"`
	Path   string    `json:"path"`
	Sid    int       `json:"sid"`
	Params ReqParams `json:"params"`
}
type TtsReponse struct {
	Txt string `json:"txt"`
}

func TtsHandler(w http.ResponseWriter, r *http.Request) {
	var req TtsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		libs.Error(w, err.Error())
		return
	}
	voicePath, err := libs.GetVoiceDir()
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	// 生成随机文件名并保留原扩展名
	randomName := libs.GenerateRandomString(10) + ".wav"
	req.Path = filepath.Join(voicePath, randomName)
	if err := Txt2voc(req); err != nil {
		libs.Error(w, err.Error())
		return
	}
	txt, err := WavToBase64(req.Path)
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	res := ResUpload{
		Txt: txt,
	}
	libs.Success(w, res, "success")

}
func WavToBase64(filename string) (string, error) {
	// 读取文件内容
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// 将文件内容编码为Base64字符串
	base64Str := base64.StdEncoding.EncodeToString(data)

	return base64Str, nil
}
