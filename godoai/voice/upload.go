package voice

import (
	"encoding/json"
	"fmt"
	"godoai/libs"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type ReqParams struct {
	Type     string   `json:"type"`
	Model    string   `json:"model"`
	Decoder  string   `json:"decoder"`
	Encoder  string   `json:"encoder"`
	Token    string   `json:"token"`
	Joiner   string   `json:"joiner"`
	Lexicon  string   `json:"lexicon"`
	RuleFsts []string `json:"ruleFsts"`
}
type ReqUpload struct {
	File   string    `json:"file"`
	Model  string    `json:"model"`
	Params ReqParams `json:"params"`
}
type ResUpload struct {
	Txt string `json:"txt"`
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	// 解析上传的文件
	err := r.ParseMultipartForm(10000 << 20) // 限制最大上传大小为100MB
	if err != nil {
		libs.Error(w, "上传文件过大"+err.Error())
		return
	}

	// 打印所有表单字段
	for key, values := range r.MultipartForm.Value {
		for _, value := range values {
			fmt.Printf("Key: %s, Value: %s\n", key, value)
		}
	}

	file, header, err := r.FormFile("file") // 假设表单字段名为"file"
	if err != nil {
		libs.Error(w, "没有找到文件")
		return
	}

	defer file.Close()
	savePath, err := saveVocieFile(file, header)
	if err != nil {
		libs.Error(w, "保存文件失败"+err.Error())
		return
	}
	// 从表单中获取params字段的值
	paramsStr := r.FormValue("params")

	// 将JSON字符串解码为ReqParams结构体
	var params ReqParams
	if err := json.Unmarshal([]byte(paramsStr), &params); err != nil {
		libs.Error(w, "无法解析参数: "+err.Error())
		return
	}
	req := ReqUpload{
		File:   savePath,
		Model:  r.FormValue("model"),
		Params: params,
	}
	//log.Printf("req: %+v", req)
	res, err := Voc2txt(req)
	defer os.Remove(savePath)
	if err != nil {
		libs.Error(w, "转换失败"+err.Error())
		return
	}
	resJson := ResUpload{
		Txt: res,
	}
	libs.Success(w, resJson, "success")

}
func saveVocieFile(file multipart.File, header *multipart.FileHeader) (string, error) {

	// 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	voicePath, err := libs.GetVoiceDir()
	if err != nil {
		return "", err
	}
	if !libs.PathExists(voicePath) {
		err = os.MkdirAll(voicePath, 0755)
		if err != nil {
			return "", err
		}
	}
	// 生成随机文件名并保留原扩展名
	randomName := libs.GenerateRandomString(10) + filepath.Ext(header.Filename)
	savePath := filepath.Join(voicePath, randomName)

	out, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 将文件内容写入到服务器上的文件
	_, err = out.Write(fileBytes)
	if err != nil {
		return "", err
	}
	return savePath, nil
}
func ServeAudio(w http.ResponseWriter, r *http.Request) {
	// 从 URL 查询参数中获取音频路径
	audioPath := r.URL.Query().Get("path")
	//log.Printf("audioPath: %s", audioPath)
	// 检查音频路径是否为空或无效
	if audioPath == "" {
		libs.Error(w, "Invalid audio path")
		return
	}

	// 确保音频路径是绝对路径
	absAudioPath, err := filepath.Abs(audioPath)
	//log.Printf("absAudioPath: %s", absAudioPath)
	if err != nil {
		libs.Error(w, err.Error())
		return
	}

	// 获取文件的 MIME 类型
	mimeType := mime.TypeByExtension(filepath.Ext(absAudioPath))
	if mimeType == "" {
		mimeType = "application/octet-stream" // 如果无法识别，就用默认的二进制流类型
	}

	// 设置响应头的 MIME 类型
	w.Header().Set("Content-Type", mimeType)

	// 打开文件并读取内容
	file, err := os.Open(absAudioPath)
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	defer file.Close()

	// 将文件内容写入响应体
	_, err = io.Copy(w, file)
	if err != nil {
		libs.Error(w, err.Error())
	}
}
