package internal

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Project struct {
	Name string `xml:"name"`
	Dir  string `xml:"dir"`
	Url  string `xml:"url"`
}

type ProjectXml struct {
	XMLName  xml.Name  `xml:"projects"`
	Projects []Project `xml:"project"`
}

// Add 添加一个项目
func (p *ProjectXml) Add(project Project) {
	p.Projects = append(p.Projects, project)
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

// Save 保存配置文件
func (p *ProjectXml) Save(filename string) error {
	bytes, err := xml.MarshalIndent(p, "", "    ")
	if err != nil {
		return err
	}
	bytes = []byte(xml.Header + string(bytes))
	return ioutil.WriteFile(filename, bytes, os.ModePerm)
}
