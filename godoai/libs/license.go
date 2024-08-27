package libs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type LicenseOsInfo struct {
	OsInfo    UserOsInfo `json:"user_os_info"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
}
type CheckLicenseInfo struct {
	IsTest       bool   `json:"isTest"`
	EncodedInfo  string `json:"osInfo"`
	PublicKeyHex string `json:"publicKey"`
	SignatureHex string `json:"signature"`
}
type ServerRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Token   string `json:"token"`
	Time    int64  `json:"time"`
}

func InitLinese() error {
	lineseInfo := GetStaticLinese()
	//log.Printf("2024/08/14 14:20:49 license info: %+v", lineseInfo)
	if lineseInfo.EncodedInfo == "" || lineseInfo.SignatureHex == "" || lineseInfo.PublicKeyHex == "" {
		//log.Println("EncodedInfo, SignatureHex, or PublicKeyHex is empty, fetching dynamic license...")
		lineseInfo, err := GetGodoLinceInfo()
		if err != nil {
			log.Println("Failed to fetch dynamic license info." + err.Error())
			return err
		}
		//log.Printf("2024/08/14 14:20:49 fetched dynamic license info: %+v", lineseInfo)
		err = SetConfigByName("license", lineseInfo)
		if err != nil {
			return err
		}
	}
	return nil
}
func CheckLinese() bool {
	lineseInfo := GetStaticLinese()
	if lineseInfo.EncodedInfo == "" || lineseInfo.SignatureHex == "" || lineseInfo.PublicKeyHex == "" {
		return false
	}
	return VerifySystem(lineseInfo)
}
func GetGodoLinceInfo() (CheckLicenseInfo, error) {
	osInfo, err := GetSystemInfo()
	if err != nil {
		return CheckLicenseInfo{}, err
	}
	url := "https://godoos.com/license?info=" + osInfo
	log.Printf("url: %+v", url)
	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return CheckLicenseInfo{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return CheckLicenseInfo{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CheckLicenseInfo{}, err
	}

	// 解析响应体为 CheckLicenseInfo 结构体
	res := map[string]interface{}{}
	err = json.Unmarshal(body, &res)
	log.Printf("res: %+v", res)
	if err != nil {
		return CheckLicenseInfo{}, err
	}
	if !res["sucess"].(bool) {
		return CheckLicenseInfo{}, fmt.Errorf("license verification failed")
	}
	checkInfoMap, ok := res["data"].(map[string]interface{})
	if !ok {
		return CheckLicenseInfo{}, fmt.Errorf("invalid data type for 'data'")
	}
	var checkInfo CheckLicenseInfo
	isTestFloat, ok := checkInfoMap["isTest"].(bool)
	if !ok {
		isTestFloat = true
	}
	checkInfo.IsTest = isTestFloat

	checkInfo.EncodedInfo = checkInfoMap["osInfo"].(string)
	checkInfo.PublicKeyHex = checkInfoMap["publicKey"].(string)
	checkInfo.SignatureHex = checkInfoMap["signature"].(string)

	if !VerifySystem(checkInfo) {
		return CheckLicenseInfo{}, fmt.Errorf("license verification failed")
	}
	//log.Printf("====checkInfo: %+v", checkInfo)
	return checkInfo, nil
}
