package main

import (
	"fmt"
	"log"	
	"os"
	"io/ioutil"
	"bufio"
	"runtime"
	"flag"
	"path/filepath"
	"strings"
	"regexp"
	_ "container/list"	
)


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
	if flag.NArg()>0{
		source = flag.Args()[0]
	}
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

func withFilesInList(listfile string, filterFunc func(string) (bool, AnyFileInfo)) (listedFiles []AnyFileInfo, processed int) {
	file, _ := os.Open(listfile)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if !strings.HasPrefix(s, "#"){    		
    		if !filepath.IsAbs(s){
    			s = joinpath(listfile, s)    			
    		}
			fmt.Println(s)
    	}
	}

	if err := scanner.Err(); err != nil {
    	log.Fatal(err)
	}
	return nil, 0
}

func withFilesInDir(folder string, filterFunc func(string) (bool, AnyFileInfo)) (filesInDir []AnyFileInfo, processed int) {	
	processed = 0
	var folderNames []string = make([]string, 0, 100)
	var ignoredNames []string = make([]string, 0, 100)
	var result []AnyFileInfo = make([]AnyFileInfo, 0, 100)
	if dirlist,error := ioutil.ReadDir(folder); error == nil{
		for _, a:=range dirlist {
			if !a.IsDir() {				
				if matches, fileinfo := filterFunc(a.Name()); matches{
					result = append(result, fileinfo)
					fileinfo.path = folder
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
	return result, len(result)
	
}


func joinpath(source, target string) string {
    if filepath.IsAbs(target) {
        return target
    }
    return filepath.Join(filepath.Dir(source), target)
}

func main() {
	fmt.Printf("Starting...\n")	
	finfo, err := os.Stat(source)
	if err != nil { // no such file or dir
		log.Fatalf("File or folder %s does not exist: %v", source, err)
    	return
	}

	if finfo.IsDir() {
    	withFilesInDir(source, filterNumbered)    	
	} else {
		fmt.Printf("Opening file %s, setting temp folder to %s\n", source, filepath.Dir(source))
    	withFilesInList(source, filterOrdered)
	}
	
	fmt.Printf("Finished!\n")
}