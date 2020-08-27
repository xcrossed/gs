package internal

import (
	"encoding/xml"
)

type Project struct {
	Name    string `xml:"name"`
	Subtree string `xml:"subtree"`
	Url     string `xml:"url"`
}

type ProjectXml struct {
	XMLName  xml.Name  `xml:"projects"`
	Projects []Project `xml:"project"`
}
