package libs

import (
	"bytes"
	"errors"
	"fmt"
	"godoai/libs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunXpdf(pdfPath string) (string, error) {
	tempDir, err := GetTempDir("xpdf-dirs")
	if err != nil {
		return "", err
	}
	tempDirSlash := tempDir
	if !strings.HasSuffix(tempDir, string(filepath.Separator)) { // 检查路径是否已经以分隔符结尾
		tempDirSlash = tempDir + string(filepath.Separator) // 添加分隔符
	}
	runFile, err := getXpdfDir("pdftopng")
	if err != nil {
		return "", err
	}
	// 构建命令
	cmdArgs := []string{
		runFile,
		"-mono",
		pdfPath,
		tempDirSlash,
	}
	// 打印将要执行的命令行
	cmdStr := strings.Join(cmdArgs, " ")
	log.Printf("Executing command: %s\n", cmdStr)
	// 使用Command构造命令
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	//var out bytes.Buffer
	var stderr bytes.Buffer
	//cmd.Stdout = &out    // 捕获标准输出
	cmd.Stderr = &stderr // 捕获标准错误
	// 执行命令
	err = cmd.Run()
	if err != nil {
		// 打印错误信息
		log.Printf("执行命令时出错: %v, stderr: %s", err, stderr.String())
		return "", err
	}
	// 输出命令结果
	// outputStr := out.String()
	// log.Printf("Output command: %s\n", outputStr)
	err = GetImages(pdfPath)
	if err != nil {
		log.Println("Failed to get images:", err)
		return "", err
	}

	dir, err := os.ReadDir(tempDir)
	if err != nil {
		log.Println("Failed to read directory:", err)
		return "", err
	}
	imagePaths := []string{}
	for _, entry := range dir {
		absPath := filepath.Join(tempDir, entry.Name())
		//log.Println(absPath)
		imagePaths = append(imagePaths, absPath)
	}
	//log.Printf("imagePaths: %v\n", imagePaths)
	if len(imagePaths) < 1 {
		return "", errors.New("no images found")
	}
	text, err := RunRapid(imagePaths)
	if err != nil {
		log.Println("Failed to run rapid:", err)
		return "", err
	}

	defer func() {

		if err := os.RemoveAll(tempDir); err != nil {
			log.Printf("Error removing temp dir: %s", err)
		}
	}()
	// go func(pdfPath string) {

	// }(pdfPath)
	return text, nil
}
func GetImages(pdfPath string) error {
	cacheDir, err := libs.GetCacheDir()
	if err != nil {
		return fmt.Errorf("failed to create temporary cache directory: %w", err)
	}
	tempDirSlash := cacheDir
	if !strings.HasSuffix(cacheDir, string(filepath.Separator)) { // 检查路径是否已经以分隔符结尾
		tempDirSlash = cacheDir + string(filepath.Separator) // 添加分隔符
	}
	log.Printf("tempDirSlash: %s\n", tempDirSlash)
	runFile, err := getXpdfDir("pdfimages")
	if err != nil {
		return err
	}
	cmdArgs := []string{
		runFile,
		"-j",
		pdfPath,
		tempDirSlash,
	}
	log.Printf("Executing command: %s\n", strings.Join(cmdArgs, " "))
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	var stderr bytes.Buffer
	//cmd.Stdout = &out    // 捕获标准输出
	cmd.Stderr = &stderr // 捕获标准错误
	if err := cmd.Run(); err != nil {
		log.Printf("执行命令时出错: %v, stderr: %s", err, stderr.String())
		return fmt.Errorf("failed to run pdfimages: %w", err)
	}
	return nil
}
