// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package lint

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type ImageLintConfig struct {
	IncludeExts       []string               `yaml:"includeExts"`
	IncludeFiles      []string               `yaml:"includeFiles"`
	IncludeLines      []string               `yaml:"includeLines"`
	ExcludeFiles      []string               `yaml:"excludeFiles"`
	SuccessValidators []string               `yaml:"succesValidators"`
	FailureValidators []string               `yaml:"failureValidators"`
	ImageMap          map[string][]ImageLint // consists map as the key and file details as values
}

func New(configFile string) (*ImageLintConfig, error) {
	file, err := ioutil.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	ilc := &ImageLintConfig{}

	err = yaml.Unmarshal([]byte(file), ilc)
	if err != nil {
		return nil, err
	}
	ilc.ImageMap = make(map[string][]ImageLint)
	return ilc, nil
}

func NewFromContent(content []byte) (*ImageLintConfig, error) {
	ilc := &ImageLintConfig{}
	err := yaml.Unmarshal(content, ilc)
	if err != nil {
		return nil, err
	}
	ilc.ImageMap = make(map[string][]ImageLint)
	return ilc, nil
}

type ImageLint struct {
	Path     string
	Status   string
	Position Position
}
type Position struct {
	Row, Col int
}

func (imc *ImageLintConfig) Init(dir string) error {
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			for _, exclude := range imc.ExcludeFiles {
				if strings.HasPrefix(path, exclude) {
					goto jumpOut

				} else if strings.HasPrefix(exclude, "*.") {
					if filepath.Ext(path) == filepath.Ext(exclude) {
						goto jumpOut
					}

				} else if string(exclude[len(exclude)-1]) != "/" { // its a file
					if path == exclude {
						goto jumpOut
					}
				}
			}
			for _, ext := range imc.IncludeExts {
				if ext == filepath.Ext(path) {
					imc.ReadFile(path)
				}
			}

			for _, f := range imc.IncludeFiles {
				if f == path {
					imc.ReadFile(path)
				}
			}
		jumpOut:
			return nil
		})
	if err != nil {
		return err
	}
	return err
}

func (imc *ImageLintConfig) ReadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	count := 1
	skip := false
	//fmt.Println(f)
	for s.Scan() {
		// ignore lines
		// if the line is commented then skip it

		line := strings.Trim(s.Text(), " ")

		// This is for go or programming comments only
		if len(line) >= 2 && line[:2] == "//" {
			continue
		}
		// comments for yaml or yml files. Do not consider that line if line start with a comment
		if len(line) > 1 && line[:1] == "#" {
			continue
		}
		// comments for yaml or yml files. If there is comment in the line take only uncommented part
		index := strings.Index(line, "#")
		if index > 0 {
			line = line[0:index]
		}
		// This is for go or programming code only as comments in yaml files start with #
		// start
		if len(line) >= 2 && line[:2] == "/*" {
			//TODO here
			skip = true
		}
		if strings.Contains(line, "*/") {
			skip = false
		}
		if skip {
			continue
		}

		for _, searchterm := range imc.IncludeLines {
			if strings.Contains(line, searchterm) {
				index := strings.Index(line, searchterm) + len(searchterm)
				if strings.Trim(line[index:], " ") != "" {
					// when to ignore ?
					ln := removeChars(line[index:])
					if canIgnore(ln) {
						continue
					}
					ilints := imc.ImageMap[ln]
					imc.ImageMap[ln] = append(ilints, ImageLint{Path: path, Position: Position{Row: count, Col: index}, Status: "YetToLint"})
				}
			}
		}
		count++
	}
	err = s.Err()
	if err != nil {
		//fmt.Println(err)
		return err
	}
	return nil
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func canIgnore(line string) bool {
	ignores := []string{"%", "$", "%", "{", "}", "...", ",", " "}
	if len(line) < 5 {
		return true
	}
	for _, s := range ignores {
		if strings.Contains(strings.Trim(line, " "), s) {
			return true
		}
	}
	return false
}
func removeChars(line string) string {
	line = strings.Trim(line, " ")
	line = strings.Trim(line, `"`)
	line = strings.Trim(line, `'`)
	return line
}
