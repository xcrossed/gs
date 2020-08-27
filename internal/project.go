package internal

import (
	"encoding/xml"
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
