package libs

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"runtime"
	"time"
)

// VerifySignature 使用公钥验证签名
func VerifySignature(data []byte, signatureHex string, publicKeyHex string) bool {
	// 将签名从十六进制字符串转换为字节切片
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		log.Printf("Error decoding signature: %v", err)
		return false
	}

	// 将公钥从十六进制字符串转换为字节切片
	publicKey, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		log.Printf("Error decoding public key: %v", err)
		return false
	}

	// 使用 ed25519.Verify 验证签名
	return ed25519.Verify(ed25519.PublicKey(publicKey), data, signature)
}

// VerifySystem 验证时间戳并校验签名
func VerifySystem(linceInfo CheckLicenseInfo) bool {
	//log.Printf("VerifySystem====%v", linceInfo)
	if linceInfo.EncodedInfo == "" || linceInfo.SignatureHex == "" || linceInfo.PublicKeyHex == "" {
		log.Println("linceInfo is empty")
		return false
	}

	// Base64解码
	plaintext, err := base64.StdEncoding.DecodeString(linceInfo.EncodedInfo)
	if err != nil {
		log.Println("base64 decode failed")
		return false
	}
	// 反序列化JSON为UserOsInfo
	var licenseOsInfo LicenseOsInfo
	if err := json.Unmarshal(plaintext, &licenseOsInfo); err != nil {
		log.Printf("反序列化JSON错误: %v", err)
		return false
	}
	// 验证时间有效性
	currentTime := time.Now()
	if currentTime.Before(licenseOsInfo.StartTime) || currentTime.After(licenseOsInfo.EndTime) {
		log.Println("current time is not in the valid time range")
		return false
	}
	// 校验签名
	if !VerifySignature(plaintext, linceInfo.SignatureHex, linceInfo.PublicKeyHex) {
		log.Println("signature verification failed")
		return false
	}
	systemInfo := licenseOsInfo.OsInfo
	// 获取当前系统MAC地址
	currentMAC, err := getMACAddress()
	if err != nil {
		log.Println("get mac address failed")
		return false
	}

	// 检查解密后的MAC地址是否与当前系统MAC地址一致
	if systemInfo.MAC != currentMAC {
		log.Println("mac address is not consistent")
		return false
	}
	if systemInfo.OS != runtime.GOOS {
		log.Println("os is not consistent")
		return false
	}
	if systemInfo.Arch != runtime.GOARCH {
		log.Println("arch is not consistent")
		return false
	}

	return true
}
