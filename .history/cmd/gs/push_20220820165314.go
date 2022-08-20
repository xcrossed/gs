package gs

import "github.com/go-spring/gs/internal"

// push 推送远程项目
func push(rootDir string) {

	_, dir, project := validProject(arg(2))
	internal.SafeStash(rootDir, func() {

		// 将修改提交到远程项目，不需要往回合并
		if p, ok := projectXml.Find(project); ok {
			internal.Push(rootDir, project, dir, p.Branch)
		}
	})
}
