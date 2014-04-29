package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	_ "container/list"
	"strings"
	"log"
	"runtime"
	"flag"
)

//import "os"

type NumberedFileInfo interface {
	Name() string
	Ext() string
	Number() int
	NewNumber() int
}

type AnyFileInfo struct {
	// these identifies file on a disk
	path string
	name string
	ext string

	// initial position of a file in a list
	number int
	// new posotion of file in a list
	newNumber int

	// file name without any numbering.
	stem string
}

var numberedFileRegex = regexp.MustCompile(`(.*?)(\d+)\.([^\.]+)`)
var flagcmd string
var flagdest string
var source string = "c:/dev/docs/idcards"

func init() {
	flag.StringVar(&flagcmd, "cmd", "copy", "copy|renum")
	flag.StringVar(&flagdest, "dest", "./", "Destination foler")
	flag.Bool("emulate", false, "Do not make any file system changes")
	flag.Parse()
	
	fmt.Printf("Command is %s, args are %s\n", flagcmd, flag.Args())
}


/**
Filters file from a folder. Reject files that dont end with number.
Used for filtering photoes or music in a folder.
*/
func filterNumbered(filename string) (matches bool, fileinfo AnyFileInfo){
	if numberedFileRegex.MatchString(filename){
		matches = true
		fileinfo = AnyFileInfo{name: filename}		
	}else{
		matches = false
	}
	return
}

/**
Accepts any file.
Used for 'filtering' files listed in winamp playlist (or any other file list)
*/
func filterOrdered(filename string)  (matches bool, fileinfo AnyFileInfo){
	matches = true
	fileinfo = AnyFileInfo{name: filename}
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
		_, file, line, _ := runtime.Caller(1)
		log.Printf("Error at %s: %d", file, line)
		log.Fatalf("Unable to read %s:\n %v \nDetails:\n %#v", folder, error, error)
		//fmt.Printf("Unable to read %s, error code %s", folder, error.Error())
	}	
	return filesInDir, len(result)
	
}


func main() {
	fmt.Printf("Starting...\n")	
	withFilesInDir(source, filterNumbered)
	fmt.Printf("Finished!\n")
}