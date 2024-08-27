package libs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func GetAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".godoos"), nil
}
func GetExeDir() string {
	// 获取当前用户主目录
	homeDir, err := GetAppDir()
	if err != nil {
		return ".godoos"
	}
	return filepath.Join(homeDir, "ai")
}
func GetRunDir() (string, error) {
	exeDir := GetExeDir()
	var osType string
	switch runtime.GOOS {
	case "windows":
		osType = "windows"
	case "darwin": // macOS
		osType = "darwin"
	default: // 包含了Linux和其他未明确列出的系统
		osType = "linux"
	}
	runDir := filepath.Join(exeDir, osType)
	// if !PathExists(runDir) {
	// 	os.MkdirAll(runDir, 0755)
	// }
	return runDir, nil
}
func GetCmdPath(path string, name string) (string, error) {
	// 根据操作系统添加.exe后缀
	binaryExt := ""
	if runtime.GOOS == "windows" {
		binaryExt = ".exe"
	}
	scriptName := name + binaryExt
	exeDir, err := GetRunDir() // 获取执行程序的目录
	if err != nil {
		return "", err
	}
	scriptPath := filepath.Join(exeDir, path, scriptName)
	if !PathExists(scriptPath) {
		return "", fmt.Errorf("script %s not found", scriptPath)
	}
	return scriptPath, nil
}
func InitDataDir() (string, error) {
	appDir, err := GetAppDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	dataDir := filepath.Join(appDir, "aidata")
	if !PathExists(dataDir) {
		os.MkdirAll(dataDir, 0755)
	}
	return dataDir, nil
}
func InitConfig() error {
	exist := ExistConfig("dataDir")
	if !exist {
		dataDir, err := InitDataDir()
		if err != nil {
			return err
		}
		osData := ReqBody{
			Name:  "dataDir",
			Value: dataDir,
		}
		SetConfig(osData)
		info := GenerateSystemInfo()
		osInfo := ReqBody{
			Name:  "osInfo",
			Value: info,
		}
		SetConfig(osInfo)
	}

	return nil
}
func GetDataDir() (string, error) {
	dataDirAny, _ := GetConfig("dataDir")

	dataDir, ok := dataDirAny.(string)
	if !ok {
		return "", fmt.Errorf("unexpected type for dataDir: %T", dataDirAny)
	}

	if dataDir == "" {
		dataDir, _ = InitDataDir()
	}

	return dataDir, nil
}
func GetIpList() []string {
	var res []string
	ipListAny, _ := GetConfig("ipList")
	ipList, ok := ipListAny.([]string)
	if !ok {
		return res
	}
	return ipList
}

// GetStaticLinese 函数
func GetStaticLinese() CheckLicenseInfo {
	var res CheckLicenseInfo
	checkInfoAny, exists := GetConfig("license")
	//log.Printf("checkInfoAny: %v", checkInfoAny)
	if !exists {
		return res
	}
	// 将 checkInfoAny 转换为 JSON 字符串
	jsonBytes, err := json.Marshal(checkInfoAny)
	if err != nil {
		log.Println("Failed to marshal to JSON:", err)
		return res
	}

	// 解析 JSON 字符串到 CheckLicenseInfo 结构体
	err = json.Unmarshal(jsonBytes, &res)
	if err != nil {
		log.Println("Failed to unmarshal from JSON:", err)
		return res
	}

	//log.Printf("checkInfo: %v", res)
	return res
}
func GetHfModelDir() (string, error) {
	dataDir, _ := GetDataDir()
	return filepath.Join(dataDir, "hfmodels"), nil
}
func GetOllamaModelDir() string {
	dataDir, _ := GetDataDir()
	return filepath.Join(dataDir, "models")
}
func GetUploadDir() (string, error) {
	dataDir, _ := GetDataDir()
	return filepath.Join(dataDir, "upload"), nil
}
func GetVoiceDir() (string, error) {
	dataDir, _ := GetDataDir()
	return filepath.Join(dataDir, "voice"), nil
}
func GetCacheDir() (string, error) {
	uploadDir, err := GetUploadDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(uploadDir, "cache")
	if PathExists(cacheDir) {
		err = os.RemoveAll(cacheDir)
		if err != nil {
			return "", err
		}
	}
	err = os.MkdirAll(cacheDir, 0755)
	return cacheDir, err
}

func GetTrueCacheDir() (string, error) {
	uploadDir, err := GetUploadDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(uploadDir, "cache")
	return cacheDir, nil
}

func PathExists(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		//log.Println("文件夹存在")
		return true
	} else if os.IsNotExist(err) {
		//log.Println("文件夹不存在")
		return false
	} else if os.IsExist(err) {
		//log.Println("文件夹存在")
		return true
	} else {
		//log.Println("发生错误:", err)
		return false
	}
}
