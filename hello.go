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
	"strconv"
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

func (fi AnyFileInfo) NewFileName(precision int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(precision)+"d %s%s", fi.newNumber, fi.stem, fi.ext)
}

// for sorting
type ByNumber []AnyFileInfo

func (b ByNumber) Len() int           { return len(b) }
func (b ByNumber) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByNumber) Less(i, j int) bool { return b[i].number < b[j].number }

const ARG_DEST_DEFAULT = "./"

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

var numberedFileRegex = regexp.MustCompile(`(.*?)[\.\s\-\_]+(\d+)$`)

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
	listFileParentFolder := filepath.Dir(listfile)
	for scanner.Scan() {
		s := scanner.Text()
		if !strings.HasPrefix(s, "#") {
			if !filepath.IsAbs(s) {
				s = joinpath(listFileParentFolder, s)
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
	return filepath.Join(source, target)
}

func mainRoutine() {
	fmt.Printf("Starting...\n")
	finfo, err := os.Stat(argsource)
	if err != nil { // no such file or dir
		log.Fatalf("File or folder %s does not exist: %v", argsource, err)
		return
	}

	fdest, err := os.Stat(argdest)
	if err != nil || !fdest.IsDir() {
		log.Fatalf("%s does not exist or is not a folder: %v", argdest, err)
		return
	}

	var listedFiles []AnyFileInfo
	var processed int
	if finfo.IsDir() {
		listedFiles, processed = doDir(argsource)
	} else {
		listedFiles, processed = doList(argsource)
	}

	sort.Sort(ByNumber(listedFiles))

	for lf := range listedFiles {
		fmt.Println(joinpath(argdest, listedFiles[lf].NewFileName(2)))
	}

	fmt.Printf("Destination folder is set to %s\n", filepath.Dir(argdest))
	fmt.Printf("Finished! %v files processed \n", processed)
}

func doDir(source string) (filesInDir []AnyFileInfo, processed int) {
	return withFilesInDir(source, filterNumbered)
}

func doList(source string) (filesInDir []AnyFileInfo, processed int) {
	fmt.Printf("Opening file %s, setting source folder to %s\n", source, filepath.Dir(source))
	listedFiles, processed := withFilesInList(source, analyzeListEntry)
	fmt.Printf("Total %v(%v)\n", len(listedFiles), processed)
	return listedFiles, processed
}

func main() {
	mainRoutine()
}
