package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

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

// arg 获取命令行参数
func arg(index int) string {
	if len(os.Args) > index {
		return os.Args[index]
	}
	panic("not enough arg")
}

var springProject = regexp.MustCompile("spring-.*")
var starterProject = regexp.MustCompile("starter-.*")

// validProject 项目名称是否有效
func validProject(project string) (prefix string, _ string) {
	if !springProject.MatchString(project) && !starterProject.MatchString(project) {
		panic("error project " + project)
	}
	return strings.Split(project, "-")[0], project
}

var commands = map[string]func(rootDir string){
	"pull":    cPull,
	"push":    cPush,
	"remove":  cRemove,
	"release": cRelease,
}

func cPull(rootDir string) {
	prefix, project := validProject(arg(2))
	stash(rootDir, func() {
		pull(rootDir, project, prefix)
	})
}

func cPush(rootDir string) {
	prefix, project := validProject(arg(2))
	stash(rootDir, func() {
		push(rootDir, project, prefix)
	})
}

func cRemove(rootDir string) {
	prefix, project := validProject(arg(2))
	remove(rootDir, project, prefix)
}

func cRelease(rootDir string) {
	tag := arg(2)
	release(rootDir, tag)
}

func main() {
	fmt.Println(help)

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			os.Exit(-1)
		}
	}()

	// 获取工作目录
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 执行任何命令之前一定要先备份本地文件
	zip(rootDir)

	command := arg(1)
	if fn, ok := commands[command]; !ok {
		panic("error command " + command)
	} else {
		fn(rootDir)
	}
}

// zip 备份本地文件
func zip(rootDir string) {
	backupDir := path.Dir(rootDir)
	baseName := path.Base(rootDir)
	now := time.Now().Format("20060102150405")
	zipFile := fmt.Sprintf("%s-%s.zip", baseName, now)
	cmd := internal.NewCommand("zip", "-r", zipFile, "./"+baseName)
	if _, err := cmd.Run(backupDir); err != nil {
		panic(err)
	}
}

/////////////////////////////////// git //////////////////////////////////////////

// remotes 返回远程项目列表
func remotes(rootDir string) []string {
	// git remote
	cmd := internal.NewCommand("git", "remote")
	if r, err := cmd.Run(rootDir); err == nil {
		ss := strings.Split(r, "\n")
		return ss[:len(ss)-1]
	} else {
		panic(err)
	}
}

func sync(rootDir, project, prefix string) {
	// git subtree pull --prefix=spring/spring-message spring-message master
	cmd := internal.NewCommand("git", "subtree", "pull", fmt.Sprintf("--prefix=%s/%s", prefix, project), project, "master")
	if _, err := cmd.Run(rootDir); err != nil {
		panic(err)
	}
}

// pull 新增并拉取项目
func pull(rootDir, project, prefix string) {

	if remotes := remotes(rootDir); internal.ContainsString(remotes, project) < 0 {

		add := false

		defer func() {
			if !add {
				remove(rootDir, project, prefix)
			}
		}()

		// git remote add -f spring-message https://github.com/go-spring/spring-message.git
		cmd := internal.NewCommand("git", "remote", "add", "-f", project, fmt.Sprintf("https://github.com/go-spring/%s.git", project))
		if _, err := cmd.Run(rootDir); err != nil {
			panic(err)
		}

		// git subtree add --prefix=spring/spring-message spring-message master
		cmd = internal.NewCommand("git", "subtree", "add", fmt.Sprintf("--prefix=%s/%s", prefix, project), project, "master")
		if _, err := cmd.Run(rootDir); err != nil {
			panic(err)
		}

		add = true
	}

	sync(rootDir, project, prefix)
}

// push 推送并同步项目
func push(rootDir, project, prefix string) {

	// git subtree push --prefix=spring/spring-message spring-message master
	cmd := internal.NewCommand("git", "subtree", "push", fmt.Sprintf("--prefix=%s/%s", prefix, project), project, "master")
	if _, err := cmd.Run(rootDir); err != nil {
		panic(err)
	}

	sync(rootDir, project, prefix)
}

// remove 删除项目
func remove(rootDir, project, prefix string) {

	cmd := internal.NewCommand("git", "remote", "remove", project)
	_, _ = cmd.Run(rootDir)

	log.Println(remotes(rootDir))

	projectDir := path.Join(rootDir, prefix, project)
	_ = os.RemoveAll(projectDir)

	cmd = internal.NewCommand("ls", projectDir)
	_, _ = cmd.Run(rootDir)
}

// release 发布版本
func release(rootDir, tag string) {

}

func stash(rootDir string, fn func()) {
	stash := stashIn(rootDir)
	defer func() {
		if stash {
			stashOut(rootDir)
		}
	}()
	fn()
}

func stashIn(rootDir string) bool {
	// git stash
	cmd := internal.NewCommand("git", "stash")
	if _, err := cmd.Run(rootDir); err != nil {
		return false
	}
	return true
}

func stashOut(rootDir string) {
	// git stash pop stash@{0}
	cmd := internal.NewCommand("git", "stash", "pop", "stash@{0}")
	_, _ = cmd.Run(rootDir)
}
