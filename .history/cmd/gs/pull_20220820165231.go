package gs

import (
	"os"

	"github.com/go-spring/gs/internal"
)

// pull 拉取远程项目
func pull(rootDir string) {

	_, dir, project := validProject(arg(2))
	internal.SafeStash(rootDir, func() {

		branch := "main"
		if len(os.Args) > 3 {
			branch = os.Args[3]
		}

		remotes := internal.Remotes(rootDir)
		if internal.ContainsString(remotes, project) < 0 {
			add := false
			defer func() {
				if !add {
					remove(rootDir)
				}
			}()
			repository := internal.Add(rootDir, project, dir, branch)
			projectXml.Add(internal.Project{
				Name:   project,
				Dir:    dir,
				Url:    repository,
				Branch: branch,
			})
			add = true
		}

		internal.Sync(rootDir, project, dir, branch)
	})
}
