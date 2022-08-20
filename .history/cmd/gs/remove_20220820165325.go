package gs

import (
	"os"
	"path"

	"github.com/go-spring/gs/internal"
)

// remove 移除远程项目
func remove(rootDir string) {

	_, dir, project := validProject(arg(2))
	internal.Remove(rootDir, project)

	projectDir := path.Join(rootDir, dir)
	_ = os.RemoveAll(projectDir)

	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		panic(err)
	}

	projectXml.Remove(project)
	internal.Remotes(rootDir)
}
