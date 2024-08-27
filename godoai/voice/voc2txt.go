package voice

import (
	"fmt"
	"godoai/libs"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func Voc2txt(req ReqUpload) (string, error) {
	scriptPath, err := libs.GetCmdPath("voice", "non-streaming-decode-files")
	if err != nil {
		return "", err
	}
	log.Printf("scriptPath: %s", scriptPath)
	hfModelDir, err := libs.GetHfModelDir()
	if err != nil {
		return "", err
	}
	modelDir := filepath.Join(hfModelDir, req.Model)
	tokenDir := filepath.Join(modelDir, req.Params.Token)
	if !libs.PathExists(tokenDir) {
		return "", fmt.Errorf("token file not found")
	}
	cmdParams := []string{
		"--num-threads", strconv.Itoa(runtime.NumCPU()),
		"--tokens", tokenDir,
	}
	params := req.Params
	switch params.Type {
	case "paraformer":
		modelPath := filepath.Join(modelDir, params.Model)
		if !libs.PathExists(modelPath) {
			return "", fmt.Errorf("model file not found")
		}
		cmdParams = append(cmdParams,
			"--paraformer", modelPath,
			"--model-type", "paraformer",
			"--debug", "0",
		)
	case "whisper":
		encodePath := filepath.Join(modelDir, params.Encoder)
		decodePath := filepath.Join(modelDir, params.Decoder)
		if !libs.PathExists(encodePath) || !libs.PathExists(decodePath) {
			return "", fmt.Errorf("model file not found")
		}
		cmdParams = append(cmdParams,
			"--whisper-encoder", encodePath,
			"--whisper-decoder", decodePath,
		)
	case "zipformer":
		encodePath := filepath.Join(modelDir, params.Encoder)
		decodePath := filepath.Join(modelDir, params.Decoder)
		joinerPath := filepath.Join(modelDir, params.Joiner)
		if !libs.PathExists(encodePath) || !libs.PathExists(decodePath) || !libs.PathExists(joinerPath) {
			return "", fmt.Errorf("models file not found")
		}
		cmdParams = append(cmdParams,
			"--encoder", encodePath,
			"--decoder", decodePath,
			"--joiner", joinerPath,
			"--model-type", "transducer",
			"--debug", "0",
		)

	case "nemo":
		modelPath := filepath.Join(modelDir, params.Model)
		if !libs.PathExists(modelPath) {
			return "", fmt.Errorf("model file not found")
		}
		cmdParams = append(cmdParams,
			"--nemo-ctc", modelPath,
			"--model-type", "nemo_ctc",
			"--debug", "0",
		)
	case "telespeech":
		modelPath := filepath.Join(modelDir, params.Model)
		if !libs.PathExists(modelPath) {
			return "", fmt.Errorf("model file not found")
		}
		cmdParams = append(cmdParams,
			"--telespeech-ctc", modelPath,
			"--model-type", "telespeech_ctc",
			"--debug", "0",
		)
	}
	cmdParams = append(cmdParams,
		req.File,
	)
	log.Printf("cmdParams: %v", cmdParams)
	cmd := exec.Command(scriptPath, cmdParams...)

	// 创建管道来读取命令的输出
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create stdout pipe: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Failed to create stderr pipe: %v", err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	// 读取标准输出和标准错误
	stdoutBytes, err := io.ReadAll(stdoutPipe)
	if err != nil {
		log.Fatalf("Failed to read from stdout pipe: %v", err)
	}
	stderrBytes, err := io.ReadAll(stderrPipe)
	if err != nil {
		log.Fatalf("Failed to read from stderr pipe: %v", err)
	}

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Command failed: %v\n%s", err, string(stderrBytes))
	}

	// 输出命令结果
	fmt.Printf("Command output:\n%s", string(stdoutBytes))
	fmt.Printf("Command stderr:\n%s", string(stderrBytes))
	res := ExtactText(string(stderrBytes))
	if res != "" {
		return res, nil
	} else {
		res := ExtactText(string(stdoutBytes))
		if res != "" {
			return res, nil
		}
	}

	return "", nil
}
func ExtactText(outputStr string) string {
	// 分割输出为多行
	lines := strings.Split(outputStr, "\n")

	// 标记是否已经找到 "Decoding done!"
	foundDecodingDone := false

	// 遍历每一行
	for _, line := range lines {
		if strings.Contains(line, "Decoding done!") {
			// 当前行包含 "Decoding done!"，设置标记
			foundDecodingDone = true
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "Wave duration:") {
			// 当前行以 "Wave duration:" 开始，停止循环
			break
		}
		if foundDecodingDone {
			// 如果已经找到了 "Decoding done!" 并且当前行不是 "Wave duration:"，
			// 则这行就是我们要找的文本。
			// 去掉时间戳
			re := regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d+\s`)
			cleanLine := re.ReplaceAllString(line, "")
			return strings.TrimSpace(cleanLine)
		}
	}
	return ""
}
