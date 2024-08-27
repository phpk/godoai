package serv

import (
	"encoding/base64"
	"encoding/json"
	"godoai/libs"
	"net/http"
)

func GetLicenseHandle(w http.ResponseWriter, r *http.Request) {
	lineseInfo, err := libs.GetSystemInfo()
	if err != nil {
		libs.Error(w, "failed to generate Linese info.")
		return
	}
	w.WriteHeader(http.StatusOK) // 返回200状态码
	libs.Success(w, lineseInfo, "sucess")
}
func SetLicenseHandler(w http.ResponseWriter, r *http.Request) {
	type LicenseCode struct {
		LicenseCode string `json:"licenseCode"`
	}
	var req LicenseCode
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libs.Error(w, "first Decode request body error")
		return
	}
	// Base64解码
	plaintext, err := base64.StdEncoding.DecodeString(req.LicenseCode)
	if err != nil {
		libs.Error(w, "base64 decode failed")
		return
	}
	// 反序列化JSON为UserOsInfo
	var licenseInfo libs.CheckLicenseInfo
	if err := json.Unmarshal(plaintext, &licenseInfo); err != nil {
		libs.Error(w, "反序列化JSON错误:"+err.Error())
		return
	}
	if !libs.VerifySystem(licenseInfo) {
		libs.Error(w, "licenseCode is invalid")
		return
	}
	err = libs.SetConfigByName("license", licenseInfo)
	if err != nil {
		libs.Error(w, "set licenseCode failed")
		return
	}
	w.WriteHeader(http.StatusOK) // 返回200状态码
	libs.Success(w, nil, "sucess")
}
