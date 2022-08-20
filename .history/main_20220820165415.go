package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/go-spring/gs/cmd/gs"
	"github.com/go-spring/gs/internal"
)

func test() {
	fmt.Println(help)
	defer func() { fmt.Println() }()

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			os.Exit(-1)
		}
	}()

	command := arg(1)
	cmd, ok := commands[command]
	if !ok {
		panic("error command " + command)
	}

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

	fmt.Print(os.Args, " 输入 Yes 执行该命令: ")
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if strings.TrimSpace(input) != "Yes" {
		os.Exit(-1)
	}

	// 备份本地文件
	if cmd.backup {
		internal.Zip(rootDir)
	}

	// 执行命令
	if cmd.fn != nil {
		cmd.fn(rootDir)
	}
}

func main() {
	gs.Execute()
}
