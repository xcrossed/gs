package main

// import (
// 	"bufio"
// 	"bytes"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"path"
// 	"path/filepath"
// 	"regexp"
// 	"strings"
// 	"time"

// 	"github.com/go-spring/gs/internal"
// )

// func init() {
// 	log.SetFlags(log.Lshortfile)
// }

// // help 展示命令行用法
// const help = `command:
//   gs pull spring-*/starter-* [branch]
//   gs repair spring-*/starter-* [branch]
//   gs push spring-*/starter-*
//   gs remove spring-*/starter-*
//   gs release tag`

// var springProject = regexp.MustCompile("spring-.*")
// var starterProject = regexp.MustCompile("starter-.*")

// // validProject 项目名称是否有效，返回项目前缀、项目目录、项目名称
// func validProject(project string) (prefix string, dir string, _ string) {
// 	if !springProject.MatchString(project) && !starterProject.MatchString(project) {
// 		panic("error project " + project)
// 	}
// 	prefix = strings.Split(project, "-")[0]
// 	return prefix, fmt.Sprintf("%s/%s", prefix, project), project
// }

// type Command struct {
// 	backup bool // 是否需要备份
// 	fn     func(rootDir string)
// }

// // commands 命令与处理函数的映射
// var commands = map[string]Command{
// 	"pull":    {backup: true, fn: pull},    // 拉取单个远程项目
// 	"repair":  {backup: false, fn: repair}, // 拉取单个远程项目
// 	"push":    {backup: true, fn: push},    // 推送单个远程项目
// 	"remove":  {backup: true, fn: remove},  // 移除单个远程项目
// 	"release": {backup: true, fn: release}, // 发布所有远程项目
// 	"backup":  {backup: true, fn: nil},     // 备份本地项目文件
// }

// // arg 获取命令行参数
// func arg(index int) string {
// 	if len(os.Args) > index {
// 		return os.Args[index]
// 	}
// 	panic("not enough arg")
// }

// // 配置文件
// var projectXml internal.ProjectXml

// func main() {
// 	fmt.Println(help)
// 	defer func() { fmt.Println() }()

// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Println(r)
// 			os.Exit(-1)
// 		}
// 	}()

// 	command := arg(1)
// 	cmd, ok := commands[command]
// 	if !ok {
// 		panic("error command " + command)
// 	}

// 	// 获取工作目录
// 	rootDir, err := os.Getwd()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 加载 project.xml 配置文件
// 	projectFile := path.Join(rootDir, "project.xml")
// 	err = projectXml.Read(projectFile)
// 	if err != nil {
// 		panic(err)
// 	}

// 	count := len(projectXml.Projects)
// 	defer func() {
// 		if count != len(projectXml.Projects) {
// 			// 保存 project.xml 配置文件
// 			err = projectXml.Save(projectFile)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}()

// 	fmt.Print(os.Args, " 输入 Yes 执行该命令: ")
// 	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
// 	if strings.TrimSpace(input) != "Yes" {
// 		os.Exit(-1)
// 	}

// 	// 备份本地文件
// 	if cmd.backup {
// 		internal.Zip(rootDir)
// 	}

// 	// 执行命令
// 	if cmd.fn != nil {
// 		cmd.fn(rootDir)
// 	}
// }

// // pull 拉取远程项目
// func pull(rootDir string) {

// 	_, dir, project := validProject(arg(2))
// 	internal.SafeStash(rootDir, func() {

// 		branch := "main"
// 		if len(os.Args) > 3 {
// 			branch = os.Args[3]
// 		}

// 		remotes := internal.Remotes(rootDir)
// 		if internal.ContainsString(remotes, project) < 0 {
// 			add := false
// 			defer func() {
// 				if !add {
// 					remove(rootDir)
// 				}
// 			}()
// 			repository := internal.Add(rootDir, project, dir, branch)
// 			projectXml.Add(internal.Project{
// 				Name:   project,
// 				Dir:    dir,
// 				Url:    repository,
// 				Branch: branch,
// 			})
// 			add = true
// 		}

// 		internal.Sync(rootDir, project, dir, branch)
// 	})
// }

// // repair 修复远程项目的链接
// func repair(rootDir string) {
// 	branch := "main"
// 	if len(os.Args) > 3 {
// 		branch = os.Args[3]
// 	}
// 	_, dir, project := validProject(arg(2))
// 	internal.Add(rootDir, project, dir, branch)
// }

// // push 推送远程项目
// func push(rootDir string) {

// 	_, dir, project := validProject(arg(2))
// 	internal.SafeStash(rootDir, func() {

// 		// 将修改提交到远程项目，不需要往回合并
// 		if p, ok := projectXml.Find(project); ok {
// 			internal.Push(rootDir, project, dir, p.Branch)
// 		}
// 	})
// }

// // remove 移除远程项目
// func remove(rootDir string) {

// 	_, dir, project := validProject(arg(2))
// 	internal.Remove(rootDir, project)

// 	projectDir := path.Join(rootDir, dir)
// 	_ = os.RemoveAll(projectDir)

// 	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
// 		panic(err)
// 	}

// 	projectXml.Remove(project)
// 	internal.Remotes(rootDir)
// }

// // release 发布所有远程项目
// func release(rootDir string) {
// 	tag := arg(2)
// 	err := filepath.Walk(rootDir, func(walkFile string, _ os.FileInfo, err error) error {

// 		if err != nil {
// 			return err
// 		}

// 		if path.Base(walkFile) != "go.mod" {
// 			return nil
// 		}

// 		fmt.Println(walkFile)
// 		fileData, e0 := ioutil.ReadFile(walkFile)
// 		if e0 != nil {
// 			return nil
// 		}

// 		outBuf := bytes.NewBuffer(nil)
// 		r := bufio.NewReader(strings.NewReader(string(fileData)))
// 		for {
// 			line, isPrefix, e1 := r.ReadLine()
// 			if len(line) > 0 && e1 != nil {
// 				panic(e1)
// 			}
// 			if isPrefix {
// 				panic(errors.New("ReadLine returned prefix"))
// 			}
// 			if e1 != nil {
// 				if e1 != io.EOF {
// 					panic(err)
// 				}
// 				break
// 			}
// 			s := strings.TrimSpace(string(line))
// 			if strings.HasPrefix(s, "github.com/go-spring/spring-") ||
// 				strings.HasPrefix(s, "github.com/go-spring/starter-") {
// 				index := strings.LastIndexByte(s, ' ')
// 				if index <= 0 {
// 					panic(errors.New(s))
// 				}
// 				b := append(line[:index+2], []byte(tag)...)
// 				outBuf.Write(b)
// 			} else {
// 				outBuf.Write(line)
// 			}
// 			outBuf.WriteString("\n")
// 		}

// 		fmt.Println(outBuf.String())
// 		return ioutil.WriteFile(walkFile, outBuf.Bytes(), os.ModePerm)
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	// 提交代码更新
// 	internal.Commit(rootDir, "publish "+tag)

// 	// 遍历所有项目，推送远程更新
// 	for _, project := range projectXml.Projects {
// 		_, dir, _ := validProject(project.Name)
// 		internal.Push(rootDir, project.Name, dir, project.Branch)
// 	}

// 	// 创建临时目录
// 	now := time.Now().Format("20060102150405")
// 	buildDir := path.Join(rootDir, "..", "go-spring-build-"+now)
// 	err = os.MkdirAll(buildDir, os.ModePerm)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 遍历所有项目，推送远程标签
// 	for _, project := range projectXml.Projects {
// 		projectDir := internal.Clone(buildDir, project.Name, project.Url)
// 		internal.Release(projectDir, tag)
// 	}
// }
