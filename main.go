package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
)

type Variable struct {
	Name        string  `hcl:"variable,label"`
	Description *string `hcl:"description,attr"`
	Default     *string `hcl:"default,attr"`
}

type Sample struct {
	A string   `hcl:"a"`
	B []string `hcl:"b1,optional"`
}

type Config struct {
	Variables []Variable `hcl:"variable,block"`
	Remain    hcl.Body   `hcl:",remain"`
}

func main() {
	// processtffiles()
	parsehclfile("test.hcl")
}

func parsehclfile(filename string) {
	fmt.Println("Processing .." + filename + " file")
	parser := hclparse.NewParser()
	f, parseDiags := parser.ParseHCLFile("test.hcl")
	if parseDiags.HasErrors() {
		log.Fatal(parseDiags.Error())
	}

	var sampleInstance Sample
	decodeDiags := gohcl.DecodeBody(f.Body, nil, &sampleInstance)
	if decodeDiags.HasErrors() {
		log.Fatal(decodeDiags.Error())
	}

	fmt.Printf("%#v", sampleInstance)
}
func processtffiles() {
	var conf Config
	var vars []Variable

	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	parser := hclparse.NewParser()
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			fmt.Println(info.Name())

			if strings.HasSuffix(info.Name(), ".tf") {
				rel, err := filepath.Rel(root, path)
				if err != nil {
					panic(err)
				}
				f, parseDiags := parser.ParseHCLFile(rel)
				if parseDiags.HasErrors() {
					panic(parseDiags.Error())
				}

				decodeDiags := gohcl.DecodeBody(f.Body, nil, &conf)
				if decodeDiags.HasErrors() {
					panic(decodeDiags.Error())
				}
				vars = append(vars, conf.Variables...)
			}
			return nil
		})
	if err != nil {
		panic(nil)
	}
	fmt.Println("--------------------")
	for i := 0; i < len(vars); i++ {
		fmt.Println(*vars[i].Description)
	}
}
