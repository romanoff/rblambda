package main

import (
	"errors"
	"strings"
	"regexp"
)

func byteEquals(s string, i int, content []byte) bool {
	buffer := ""
	for j := i; j < len(content); j++ {
		buffer += string(content[j])
		if len(buffer) == len(s) {
			break
		}
	}
	return s == buffer
}

type LambdaToken struct {
	OldSyntax bool
	DoBlock   bool
	Params    []string
	Content   []byte
}

var blockVariablesRegexp *regexp.Regexp = regexp.MustCompile("(?s)\\s*\\|([\\w\\,\\s]+)\\|(.*)")

func (self *LambdaToken) Parse(i int, content []byte) (int, error) {
	lambda := false
	openclause := false
	buffer := ""
	for j := i; j < len(content); j++ {
		s := string(content[j])
		if s == "(" || s == "{" || (s == "d" && j+1 < len(content) && string(content[j+1]) == "o") {
			i = j
			break
		}
		buffer += s
	}
	if strings.TrimSpace(buffer) == "->" || strings.TrimSpace(buffer) == "lambda" {
		lambda = true
		if strings.TrimSpace(buffer) == "lambda" {
			self.OldSyntax = true
		}
	}
	if !lambda {
		return 0, errors.New("Could not find lambda syntax")
	}
	// Reading parameters -> (a,b) {}
	if string(content[i]) == "(" {
		buffer = ""
		for j := i + 1; j < len(content); j++ {
			s := string(content[j])
			if s == ")" {
				i = j + 1
				if buffer == "" {
					self.Params = []string{}
				} else {
					self.Params = strings.Split(buffer, ",")
				}
				break
			}
			buffer += s
		}
	}
	// Reading lambda block open clause
	for j := i; j < len(content); j++ {
		s := string(content[j])
		if s == "{" || (s == "d" && j+1 < len(content) && string(content[j+1]) == "o") {
			i = j + 1
			if s == "d" {
				self.DoBlock = true
				i += 1
			}
			openclause = true
			break
		}
	}
	if !openclause {
		return 0, errors.New("Could not find lambda open clause")
	}

	clauseCounter := 1
	buffer = ""
	for j := i; j < len(content); j++ {
		if !self.DoBlock {
			if byteEquals("{", j, content) {
				clauseCounter++
			}
			if byteEquals("}", j, content) {
				clauseCounter--
			}
		} else {
			if byteEquals("do", j, content) {
				clauseCounter++
				j += 1
				buffer += "d"
			}
			if byteEquals("begin", j, content) {
				clauseCounter++
				j += 4
				buffer += "begi"
			}
			// TODO: check for if and end
			// if byteEquals("if", j, content) {
			// 	clauseCounter++
			// 	i += 1
			// 	buffer += "i"
			// }
			if byteEquals("end", j, content) {
				clauseCounter--
				j += 2
				if clauseCounter != 0 {
					buffer += "en"
				}
			}
		}
		if clauseCounter == 0 {
			i = j + 1
			break
		}
		buffer += string(content[j])
	}
	if clauseCounter != 0 {
		return 0, errors.New("Could not find lambda close clause")
	}
	matches := blockVariablesRegexp.FindStringSubmatch(buffer)
	if len(matches) == 3 {
		self.Params = strings.Split(matches[1], ",")
		buffer = matches[2]
	}
	self.Content = []byte(buffer)
	return i, nil
}

func (self *LambdaToken) OpenBlockToken() []byte {
	if self.OldSyntax && self.DoBlock {
		return []byte("do")
	}
	return []byte("{")
}

func (self *LambdaToken) CloseBlockToken() []byte {
	if self.OldSyntax && self.DoBlock {
		return []byte("end")
	}
	return []byte("}")
}

func (self *LambdaToken) GetContent() []byte {
	content := []byte{}
	if self.OldSyntax {
		content = append(content, []byte("lambda ")...)
		content = append(content, self.OpenBlockToken()...)
		if len(self.Params) > 0 {
			content = append(content, []byte(" |"+strings.Join(self.Params, ", ")+"|")...)
		}
	} else {
		content = append(content, []byte("->")...)
		if len(self.Params) > 0 {
			content = append(content, []byte(" ("+strings.Join(self.Params, ", ")+") ")...)
		}
		content = append(content, self.OpenBlockToken()...)
	}
	content = append(content, self.Content...)
	content = append(content, self.CloseBlockToken()...)
	return content
}
