package voice

import (
	"fmt"
	"godoai/libs"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func Txt2voc(req TtsRequest) error {
	scriptPath, err := libs.GetCmdPath("voice", "non-streaming-tts")
	if err != nil {
		return err
	}
	hfModelDir, err := libs.GetHfModelDir()
	if err != nil {
		return err
	}
	modelDir := filepath.Join(hfModelDir, req.Model)
	tokenDir := filepath.Join(modelDir, req.Params.Token)
	if !libs.PathExists(tokenDir) {
		return fmt.Errorf("token file not found")
	}
	cmdParams := []string{
		"--num-threads", strconv.Itoa(runtime.NumCPU()),
	}
	params := req.Params
	switch params.Type {
	case "vits":
		modelPath := filepath.Join(modelDir, params.Model)
		if !libs.PathExists(modelPath) {
			return fmt.Errorf("model file not found")
		}
		lexicon := filepath.Join(modelDir, params.Lexicon)
		if !libs.PathExists(lexicon) {
			return fmt.Errorf("lexicon file not found")
		}
		cmdParams = append(cmdParams,
			"--vits-model", modelPath,
			"--vits-lexicon", lexicon,
			"--vits-tokens", tokenDir,
		)
		if len(req.Params.RuleFsts) > 0 {
			arr := []string{}
			for _, v := range req.Params.RuleFsts {
				fstsPath := filepath.Join(modelDir, v)
				if !libs.PathExists(fstsPath) {
					return fmt.Errorf("fsts file not found:" + fstsPath)
				}
				arr = append(arr, filepath.Join(modelDir, v))
			}
			cmdParams = append(cmdParams, "--tts-rule-fsts", strings.Join(arr, ","))
		}
	}
	cmdParams = append(cmdParams,
		"--sid", strconv.Itoa(req.Sid),
		"--debug", "1",
		"--output-filename", req.Path,
		req.Text,
	)
	cmd := exec.Command(scriptPath, cmdParams...)
	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}
	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	return nil
}
