package internal

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

type Project struct {
	Name   string `xml:"name"`
	Dir    string `xml:"dir"`
	Url    string `xml:"url"`
	Branch string `xml:"branch"`
}

type ProjectXml struct {
	XMLName  xml.Name  `xml:"projects"`
	Projects []Project `xml:"project"`
}

// Add 添加一个项目
func (p *ProjectXml) Add(project Project) {
	p.Projects = append(p.Projects, project)
}

// Find 查找一个项目，成功返回 true，失败返回 false
func (p *ProjectXml) Find(project string) (Project, bool) {
	for _, t := range p.Projects {
		if t.Name == project {
			return t, true
		}
	}
	return Project{}, false
}

// Remove 移除一个项目
func (p *ProjectXml) Remove(project string) {
	index := -1
	for i := 0; i < len(p.Projects); i++ {
		if p.Projects[i].Name == project {
			index = i
			break
		}
	}
	if index >= 0 {
		projects := append(p.Projects[:index], p.Projects[index+1:]...)
		fmt.Println(projects)
		p.Projects = projects
	}
}

// Read 读取配置文件
func (p *ProjectXml) Read(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return xml.Unmarshal(bytes, p)
}

type projectSlice []Project

func (s projectSlice) Len() int {
	return len(s)
}

func (s projectSlice) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s projectSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Save 保存配置文件
func (p *ProjectXml) Save(filename string) error {
	sort.Sort(projectSlice(p.Projects))
	bytes, err := xml.MarshalIndent(p, "", "    ")
	if err != nil {
		return err
	}
	bytes = []byte(xml.Header + string(bytes))
	return ioutil.WriteFile(filename, bytes, os.ModePerm)
}
