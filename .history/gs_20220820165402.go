package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/go-spring/gs/internal"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

// help 展示命令行用法
const help = `command:
  gs pull spring-*/starter-* [branch]
  gs repair spring-*/starter-* [branch]
  gs push spring-*/starter-*
  gs remove spring-*/starter-*
  gs release tag`

type Command struct {
	backup bool // 是否需要备份
	fn     func(rootDir string)
}

// commands 命令与处理函数的映射
var commands = map[string]Command{
	"pull":    {backup: true, fn: pull},    // 拉取单个远程项目
	"repair":  {backup: false, fn: repair}, // 拉取单个远程项目
	"push":    {backup: true, fn: push},    // 推送单个远程项目
	"remove":  {backup: true, fn: remove},  // 移除单个远程项目
	"release": {backup: true, fn: release}, // 发布所有远程项目
	"backup":  {backup: true, fn: nil},     // 备份本地项目文件
}
