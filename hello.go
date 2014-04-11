package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	_ "container/list"
	"strings"
)

//import "os"

type NumberedFileInfo interface {
	Name() string
	Ext() string
	Number() int
	NewNumber() int
}


var numberedFileRegex = regexp.MustCompile(`.*?\d+\.[^\.]+`)

func withFilesInDir(folder string) (filesInDir, processed int) {
	filesInDir = 0
	processed = 0
	var folderNames []string = make([]string, 0, 100)
	var ignoredNames []string = make([]string, 0, 100)
	if dirlist,error := ioutil.ReadDir(folder); error == nil{
		for i := 0; i<len(dirlist); i++ {
			if a := dirlist[i]; !a.IsDir() {
				filesInDir++
				if (numberedFileRegex.MatchString(a.Name())){
					processed++
					fmt.Printf("Accepted %s\n",a.Name())
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
		fmt.Printf("Unable to read %s, error code %s", folder, error.Error())
	}	
	return filesInDir, processed
	
}

//func init() {} // only called once

func main() {
	fmt.Printf("Starting...\n")	
	withFilesInDir("c:/dev/docs/idcards")
	fmt.Printf("Finished!\n")
}