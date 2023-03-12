package main

import (
	"fmt"
	"path/filepath"
	"os"
	"log"
)


func main() {
	currentDirectory, err := os.Getwd()

	if err != nil {
        log.Fatal(err)
    }

	targetDirectory := filepath.Join(currentDirectory, "xlsx")

	err = filepath.Walk(targetDirectory, lookForXlsx)

	if err != nil {
		fmt.Printf("Error while scanning  %q: %v\n", targetDirectory, err)
		return
	}
}


func lookForXlsx(path string, info os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if !info.IsDir(){
		if filepath.Ext(path) == ".xlsx"{
			fmt.Println(filepath.Base(path))
		}
	}

	return nil
}
