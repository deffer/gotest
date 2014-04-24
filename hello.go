package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	_ "container/list"
	"strings"
	"log"
)

//import "os"

type NumberedFileInfo interface {
	Name() string
	Ext() string
	Number() int
	NewNumber() int
}

type AnyFileInfo struct {
	path string
	name string
	ext string
	number int
	newNumber int
}

var numberedFileRegex = regexp.MustCompile(`(.*?)(\d+)\.([^\.]+)`)

func filterNumbered(filename string) (matches bool, fileinfo AnyFileInfo){
	if numberedFileRegex.MatchString(filename){
		matches = true
		fileinfo = AnyFileInfo{name: filename}		
	}else{
		matches = false
	}
	return
}

func withFilesInDir(folder string, filterFunc func(string) (bool, AnyFileInfo)) (filesInDir, processed int) {
	filesInDir = 0
	processed = 0
	var folderNames []string = make([]string, 0, 100)
	var ignoredNames []string = make([]string, 0, 100)
	var result []AnyFileInfo = make([]AnyFileInfo, 0, 100)
	if dirlist,error := ioutil.ReadDir(folder); error == nil{
		for _, a:=range dirlist {
			if !a.IsDir() {
				filesInDir++
				if matches, fileinfo := filterFunc(a.Name()); matches{
					result = append(result, fileinfo)
					fmt.Printf("Accepted %s\n",fileinfo.name)
				}else{
					ignoredNames = append(ignoredNames, a.Name())
				}
			}else{
				folderNames = append(folderNames, a.Name())
			}
		}
		
		if len(folderNames) > 0 {
			fmt.Printf("Skipped FOLDERS: \n    %s\n", strings.Join(folderNames, "\n    "))
			fmt.Printf("Ignored files: \n    %s\n", strings.Join(ignoredNames, "\n    "))
		}
	}else{	
		log.Fatalf("Unable to read %s:\n %v \nDetails:\n %#v", folder, error, error)
		//fmt.Printf("Unable to read %s, error code %s", folder, error.Error())
	}	
	return filesInDir, len(result)
	
}

//func init() {} // only called once

func main() {
	fmt.Printf("Starting...\n")	
	withFilesInDir("c:/dev/docs/idcards2", filterNumbered)
	fmt.Printf("Finished!\n")
}