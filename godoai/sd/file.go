package sd

import (
	"fmt"
	"godoai/libs"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func GetRandImgs(num int) ([]string, error) {
	var seedList []string

	imagePath, err := GetImageDir()
	if err != nil {
		return seedList, err
	}
	//根据requestBody.Num 生成num个随机数切片
	nowNum := strconv.FormatInt(time.Now().UnixNano(), 10)
	for i := 1; i <= num; i++ {
		filename := fmt.Sprintf("txt2img_%s.png", nowNum)
		if i > 1 {
			filename = fmt.Sprintf("txt2img_%s_%d.png", nowNum, i)
		}
		numPath := filepath.Join(imagePath, filename)
		seedList = append(seedList, numPath)
	}

	return seedList, nil
}
func GetOutputFiles(num int) ([]string, error) {
	if num < 1 {
		return nil, fmt.Errorf("num must be at least 1")
	}
	prefix := "output"
	var filenames []string
	for i := 1; i <= num; i++ {
		suffix := ""
		if i > 1 {
			suffix = fmt.Sprintf("_%d", i)
		}
		filename := fmt.Sprintf("%s%s.png", prefix, suffix)
		tmpfile, err := os.CreateTemp("", filename)
		if err != nil {
			// If any creation fails, clean up and return the error.
			for _, file := range filenames {
				os.Remove(file)
			}
			return nil, err
		}
		defer tmpfile.Close() // Defer closing until after the loop or error handling.

		absFilePath, _ := filepath.Abs(tmpfile.Name())
		filenames = append(filenames, absFilePath)
	}

	return filenames, nil
}

func GetModelPath(modelPath string, fileName string) (string, error) {
	baseDir, err := libs.GetHfModelDir()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(baseDir, modelPath, fileName)
	if !libs.PathExists(filePath) {
		return "", fmt.Errorf("model not found")
	}
	return filePath, nil
}

func GetRunDir() (string, error) {
	appDir, err := libs.GetRunDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	runDir := filepath.Join(appDir, "sd")
	return runDir, nil
}
func GetImageDir() (string, error) {
	appDir, err := libs.GetUploadDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	savePath := filepath.Join(appDir, "Photo", time.Now().Format("2006-01-02"))
	if !libs.PathExists(savePath) {
		os.MkdirAll(savePath, 0755)
	}
	return savePath, nil
}
