package internal

import (
	"fmt"
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
	_, err := cmd.Run(rootDir)
	return err == nil
}

// stashOut 恢复暂存的提交
func stashOut(rootDir string) {
	// git stash pop stash@{0}
	cmd := NewCommand("git", "stash", "pop", "stash@{0}")
	_, _ = cmd.Run(rootDir)
}

// Remotes 返回远程项目列表
func Remotes(rootDir string) []string {
	// git remote
	cmd := NewCommand("git", "remote")
	if r, err := cmd.Run(rootDir); err == nil {
		ss := strings.Split(r, "\n")
		return ss[:len(ss)-1]
	} else {
		panic(err)
	}
}

// Add 添加远程项目
func Add(rootDir, project, dir string) (repository string) {

	// git remote add -f spring-message https://github.com/go-spring/spring-message.git
	repository = fmt.Sprintf("https://github.com/go-spring/%s.git", project)
	cmd := NewCommand("git", "remote", "add", "-f", project, repository)
	if _, err := cmd.Run(rootDir); err != nil {
		panic(err)
	}

	// git subtree add --prefix=spring/spring-message spring-message master
	prefixArg := fmt.Sprintf("--prefix=%s", dir)
	cmd = NewCommand("git", "subtree", "add", prefixArg, project, "master", "--squash")
	if _, err := cmd.Run(rootDir); err != nil {
		panic(err)
	}

	return repository
}

// Remove 删除远程项目
func Remove(rootDir, project string) {
	// git remote remove spring-message
	cmd := NewCommand("git", "remote", "remove", project)
	_, _ = cmd.Run(rootDir)
}

// Sync 同步远程项目
func Sync(rootDir, project, dir string) {
	// git subtree pull --prefix=spring/spring-message spring-message master
	prefixArg := fmt.Sprintf("--prefix=%s", dir)
	cmd := NewCommand("git", "subtree", "pull", prefixArg, project, "master", "--squash")
	if _, err := cmd.Run(rootDir); err != nil {
		panic(err)
	}
}

// Push 推送远程项目
func Push(rootDir, project, dir string) {
	// git subtree push --prefix=spring/spring-message spring-message master
	prefixArg := fmt.Sprintf("--prefix=%s", dir)
	cmd := NewCommand("git", "subtree", "push", prefixArg, project, "master")
	if _, err := cmd.Run(rootDir); err != nil {
		panic(err)
	}
}
