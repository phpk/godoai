package sd

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type ResJson struct {
	ImageList []string `json:"image_list"`
	Message   []string `json:"message"`
}

func CreateImage(w http.ResponseWriter, r *http.Request) {
	var req CLIConfig
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Model == "" || req.FileName == "" {
		log.Printf("error get modelpath")
		http.Error(w, "the model or FileName is empty!", http.StatusBadRequest)
		return
	}

	runerFile, err := GetRuner()
	if err != nil {
		log.Printf("error get runer: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//log.Printf("runerFile: %v", runerFile)
	if req.BatchCount == 0 {
		req.BatchCount = 1
	}
	imageList, err := GetRandImgs(req.BatchCount)
	if err != nil {
		log.Printf("error get imageList: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Output = imageList[0]

	params, err := ApplyDefaults(&req)
	if err != nil {
		log.Printf("error apply defaults: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//log.Printf("params: %v", params)
	var outBuff bytes.Buffer
	cmd := exec.Command(runerFile, params...)
	// 重定向标准输出到outBuff
	cmd.Stdout = &outBuff
	if err := cmd.Start(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("Streaming unsupported")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// 从outBuff读取所有输出
				output := outBuff.String()

				if output != "" {
					//log.Printf("Command output:\n%s", output)
					// 将输出发送到客户端
					//w.Write([]byte(output))
					resJson := ResJson{
						ImageList: imageList,
						Message:   strings.Split(output, "\n"),
					}
					json.NewEncoder(w).Encode(resJson)
					flusher.Flush()
					outBuff.Reset()
				}

			case <-done: // 等待命令完成的信号
				return // 收到信号后退出goroutine
			}
		}

	}()
	if err := cmd.Wait(); err != nil {
		log.Printf("Command execution failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		close(done) // 命令执行失败也要关闭done
		return
	}

	close(done)
}
