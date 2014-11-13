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

func forceOldSyntax(content []byte) ([]byte, bool) {
	return content, false
}

func walkRubyFiles(callback func(string, os.FileMode)) {
	filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && filepath.Ext(path) == ".rb" {
			callback(path, f.Mode())
		}
		return nil
	})
}
