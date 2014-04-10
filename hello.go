package main

import "fmt"
import "io/ioutil"
//import "os"

func main() {
	fmt.Printf("Starting...\n")
	var folder string = "./"
	dirlist,error := ioutil.ReadDir(folder)
	if (error == nil){
		fmt.Printf("%s\n",dirlist[0].Name())
	}else{
		fmt.Printf("Unable to read %s, error code %s", folder, error.Error())
	}
	fmt.Printf("Finished!\n")
}