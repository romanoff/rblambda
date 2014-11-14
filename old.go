package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func updateToOldSyntax() {
	walkRubyFiles(func(path string, perm os.FileMode) {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("Error while reading '" + path + "' : " + err.Error())
		}
		content, changed := forceOldSyntax(content)
		if changed {
			fmt.Println(path)
		}
		err = ioutil.WriteFile(path, content, perm)
		if err != nil {
			fmt.Println("Error while writing '" + path + "' : " + err.Error())
		}
	})
}

type Contentable interface {
	GetContent() []byte
}

type Content struct {
	Tokens []Contentable
}

func (self *Content) GetContent() []byte {
	contentBytes := []byte{}
	for _, token := range self.Tokens {
		contentBytes = append(contentBytes, token.GetContent()...)
	}
	return contentBytes
}

type ContentToken struct {
	Content []byte
}

func (self *ContentToken) GetContent() []byte {
	return self.Content
}

func forceOldSyntax(fileContent []byte) ([]byte, bool) {
	content := &Content{Tokens: make([]Contentable, 0)}
	token := &ContentToken{Content: []byte{}}
	content.Tokens = append(content.Tokens, token)
	fileContentLength := len(fileContent)
	changed := false
	for i := 0; i < fileContentLength; i++ {
		if fileContentLength > i+3 {
			twoSymbols := string(fileContent[i]) + string(fileContent[i+1])
			if twoSymbols == "->" || twoSymbols == "la" {
				lt := &LambdaToken{}
				j, err := lt.Parse(i, fileContent)
				if err == nil && lt.OldSyntax == false {
					changed = true
				}
				if err == nil {
					lt.OldSyntax = true
					content.Tokens = append(content.Tokens, lt)
					token = &ContentToken{Content: []byte{}}
					content.Tokens = append(content.Tokens, token)
					i = j
				}
			}
		}
		if i < fileContentLength {
			token.Content = append(token.Content, fileContent[i])
		}
	}
	return content.GetContent(), changed
}

func walkRubyFiles(callback func(string, os.FileMode)) {
	filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && filepath.Ext(path) == ".rb" {
			callback(path, f.Mode())
		}
		return nil
	})
}
