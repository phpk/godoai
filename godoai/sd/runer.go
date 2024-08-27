package sd

import (
	"fmt"
	"godoai/libs"
	"path/filepath"
	"runtime"
)

func GetRuner() (string, error) {
	baseDir, err := GetRunDir()
	if err != nil {
		return "", err
	}

	runerFileName := "sdlibs"
	if runtime.GOOS == "windows" {
		runerFileName = runerFileName + ".exe"
	}
	runerFile := filepath.Join(baseDir, runerFileName)

	if !libs.PathExists(runerFile) {
		// Write the content to the file
		// if err := os.WriteFile(runerFile, embeddedSdlibs, 0755); err != nil {
		// 	return "", fmt.Errorf("failed to write file: %w", err)
		// }
		return "", fmt.Errorf("runer file not found")
	}
	return runerFile, nil

}
