package internal

import (
	"fmt"
	"path"
	"time"
)

// zip 备份本地文件
func Zip(rootDir string) {
	backupDir := path.Dir(rootDir)
	baseName := path.Base(rootDir)
	now := time.Now().Format("20060102150405")
	zipFile := fmt.Sprintf("%s-%s.zip", baseName, now)
	cmd := NewCommand("zip", "-qr", "-x=*/vendor/*", zipFile, "./"+baseName)
	if err := cmd.RunOnConsole(backupDir); err != nil {
		panic(err)
	}
}
