package gs

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/go-spring/gs/internal"
	"github.com/spf13/cobra"
)

// 配置文件
var projectXml internal.ProjectXml
var springProject = regexp.MustCompile("spring-.*")
var starterProject = regexp.MustCompile("starter-.*")

// validProject 项目名称是否有效，返回项目前缀、项目目录、项目名称
func validProject(project string) (prefix string, dir string, _ string) {
	if !springProject.MatchString(project) && !starterProject.MatchString(project) {
		panic("error project " + project)
	}
	prefix = strings.Split(project, "-")[0]
	return prefix, fmt.Sprintf("%s/%s", prefix, project), project
}

// arg 获取命令行参数
func arg(index int) string {
	if len(os.Args) > index {
		return os.Args[index]
	}
	panic("not enough arg")
}

var rootCmd = &cobra.Command{
	Use:   "gs",
	Short: "gs - a simple CLI to transform and inspect strings",
	Long: `gs is a super fancy CLI (kidding)
   
One can use gs to add or modfiy go spring project from the terminal`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func prepare() {
	// fmt.Println(help)
	// defer func() { fmt.Println() }()

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Println(r)
	// 		os.Exit(-1)
	// 	}
	// }()

	// command := arg(1)
	// cmd, ok := commands[command]
	// if !ok {
	// 	panic("error command " + command)
	// }

	// 获取工作目录
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 加载 project.xml 配置文件
	projectFile := path.Join(rootDir, "project.xml")
	err = projectXml.Read(projectFile)
	if err != nil {
		panic(err)
	}

	count := len(projectXml.Projects)
	defer func() {
		if count != len(projectXml.Projects) {
			// 保存 project.xml 配置文件
			err = projectXml.Save(projectFile)
			if err != nil {
				panic(err)
			}
		}
	}()

	// fmt.Print(os.Args, " 输入 Yes 执行该命令: ")
	// input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	// if strings.TrimSpace(input) != "Yes" {
	// 	os.Exit(-1)
	// }

	// 备份本地文件
	// if cmd.backup {
	// 	internal.Zip(rootDir)
	// }

	// 执行命令
	// if cmd.fn != nil {
	// 	cmd.fn(rootDir)
	// }
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
