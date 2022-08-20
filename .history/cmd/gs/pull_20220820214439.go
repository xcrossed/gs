package gs

import (
	"os"

	"github.com/go-spring/gs/internal"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:     "pull",
	Aliases: []string{"pl"},
	Short:   "pull remote code",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// i := args[0]
		// res, kind := stringer.Inspect(i, false)

		// pluralS := "s"
		// if res == 1 {
		// 	pluralS = ""
		// }
		// fmt.Printf("'%s' has a %d %s%s.\n", i, res, kind, pluralS)
		// 获取工作目录
		rootDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		pull(rootDir, projectName string)

	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

// pull 拉取远程项目
func pull(rootDir string, projectName string) {

	// _, dir, project := validProject(arg(2))
	_, dir, project := validProject(projectName)
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
