package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/go-spring/gs/internal"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

// help 展示命令行用法
const help = `command:
  gs pull spring-*/starter-*
  gs push spring-*/starter-*
  gs remove spring-*/starter-*
  gs release tag`

var springProject = regexp.MustCompile("spring-.*")
var starterProject = regexp.MustCompile("starter-.*")

// validProject 项目名称是否有效
func validProject(project string) (prefix string, _ string) {
	if !springProject.MatchString(project) && !starterProject.MatchString(project) {
		panic("error project " + project)
	}
	return strings.Split(project, "-")[0], project
}

// commands 命令与处理函数的映射
var commands = map[string]func(rootDir string){
	"pull":    pull,    // 拉取单个远程项目
	"push":    push,    // 推送单个远程项目
	"remove":  remove,  // 移除单个远程项目
	"release": release, // 发布所有远程项目
}

// arg 获取命令行参数
func arg(index int) string {
	if len(os.Args) > index {
		return os.Args[index]
	}
	panic("not enough arg")
}

func main() {
	fmt.Println(help)

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			os.Exit(-1)
		}
	}()

	command := arg(1)
	fn, ok := commands[command]
	if !ok {
		panic("error command " + command)
	}

	// 获取工作目录
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 备份本地文件
	internal.Zip(rootDir)

	// 执行命令
	fn(rootDir)
}

// pull 拉取远程项目
func pull(rootDir string) {

	prefix, project := validProject(arg(2))
	internal.SafeStash(rootDir, func() {

		remotes := internal.Remotes(rootDir)
		if internal.ContainsString(remotes, project) < 0 {
			add := false
			defer func() {
				if !add {
					remove(rootDir)
				}
			}()
			internal.Add(rootDir, project, prefix)
			add = true
		}

		internal.Sync(rootDir, project, prefix)
	})
}

// push 推送远程项目
func push(rootDir string) {

	prefix, project := validProject(arg(2))
	internal.SafeStash(rootDir, func() {

		// 将修改提交到远程项目
		internal.Push(rootDir, project, prefix)

		// 由于往回合并提交数翻倍，所以去掉试试看
		// internal.Sync(rootDir, project, prefix)
	})
}

// remove 移除远程项目
func remove(rootDir string) {

	prefix, project := validProject(arg(2))
	internal.Remove(rootDir, project)

	projectDir := path.Join(rootDir, prefix, project)
	_ = os.RemoveAll(projectDir)

	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		panic(err)
	}

	internal.Remotes(rootDir)
}

// release 发布所有远程项目
func release(rootDir string) {
	tag := arg(2)
	fmt.Println(rootDir, tag)
}
