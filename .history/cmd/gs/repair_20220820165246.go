package gs

import (
	"os"

	"github.com/go-spring/gs/internal"
)

// repair 修复远程项目的链接
func repair(rootDir string) {
	branch := "main"
	if len(os.Args) > 3 {
		branch = os.Args[3]
	}
	_, dir, project := validProject(arg(2))
	internal.Add(rootDir, project, dir, branch)
}
