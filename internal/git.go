package internal

import (
	"fmt"
	"path"
	"strings"
)

// SafeStash 安全地暂存和恢复用户未提交的修改
func SafeStash(rootDir string, fn func()) {
	stash := stashIn(rootDir)
	defer func() {
		if stash {
			stashOut(rootDir)
		}
	}()
	fn()
}

// stashIn 把未提交的修改暂存起来
func stashIn(rootDir string) bool {
	// git stash
	cmd := NewCommand("git", "stash")
	err := cmd.RunOnConsole(rootDir)
	return err == nil
}

// stashOut 恢复暂存的提交
func stashOut(rootDir string) {
	// git stash pop stash@{0}
	cmd := NewCommand("git", "stash", "pop", "stash@{0}")
	_ = cmd.RunOnConsole(rootDir)
}

// Remotes 返回远程项目列表
func Remotes(rootDir string) []string {
	// git remote
	cmd := NewCommand("git", "remote")
	if r, err := cmd.RunOnBuffer(rootDir); err == nil {
		ss := strings.Split(r, "\n")
		return ss[:len(ss)-1]
	} else {
		panic(err)
	}
}

// Add 添加远程项目，branch 默认为 main
func Add(rootDir, project, dir, branch string) (repository string) {

	// git remote add -f spring-message https://github.com/go-spring/spring-message.git
	repository = fmt.Sprintf("https://github.com/go-spring/%s.git", project)
	cmd := NewCommand("git", "remote", "add", "-f", project, repository)
	if err := cmd.RunOnConsole(rootDir); err != nil {
		panic(err)
	}

	// git subtree add --prefix=spring/spring-message spring-message main
	prefixArg := fmt.Sprintf("--prefix=%s", dir)
	cmd = NewCommand("git", "subtree", "add", prefixArg, project, branch, "--squash")
	if err := cmd.RunOnConsole(rootDir); err != nil {
		panic(err)
	}

	return repository
}

// Remove 删除远程项目
func Remove(rootDir, project string) {
	// git remote remove spring-message
	cmd := NewCommand("git", "remote", "remove", project)
	_ = cmd.RunOnConsole(rootDir)
}

// Sync 同步远程项目，branch 默认为 main
func Sync(rootDir, project, dir, branch string) {
	// git subtree pull --prefix=spring/spring-message spring-message main
	prefixArg := fmt.Sprintf("--prefix=%s", dir)
	cmd := NewCommand("git", "subtree", "pull", prefixArg, project, branch, "--squash")
	if err := cmd.RunOnConsole(rootDir); err != nil {
		panic(err)
	}
}

// Push 推送远程项目，branch 默认为 main
func Push(rootDir, project, dir, branch string) {
	// git subtree push --prefix=spring/spring-message spring-message main
	prefixArg := fmt.Sprintf("--prefix=%s", dir)
	cmd := NewCommand("git", "subtree", "push", prefixArg, project, branch)
	if err := cmd.RunOnConsole(rootDir); err != nil {
		panic(err)
	}
}

// Clone 克隆远程项目
func Clone(rootDir, project, repository string) (dir string) {
	cmd := NewCommand("git", "clone", repository)
	if err := cmd.RunOnConsole(rootDir); err != nil {
		panic(err)
	}
	return path.Join(rootDir, project)
}

// Release 发布远程项目
func Release(rootDir, tag string) {

	cmd := NewCommand("git", "tag", tag)
	if err := cmd.RunOnConsole(rootDir); err != nil {
		panic(err)
	}

	cmd = NewCommand("git", "push", "origin", tag)
	if err := cmd.RunOnConsole(rootDir); err != nil {
		panic(err)
	}
}

// Commit 提交项目修改
func Commit(rootDir, msg string) {
	cmd := NewCommand("git", "commit", "-am", "\""+msg+"\"")
	if err := cmd.RunOnConsole(rootDir); err != nil {
		panic(err)
	}
}
