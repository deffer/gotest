package main

import (
	"bufio"
	_ "container/list"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
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
	ext  string

	// initial position of a file in a list
	number int
	// new posotion of file in a list
	newNumber int

	// file name without any numbering.
	stem string
}

func (fi AnyFileInfo) String() string {
	return fmt.Sprintf("%s%s%s %v->%v (%s)", fi.path, fi.name, fi.ext, fi.number, fi.newNumber, fi.stem)
}

// for sorting
type ByNumber []AnyFileInfo

func (b ByNumber) Len() int           { return len(b) }
func (b ByNumber) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByNumber) Less(i, j int) bool { return b[i].number < b[j].number }

const ARG_DEST_DEFAULT = "./"

var numberedFileRegex = regexp.MustCompile(`(.*?)[\.\s\-\_]+(\d+)$`)
var argfrom int // when ordering desctination files, start numbering them with this
var argdest = ARG_DEST_DEFAULT
var argsource string = "c:/dev/docs/idcards"

func init() {

	flag.IntVar(&argfrom, "from", 0, "Start enumeration from this number")
	flag.Bool("emulate", false, "Do not make any file system changes")
	flag.Parse()

	if flag.NArg() > 0 {
		argsource = flag.Args()[0]
	}

	if flag.NArg() > 1 {
		argdest = flag.Args()[1]
	}
}

/**
Filters file from a folder. Reject files that dont end with number.
Used for filtering photoes or music in a folder.
*/
func filterNumbered(filename string) (matches bool, fileinfo AnyFileInfo) {
	if numberedFileRegex.MatchString(filename) {
		matches = true
		fileinfo = AnyFileInfo{name: filename}
	} else {
		matches = false
	}
	return
}

/**
Used for analyzing files listed in winamp playlist (or any other file list)
*/
var musicFileRegex = regexp.MustCompile(`\s*(\d+)[\.\s\-\_]+(.*)`)

func analyzeListEntry(filename string) (matches bool, fileinfo AnyFileInfo) {
	matches = true
	fileinfo = AnyFileInfo{}
	fileinfo.path, fileinfo.name = filepath.Split(filename)
	fileinfo.ext = filepath.Ext(fileinfo.name)
	fileinfo.name = fileinfo.name[:len(fileinfo.name)-len(fileinfo.ext)]
	if strings.TrimSpace(fileinfo.name) == "" {
		matches = false
	} else {
		var groups []string = musicFileRegex.FindStringSubmatch(fileinfo.name)
		if len(groups) > 0 {
			fileinfo.stem = groups[2]
		} else {
			groups = numberedFileRegex.FindStringSubmatch(fileinfo.name)
			if len(groups) > 0 {
				fileinfo.stem = groups[1]
			} else {
				fileinfo.stem = fileinfo.name
			}
		}
	}
	return
}

func withFilesInList(listfile string, analyzeFunc func(string) (bool, AnyFileInfo)) (listedFiles []AnyFileInfo, processed int) {
	file, _ := os.Open(listfile)
	scanner := bufio.NewScanner(file)
	var result []AnyFileInfo = make([]AnyFileInfo, 0, 100)
	var i int = 0
	for scanner.Scan() {
		s := scanner.Text()
		if !strings.HasPrefix(s, "#") {
			if !filepath.IsAbs(s) {
				s = joinpath(listfile, s)
			}
			if matches, fileinfo := analyzeFunc(s); matches {
				fileinfo.number = i
				fileinfo.newNumber = i + argfrom
				result = append(result, fileinfo)
				i++
				fmt.Println(fileinfo)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, 0
	}
	return result, i
}

func withFilesInDir(folder string, filterFunc func(string) (bool, AnyFileInfo)) (filesInDir []AnyFileInfo, processed int) {
	processed = 0
	var folderNames []string = make([]string, 0, 100)
	var ignoredNames []string = make([]string, 0, 100)
	var result []AnyFileInfo = make([]AnyFileInfo, 0, 100)
	if dirlist, error := ioutil.ReadDir(folder); error == nil {
		for _, a := range dirlist {
			if !a.IsDir() {
				if matches, fileinfo := filterFunc(a.Name()); matches {
					fileinfo.path = folder
					result = append(result, fileinfo)
					fmt.Printf("Accepted %s\n", fileinfo.name)
				} else {
					ignoredNames = append(ignoredNames, a.Name())
				}
			} else {
				folderNames = append(folderNames, a.Name())
			}
		}

		if len(folderNames) > 0 {
			fmt.Printf("Skipped FOLDERS: \n    %s\n", strings.Join(folderNames, "\n    "))
			fmt.Printf("Ignored files: \n    %s\n", strings.Join(ignoredNames, "\n    "))
		}
	} else {
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
	finfo, err := os.Stat(argsource)
	if err != nil { // no such file or dir
		log.Fatalf("File or folder %s does not exist: %v", argsource, err)
		return
	}
	var listedFiles []AnyFileInfo
	var processed int
	if finfo.IsDir() {
		withFilesInDir(argsource, filterNumbered)
	} else {
		fmt.Printf("Opening file %s, setting source folder to %s\n", argsource, filepath.Dir(argsource))
		listedFiles, processed = withFilesInList(argsource, analyzeListEntry)
		fmt.Printf("Total %v(%v)\n", len(listedFiles), processed)
	}

	sort.Sort(ByNumber(listedFiles))

	/*for lf := range listedFiles {
		fmt.Println(listedFiles[lf].number)
	}*/

	fmt.Printf("Finished!\n")
}
