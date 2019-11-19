package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	glooV0 "github.com/leosunmo/gloo-vs-upgrader/internal/gloov0"
	glooV1 "github.com/leosunmo/gloo-vs-upgrader/internal/gloov1"
	"gopkg.in/yaml.v3"
)

func readV0VSYaml(file string) (glooV0.VirtualService, error) {
	var vs glooV0.VirtualService
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return vs, fmt.Errorf("reading %s failed, err: %v ", file, err)
	}

	err = yaml.Unmarshal(yamlFile, &vs)
	if err != nil {
		return vs, fmt.Errorf("unmarshalling yaml failed: %v", err)
	}

	return vs, nil
}

func writeV1Gloo(sourceFile string, v1VirtualService glooV1.VirtualService, inplace bool) error {
	var err error
	var fileName string
	if l := strings.LastIndex(sourceFile, ".yaml"); l != -1 {
		if !inplace {
			if strings.Contains(sourceFile, "-vs.yaml") {
				fileName = strings.Replace(sourceFile, "-vs.yaml", "-vs2.yaml", 1)
			} else {
				fileName = sourceFile[:l] + "-vs2" + sourceFile[l:]
			}
		} else {
			fileName = sourceFile
		}
	} else {
		fileName = sourceFile + ".yaml"
	}

	f, fErr := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if fErr != nil {
		return err
	}
	defer f.Close()
	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	err = encoder.Encode(v1VirtualService)
	if err != nil {
		return err
	}

	fmt.Printf("writing new VirtualService as %s\n", fileName)
	return nil
}

func main() {

	kSvc := flag.Bool("k", false, "Convert Upstreams to KubeSvc routes")
	inplace := flag.Bool("i", false, "Overwrite input file in place")
	flag.Parse()
	var files []string
	if len(flag.Args()) > 0 {
		files = flag.Args()
	} else {
		fmt.Println(len(flag.Args()))
		fmt.Printf("Usage: %s [-k] [-i] <one or more YAML files>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	// Read Gloo VS yaml
	errored := false
	for _, file := range files {
		v0Vs, err := readV0VSYaml(file)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		v1Vs, err := glooV1.ConvertVirtualService(v0Vs, *kSvc)
		if err != nil {
			errored = true
			log.Printf("err while processing %s, %s\n", file, err.Error())
			continue
		}
		fErr := writeV1Gloo(file, v1Vs, *inplace)
		if fErr != nil {
			log.Printf("failed to write %s, %s\n", file, fErr.Error())
		}
	}
	if errored {
		os.Exit(1)
	}
}
