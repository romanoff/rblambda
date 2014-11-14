package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func updateToNewSyntax() {
	walkRubyFiles(func(path string, perm os.FileMode) {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("Error while reading '" + path + "' : " + err.Error())
		}
		content, changed := forceNewSyntax(content)
		if changed {
			fmt.Println(path)
		}
		err = ioutil.WriteFile(path, content, perm)
		if err != nil {
			fmt.Println("Error while writing '" + path + "' : " + err.Error())
		}
	})
}

func forceNewSyntax(fileContent []byte) ([]byte, bool) {
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
				if lt.OldSyntax == true {
					changed = true
				}
				if err == nil {
					lt.OldSyntax = false
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
