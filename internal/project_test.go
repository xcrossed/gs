package internal

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func TestProjectXml(t *testing.T) {

	doc := ProjectXml{
		XMLName: xml.Name{Space: "", Local: "projects"},
		Projects: []Project{
			{
				Name:   "spring-core",
				Dir:    "spring/spring-core",
				Url:    "https://github.com/go-spring/spring-core.git",
				Branch: "master",
			},
			{
				Name:   "spring-boot",
				Dir:    "spring/spring-boot",
				Url:    "https://github.com/go-spring/spring-boot.git",
				Branch: "main",
			},
		},
	}

	file, _ := ioutil.TempFile("", "project")
	_ = file.Close()

	if err := doc.Save(file.Name()); err != nil {
		panic(err)
	}

	b1, err := ioutil.ReadFile(file.Name())
	if err != nil {
		panic(err)
	}

	b2, err := ioutil.ReadFile("testdata/result.xml")
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(b1, b2) {
		t.Errorf("expect %s but got %s", string(b2), string(b1))
	}
}
