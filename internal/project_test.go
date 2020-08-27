package internal

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestProjectXml(t *testing.T) {

	result := `<?xml version="1.0" encoding="UTF-8"?>
<projects>
  <project>
    <name>spring-core</name>
    <subtree>spring/spring-core</subtree>
    <url>https://github.com/go-spring/spring-core.git</url>
  </project>
  <project>
    <name>spring-boot</name>
    <subtree>spring/spring-boot</subtree>
    <url>https://github.com/go-spring/spring-boot.git</url>
  </project>
</projects>`

	doc1 := ProjectXml{
		XMLName: xml.Name{Space: "", Local: "projects"},
		Projects: []Project{
			{
				Name:    "spring-core",
				Subtree: "spring/spring-core",
				Url:     "https://github.com/go-spring/spring-core.git",
			},
			{
				Name:    "spring-boot",
				Subtree: "spring/spring-boot",
				Url:     "https://github.com/go-spring/spring-boot.git",
			},
		},
	}

	bytes, _ := xml.MarshalIndent(doc1, "", "  ")
	data := xml.Header + string(bytes)
	if data != result {
		t.Errorf("%s != %s", data, result)
	}

	var doc2 ProjectXml
	_ = xml.Unmarshal([]byte(data), &doc2)
	if !reflect.DeepEqual(doc1, doc2) {
		t.Errorf("%v != %v", doc1, doc2)
	}
}
